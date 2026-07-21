package ws

const (
	MsgCreateTable = "create_table"
	MsgJoinTable   = "join_table"
	MsgSit         = "sit"
	MsgSitBot      = "sit_bot"
	MsgStand       = "stand"
	MsgRemoveBot   = "remove_bot"
	MsgSetName     = "set_name"
	MsgSetDescription = "set_description"
	MsgStart       = "start"
	MsgBid         = "bid"
	MsgPlayCard    = "play_card"
	MsgChat        = "chat"

	MsgTableCreated = "table_created"
	MsgTableJoined  = "table_joined"
	MsgSeated       = "seated"
	MsgStood        = "stood"
	MsgGameState    = "game_state"
	MsgTableState   = "table_state"
	MsgChatMessage  = "chat_message"
	MsgChatHistory  = "chat_history"
	MsgError        = "error"
)

type ClientMessage struct {
	Type       string `json:"type"`
	TableID    string `json:"tableID,omitempty"`
	Direction  string `json:"direction,omitempty"`
	Call       string `json:"call,omitempty"`
	Card       string `json:"card,omitempty"`
	Name       string `json:"name,omitempty"`
	Text       string `json:"text,omitempty"`
	Description string `json:"description,omitempty"`
	Seed       *int64 `json:"seed,omitempty"`
	Difficulty int    `json:"difficulty,omitempty"`
}

type ServerMessage struct {
	Type    string `json:"type"`
	Payload any    `json:"payload,omitempty"`
	Error   string `json:"error,omitempty"`
}
