package battleship

import (
	"errors"
	"fmt"
)

// Rules type represents specific set of rules for current game
type Rules struct {
	boardSize     int
	shipsTypesNum int
	maxHealth     int
	shipsArray    []int
}

// Specific set of rules used for the game
var rules = Rules{10, 4, 20, []int{4, 3, 2, 1}}

// Game contains the current state of the game and the connection to the enemy
type Game struct {
	ownBoard   Board
	enemyBoard Board

	ownHealth   int
	enemyHealth int

	moveNum int

	conn   Connection
	isHost bool
}

// NewGame creates a new game object and initializes it with provided board
// and connection
func NewGame(boardFile string, conn Connection, isHost bool) (*Game, error) {
	game := new(Game)

	b, err := scanBoard(boardFile)
	if err != nil {
		return game, errors.New("In " + boardFile + ":\n" + err.Error())
	}

	err = checkBoard(b)
	if err != nil {
		return game, errors.New("In " + boardFile + ":\n" + err.Error())
	}

	game.ownBoard = b
	game.ownHealth, game.enemyHealth = rules.maxHealth, rules.maxHealth
	game.moveNum = 0

	game.enemyBoard = newEnemyBoard(rules.boardSize, rules.boardSize)

	game.conn = conn
	game.isHost = isHost

	return game, nil
}

// report prints both boards
func (game *Game) report() {
	fmt.Println("\nYOUR BOARD:\n")
	game.ownBoard.print()
	fmt.Println("\nENEMY BOARD:\n")
	game.enemyBoard.print()
	fmt.Println()
}

// Play runs a general flow of the game until one of the players is dead.
// It currently reports the result. It probably shouldn't.
func (game *Game) Play() {
	var currentMove = illegalMove

	if game.isHost != true {
		fmt.Println("You start.")
		currentMove = scanMove()
		err := game.conn.SendMessage(miss, currentMove)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	for {
		result, enemyMove, err := game.conn.GetMessage()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if result == miss {
			if currentMove != illegalMove {
				game.enemyBoard[currentMove.row][currentMove.col] = empty
			}
			if game.ownBoard[enemyMove.row][enemyMove.col] == ship {
				game.ownBoard[enemyMove.row][enemyMove.col] = wreck
				game.ownHealth--
				fmt.Println("Your opponent hit your ship and moves again.")
				fmt.Println("Your health:", game.ownHealth)
				game.report()
				err := game.conn.SendMessage(hit, Move{})
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				if game.ownHealth == 0 {
					fmt.Println("All your ships are lost. Game over.")
					return
				}
			} else {
				fmt.Println("Your opponent missed.")
				game.report()
				currentMove = scanMove()
				err := game.conn.SendMessage(miss, currentMove)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
			}
		} else {
			game.enemyBoard[currentMove.row][currentMove.col] = wreck
			game.enemyHealth--
			fmt.Println("You hit an enemy ship and move again.")
			fmt.Println("Enemy health:", game.enemyHealth)
			game.report()
			if game.enemyHealth == 0 {
				fmt.Println("You destroyed all enemy ship. Congratulations.")
				return
			}
			currentMove = scanMove()
			err := game.conn.SendMessage(miss, currentMove)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
	}
}
