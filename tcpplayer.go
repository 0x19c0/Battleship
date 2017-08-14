package battleship

import (
	"net"
	"errors"
	"time"
)

var timeoutTCP = 15 * time.Second

// TCPPlayer implements Connection interface via TCP.
type TCPPlayer struct {
	conn net.Conn
}

// NewTCPPlayer method creates a new TCPPlayer from net.Conn.
func NewTCPPlayer(conn net.Conn) *TCPPlayer {
	player := new(TCPPlayer)
	player.conn = conn
	return player
}

// GetMessage gets a message over a TCP connection. Times out, reports errors.
func (player *TCPPlayer) GetMessage() (Result, Move, error) {
	buffer := make([]byte, 3)

	err := player.conn.SetReadDeadline(time.Now().Add(timeoutTCP))
	if err != nil {
		return false, Move{}, err
	}
	n, err := player.conn.Read(buffer)
	if err != nil {
		return false, Move{}, err
	}
	if n != 3 {
		return false, Move{}, errors.New("Message receiving error")
	}

	if buffer[0] == 0 {
		return miss, Move{row: int(buffer[1]), col: int(buffer[2])}, nil
	}
	return hit, Move{}, nil
}

// SendMessage sends a message over TCP connection. Times out, reports errors.
func (player *TCPPlayer) SendMessage(result Result, move Move) (error) {
	message := player.CreateMessage(result, move)

	err := player.conn.SetWriteDeadline(time.Now().Add(timeoutTCP))
	if err != nil {
		return err
	}
	n, err := player.conn.Write(message)
	if err != nil {
		return err
	}
	if n != 3 {
		return errors.New("Message sending error")
	}

	return nil
}

// CreateMessage creates a byte array encoding the result and the move.
func (player *TCPPlayer) CreateMessage(result Result, move Move) ([]byte) {
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
