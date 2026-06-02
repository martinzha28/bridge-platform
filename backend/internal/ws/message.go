package ws

const (
	MsgCreateTable = "create_table"
	MsgJoinTable   = "join_table"
	MsgSit         = "sit"
	MsgSitBot      = "sit_bot"
	MsgStart       = "start"
	MsgBid         = "bid"
	MsgPlayCard    = "play_card"

	MsgTableCreated = "table_created"
	MsgTableJoined  = "table_joined"
	MsgSeated       = "seated"
	MsgGameState    = "game_state"
	MsgError        = "error"
)

type ClientMessage struct {
	Type       string `json:"type"`
	TableID    string `json:"tableID,omitempty"`
	Direction  string `json:"direction,omitempty"`
	Call       string `json:"call,omitempty"`
	Card       string `json:"card,omitempty"`
	Seed       *int64 `json:"seed,omitempty"`
	Difficulty int    `json:"difficulty,omitempty"`
}

type ServerMessage struct {
	Type    string `json:"type"`
	Payload any    `json:"payload,omitempty"`
	Error   string `json:"error,omitempty"`
}
