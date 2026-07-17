package ws

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/martinzha28/bridge-platform/backend/internal/auth"
	"github.com/martinzha28/bridge-platform/backend/internal/repository"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Hub struct {
	mu       sync.Mutex
	tables   map[string]*Table
	gameRepo *repository.GameRepository
}

func NewHub(gameRepo *repository.GameRepository) *Hub {
	return &Hub{
		tables:   make(map[string]*Table),
		gameRepo: gameRepo,
	}
}

func (h *Hub) CreateTable() *Table {
	h.mu.Lock()
	defer h.mu.Unlock()

	id := generateTableID()
	t := NewTable(id, h.gameRepo)
	h.tables[id] = t
	return t
}

func (h *Hub) GetTable(id string) (*Table, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()

	t, ok := h.tables[id]
	return t, ok
}

func (h *Hub) RemoveTable(id string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.tables, id)
}

// reapIfEmpty drops a table once no live client connections remain
// seated, stopping its bot goroutines. Guests can create tables
// freely, so abandoned ones must not linger.
func (h *Hub) reapIfEmpty(t *Table) {
	if t.HasHumans() {
		return
	}
	h.RemoveTable(t.ID)
	t.Shutdown()
}

// HandleUpgrade upgrades an HTTP connection to a WebSocket and
// starts the client read/write pumps.
func (h *Hub) HandleUpgrade(w http.ResponseWriter, r *http.Request) {
	// A valid token cookie (via OptionalMiddleware) attaches the real
	// user ID. Guests connect without one and get a throwaway ID that
	// the table layer never reads.
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		userID = uuid.New()
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("websocket upgrade error: %v", err)
		return
	}

	client := NewClient(h, conn, userID)
	go client.writePump()
	go client.readPump()
}

func generateTableID() string {
	b := make([]byte, 4)
	rand.Read(b)
	return hex.EncodeToString(b)
}
