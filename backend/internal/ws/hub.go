package ws

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"sync"
	"time"

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
	mu        sync.Mutex
	tables    map[string]*Table
	gameRepo  *repository.GameRepository
	reapGrace time.Duration // tests shrink this right after NewHub
	// leaveGrace delays a "left" chat line so a quick reconnect (page
	// navigation, refresh) doesn't spam left/joined. Tests shrink it.
	leaveGrace time.Duration
	users      UserLookup // nil: every connection is a guest
}

func NewHub(gameRepo *repository.GameRepository) *Hub {
	return &Hub{
		tables:     make(map[string]*Table),
		gameRepo:   gameRepo,
		reapGrace:  30 * time.Second,
		leaveGrace: defaultLeaveGrace,
	}
}

// SetUserLookup lets the hub resolve logged-in users' usernames.
func (h *Hub) SetUserLookup(u UserLookup) {
	h.users = u
}

func (h *Hub) CreateTable() *Table {
	h.mu.Lock()
	defer h.mu.Unlock()

	id := generateTableID()
	t := NewTable(id, h.gameRepo)
	t.leaveGrace = h.leaveGrace
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

// reapIfEmpty schedules a table for removal once no clients are watching
// it, after a grace period (Hub.reapGrace) so a table survives a brief
// gap — a page navigation, or an invite link opened just after the
// creator closed their tab. AddObserver cancels a pending reap.
func (h *Hub) reapIfEmpty(t *Table) {
	if t.HasObservers() {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.reapTimer != nil {
		return
	}
	t.reapTimer = time.AfterFunc(h.reapGrace, func() {
		if t.HasObservers() {
			return
		}
		h.RemoveTable(t.ID)
		t.Shutdown()
	})
}

// HandleUpgrade upgrades an HTTP connection to a WebSocket and
// starts the client read/write pumps.
func (h *Hub) HandleUpgrade(w http.ResponseWriter, r *http.Request) {
	// A valid token cookie (via OptionalMiddleware) attaches the real
	// user ID. Guests connect without one and get a throwaway ID that
	// the table layer never reads.
	userID, authed := auth.UserIDFromContext(r.Context())
	if !authed {
		userID = uuid.New()
	}
	name, guest := h.lookupUsername(r.Context(), userID, authed)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("websocket upgrade error: %v", err)
		return
	}

	client := NewClient(h, conn, userID)
	client.guest = guest
	if !guest {
		client.name = name
		client.identity = "u:" + userID.String()
	}
	go client.writePump()
	go client.readPump()
}

// lookupUsername resolves a logged-in user's username. Any failure falls
// back to guest, so a DB hiccup never blocks someone from playing.
func (h *Hub) lookupUsername(ctx context.Context, id uuid.UUID, authed bool) (string, bool) {
	if !authed || h.users == nil {
		return "", true
	}
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	u, err := h.users.GetByID(ctx, id)
	if err != nil {
		log.Printf("username lookup for %v failed, connecting as guest: %v", id, err)
		return "", true
	}
	return u.Username, false
}

func generateTableID() string {
	b := make([]byte, 4)
	rand.Read(b)
	return hex.EncodeToString(b)
}
