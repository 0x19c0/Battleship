package battleship

import (
	"bufio"
	"fmt"
	"os"
)

// Move represents a single attack.
type Move struct {
	row int
	col int
}

var illegalMove = Move{row: -1, col: -1}

// scanMove gets an input from user until it's properly formatted and valid,
// then creates a Move from it.
func scanMove() Move {
	fmt.Println("Please, enter your move:")
	for {
		var r, c int
		stdin := bufio.NewReader(os.Stdin)
		_, err := fmt.Fscanf(stdin, "%d %d\n", &r, &c)
		if err != nil {
			fmt.Println("Invalid input format. Correct format is: \"row col\". Try again:")
			stdin.ReadString('\n')
			continue
		}
		if r >= rules.boardSize || c >= rules.boardSize || r < 0 || c < 0 {
			fmt.Println("Both values should be between 0 and", rules.boardSize-1, ". Try again:")
			continue
		}
		return Move{row: r, col: c}
	}
}
