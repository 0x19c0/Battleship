package battleship

// Result is used to respond to attacks and is either "hit" or "miss"
type Result bool

const (
	hit = true
	miss = false
)

// Connection is a general communication interface for moves exchange, wrapping
// either TCP or UDP connections.
type Connection interface {
	CreateMessage(Result, Move) ([]byte)
	GetMessage() (Result, Move, error)
	SendMessage(Result, Move) (error)
}
