package battleship

import (
	"errors"
	"time"
)

var timeoutChan = 7 * time.Second

// ChanPlayer implements Connection interface via local chan.
type ChanPlayer struct {
	channel chan []byte
}

// NewChanPlayer method creates a new NewChanPlayer from byte slice channel.
func NewChanPlayer(channel chan []byte) *ChanPlayer {
	player := new(ChanPlayer)
	player.channel = channel
	return player
}

// GetMessage receives a message over a channel. Times out, reports errors.
func (player *ChanPlayer) GetMessage() (Result, Move, error) {
	timeout := time.After(timeoutChan)
	select {
	case message := <-player.channel:
		if len(message) != 3 {
			return miss, Move{}, errors.New("Message receiving error")
		}
		if message[0] == 0 {
			return miss, Move{row: int(message[1]), col: int(message[2])}, nil
		}
		return hit, Move{}, nil
	case <-timeout:
		return miss, Move{}, errors.New("Message receiving timeout")
	}
}

// SendMessage sends a message over a channel. Times out.
func (player *ChanPlayer) SendMessage(result Result, move Move) error {
	message := player.CreateMessage(result, move)
	timeout := time.After(timeoutChan)
	select {
	case player.channel <-message:
		return nil
	case <-timeout:
		return errors.New("Message sending timeout")
	}
	return nil
}

// CreateMessage creates a byte array encoding the result and the move.
func (player *ChanPlayer) CreateMessage(result Result, move Move) []byte {
	var message = make([]byte, 3)
	if result == hit {
		message[0] = 1
	} else {
		message[0] = 0
	}
	message[1] = byte(move.row)
	message[2] = byte(move.col)

	return message
}
