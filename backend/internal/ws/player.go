package ws

// Player is implemented by both Client (real WebSocket connection) and
// BotClient (in-process bot). Table holds a map[Direction]Player so it
// doesn't need to distinguish between the two.
type Player interface {
	Send(data []byte)
}
