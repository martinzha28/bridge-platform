package ws

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/martinzha28/bridge-platform/backend/internal/bot"
	"github.com/martinzha28/bridge-platform/backend/internal/game"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4096
	sendBufferSize = 256
)

type Client struct {
	hub       *Hub
	conn      *websocket.Conn
	send      chan []byte
	userID    uuid.UUID
	table     *Table
	direction game.Direction
	seated    bool
}

func NewClient(hub *Hub, conn *websocket.Conn, userID uuid.UUID) *Client {
	return &Client{
		hub:    hub,
		conn:   conn,
		send:   make(chan []byte, sendBufferSize),
		userID: userID,
	}
}

func (c *Client) readPump() {
	defer func() {
		c.disconnect()
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("websocket read error: %v", err)
			}
			return
		}

		var msg ClientMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			c.sendError("invalid JSON")
			continue
		}
		c.handleMessage(msg)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) handleMessage(msg ClientMessage) {
	switch msg.Type {
	case MsgCreateTable:
		table := c.hub.CreateTable()
		c.table = table
		c.sendJSON(ServerMessage{
			Type:    MsgTableCreated,
			Payload: map[string]string{"tableID": table.ID},
		})

	case MsgJoinTable:
		table, ok := c.hub.GetTable(msg.TableID)
		if !ok {
			c.sendError("table not found: " + msg.TableID)
			return
		}
		c.table = table
		c.sendJSON(ServerMessage{
			Type:    MsgTableJoined,
			Payload: map[string]string{"tableID": table.ID},
		})

	case MsgSit:
		if c.table == nil {
			c.sendError("join a table first")
			return
		}
		if len(msg.Direction) == 0 {
			c.sendError("direction is required")
			return
		}
		dir, ok := game.ParseDirection(msg.Direction[0])
		if !ok {
			c.sendError("invalid direction: " + msg.Direction)
			return
		}
		if err := c.table.Sit(Player(c), dir); err != nil {
			c.sendError(err.Error())
			return
		}
		c.direction = dir
		c.seated = true
		c.sendJSON(ServerMessage{
			Type:    MsgSeated,
			Payload: map[string]string{"direction": dir.String()},
		})

	case MsgSitBot:
		if c.table == nil {
			c.sendError("join a table first")
			return
		}
		if len(msg.Direction) == 0 {
			c.sendError("direction is required")
			return
		}
		dir, ok := game.ParseDirection(msg.Direction[0])
		if !ok {
			c.sendError("invalid direction: " + msg.Direction)
			return
		}
		difficulty, _ := bot.ParseDifficulty(msg.Difficulty)
		if err := c.table.SitBot(dir, difficulty); err != nil {
			c.sendError(err.Error())
			return
		}
		c.sendJSON(ServerMessage{
			Type:    MsgSeated,
			Payload: map[string]string{"direction": dir.String(), "bot": "true"},
		})

	case MsgStart:
		if c.table == nil || !c.seated {
			c.sendError("not seated at a table")
			return
		}
		if err := c.table.Start(msg.Seed); err != nil {
			c.sendError(err.Error())
			return
		}

	case MsgBid:
		if c.table == nil || !c.seated {
			c.sendError("not seated at a table")
			return
		}
		call, ok := game.ParseCall(msg.Call)
		if !ok {
			c.sendError("invalid call: " + msg.Call)
			return
		}
		if err := c.table.Bid(c.direction, call); err != nil {
			c.sendError(err.Error())
			return
		}

	case MsgPlayCard:
		if c.table == nil || !c.seated {
			c.sendError("not seated at a table")
			return
		}
		card, ok := game.ParseCard(msg.Card)
		if !ok {
			c.sendError("invalid card: " + msg.Card)
			return
		}
		if err := c.table.PlayCard(c.direction, card); err != nil {
			c.sendError(err.Error())
			return
		}

	default:
		c.sendError("unknown message type: " + msg.Type)
	}
}

func (c *Client) Send(data []byte) {
	select {
	case c.send <- data:
	default:
	}
}

func (c *Client) sendJSON(msg ServerMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}
	select {
	case c.send <- data:
	default:
	}
}

func (c *Client) sendError(message string) {
	c.sendJSON(ServerMessage{
		Type:  MsgError,
		Error: message,
	})
}

func (c *Client) disconnect() {
	if c.table != nil {
		c.table.RemovePlayer(c)
		c.hub.reapIfEmpty(c.table)
	}
	close(c.send)
}
