package ws

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/martinzha28/bridge-platform/backend/internal/models"
)

const (
	minNameLen    = 2
	maxNameLen    = 20
	maxGuestIDLen = 64
)

// UserLookup resolves a logged-in user's account. *repository.UserRepository
// satisfies it; tests pass a fake.
type UserLookup interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
}

// sanitizeDisplayName trims and collapses whitespace, then accepts 2–20
// letters, digits, spaces, '_' or '-'.
func sanitizeDisplayName(s string) (string, bool) {
	name := strings.Join(strings.Fields(s), " ")
	if n := utf8.RuneCountInString(name); n < minNameLen || n > maxNameLen {
		return "", false
	}
	for _, r := range name {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != ' ' && r != '_' && r != '-' {
			return "", false
		}
	}
	return name, true
}

func randomGuestName() string {
	n, err := rand.Int(rand.Reader, big.NewInt(9000))
	if err != nil {
		return "Guest 1000"
	}
	return fmt.Sprintf("Guest %d", 1000+n.Int64())
}

// setGuestIdentity applies the nickname and browser guest ID a guest sends
// on create/join. Registered users keep their username; this is a no-op.
func (c *Client) setGuestIdentity(displayName, guestID string) {
	if !c.guest {
		return
	}
	name, ok := sanitizeDisplayName(displayName)
	if !ok {
		name = randomGuestName()
	}
	c.name = name

	guestID = strings.TrimSpace(guestID)
	if guestID != "" && len(guestID) <= maxGuestIDLen {
		c.identity = "g:" + guestID
	} else {
		c.identity = "c:" + c.userID.String()
	}
}
