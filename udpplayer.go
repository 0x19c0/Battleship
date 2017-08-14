package battleship

import (
	"net"
	"sync/atomic"
	"time"
)

var timeoutUDP = 15 * time.Second
var cooldownUDP = 2 * time.Second

// UDPPlayer implements Connection interface via UDP.
// It sents indexed messages on a background with a set cooldown
// Once a response is received, it stops sending current message and starts
// with the next one.
type UDPPlayer struct {
	connIn  *net.UDPConn
	connOut net.Conn
	index   int32
}

// NewUDPPlayer method creates a new UDPPlayer.
func NewUDPPlayer(connIn *net.UDPConn, connOut net.Conn) *UDPPlayer {
	player := new(UDPPlayer)
	player.connIn = connIn
	player.connOut = connOut
	return player
}

// GetMessage receives messages from UDP connection until a valid message
// with a valid index is received. Updates the index of connection closing the
// previous sending routine.
func (player *UDPPlayer) GetMessage() (Result, Move, error) {
	buffer := make([]byte, 4)

	defer atomic.AddInt32(&(player.index), 1)

	for {
		err := player.connIn.SetReadDeadline(time.Now().Add(timeoutUDP))
		if err != nil {
			return false, Move{}, err
		}
		n, _, err := player.connIn.ReadFrom(buffer)
		if err != nil {
			return false, Move{}, err
		}
		if n != 4 || int32(buffer[3]) != atomic.LoadInt32(&(player.index)) {
			time.Sleep(10 * time.Millisecond)
			continue
		}

		if buffer[0] == 0 {
			return miss, Move{row: int(buffer[1]), col: int(buffer[2])}, nil
		}
		return hit, Move{}, nil
	}
}

// SendMessage sends messages over UDP connection with a preset cooldown.
func (player *UDPPlayer) SendMessage(result Result, move Move) error {
	atomic.AddInt32(&(player.index), 1)

	message := player.CreateMessage(result, move)

	go func() {
		for {
			if atomic.LoadInt32(&(player.index)) != int32(message[3])+1 {
				break
			} else {
				player.connOut.Write(message)
				time.Sleep(cooldownUDP)
			}
		}
		return
	}()

	return nil
}

// CreateMessage creates a byte array encoding the result, the move and the
// index of the move.
func (player *UDPPlayer) CreateMessage(result Result, move Move) ([]byte) {
	var message = make([]byte, 3)
	if result == hit {
		message[0] = 1
	} else {
		message[0] = 0
	}
	message[1] = byte(move.row)
	message[2] = byte(move.col)
	message[3] = byte(atomic.LoadInt32(&(player.index)))

	return message
}
