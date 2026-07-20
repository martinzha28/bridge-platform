package ws

import (
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
}

func NewHub(gameRepo *repository.GameRepository) *Hub {
	return &Hub{
		tables:    make(map[string]*Table),
		gameRepo:  gameRepo,
		reapGrace: 30 * time.Second,
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
