package battleship

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
)

const (
	empty   = "-"
	ship    = "O"
	wreck   = "X"
	unknown = "?"
)

// Board represents the game board
type Board [][]string

func newBoard(rows, cols int) Board {
	b := make([][]string, rows)
	for i := range b {
		b[i] = make([]string, cols)
	}
	return b
}

func newEnemyBoard(rows, cols int) Board {
	b := make([][]string, rows)
	for i := range b {
		b[i] = make([]string, cols)
		for j := range b[i] {
			b[i][j] = unknown
		}
	}
	return b
}
// noCollision checks whether there is a ship collision on given indices
func (b Board) noCollision(i, j int) bool {
	if (i > 0 && j > 0 && b[i][j-1] != empty && b[i-1][j] != empty) ||
		(i > 0 && j < rules.boardSize-1 && b[i][j+1] != empty && b[i-1][j] != empty) ||
		(i < rules.boardSize-1 && j > 0 && b[i][j-1] != empty && b[i+1][j] != empty) ||
		(i < rules.boardSize-1 && j < rules.boardSize-1 && b[i][j+1] != empty && b[i+1][j] != empty) {
		return false
	}
	return true
}

// isShipHead indicates whether the point on given indices is a last point
// of the ship (right side of horizontal, bottom side of vertical)
func (b Board) isShipHead(i, j int) bool {
	if (i < rules.boardSize-1 && b[i+1][j] != empty) ||
		(j < rules.boardSize-1 && b[i][j+1] != empty) ||
		(b[i][j] == empty) {
		return false
	}
	return true
}

// shipSize returns a size of a ship that is finished on given indices
func (b Board) shipSize(i, j int) int {
	var size int
	if i > 0 && b[i-1][j] != empty {
		for i >= 0 && b[i][j] == ship {
			size++
			i--
		}
		return size
	}
	for j >= 0 && b[i][j] == ship {
		size++
		j--
	}
	return size
}

// checkBoard validates a board according to game rules and reports
// the violations if found.
func checkBoard(board Board) error {
	ships := make([]int, rules.shipsTypesNum)

	for i := 0; i < rules.boardSize; i++ {
		for j := 0; j < rules.boardSize; j++ {
			if board[i][j] == empty {
				continue
			}
			if board.noCollision(i, j) != true {
				return errors.New("Colliding ships: row" + strconv.Itoa(i) +
					" column " + strconv.Itoa(j))
			}
			if board.isShipHead(i, j) == true {
				ships[board.shipSize(i, j)-1]++
			}
		}
	}

	for i, val := range ships {
		if val != rules.shipsArray[i] {
			return errors.New("Wrong amount of ships of size " + strconv.Itoa(i+1) +
				", must be " + strconv.Itoa(rules.shipsArray[i]) + ", is " + strconv.Itoa(val))
		}
	}

	return nil
}

// scanBoard creates a board from file if it's properly formatted, or reports
// any formatting errors otherwise.
func scanBoard(filename string) (Board, error) {
	b := newBoard(rules.boardSize, rules.boardSize)

	inFile, err := os.Open(filename)
	if err != nil {
		return b, err
	}
	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)

	var row int
	for ; scanner.Scan(); row++ {
		if row > rules.boardSize {
			return b, errors.New("In input file " + filename + ", row " +
				strconv.Itoa(row+1) + " is too long")
		}

		str := scanner.Text()
		if len(str) != rules.boardSize {
			return b, errors.New("In input file " + filename + ", row " +
				strconv.Itoa(row+1) + " is of invalid length " +
				strconv.Itoa(len(str)) + ", must be " + strconv.Itoa(rules.boardSize))
		}

		for col, rune := range str {
			switch rune {
			case '1':
				b[row][col] = ship
			case '0':
				b[row][col] = empty
			default:
				return b, errors.New("In input file " + filename +
					", row " + strconv.Itoa(row+1) + " has illegal character " + string(rune))
			}
		}
	}

	if row < rules.boardSize {
		return b, errors.New("The input file " + filename + " is too short")
	}

	return b, nil
}

func (b Board) print() {
	for _, row := range b {
		for _, symbol := range row {
			fmt.Printf("%s", symbol)
		}
		fmt.Printf("\n")
	}
}
