package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

var wg = sync.WaitGroup{}

// -1 not a field, 0 no figure, 1 figure

var board = [7][7]int8{
	{-1, -1, 1, 1, 1, -1, -1},
	{-1, -1, 1, 1, 1, -1, -1},
	{0, 1, 1, 1, 1, 1, 1},
	{0, 0, 0, 1, 1, 1, 1},
	{0, 0, 1, 0, 1, 1, 1},
	{-1, -1, 0, 1, 1, -1, -1},
	{-1, -1, 0, 1, 1, -1, -1},
}

var numberMoves = 0
var numberGames int64 = 0
var bestResult int8 = 32
var start time.Time
var numberThreads = 0

func ident(cntMoves int) {
	for i := 0; i < cntMoves; i++ {
		fmt.Print(" ")
	}
}

func makeAllPossibleMoves(board *[7][7]int8, cntMoves int, game string, isThread bool) {

	//ident(cntMoves)
	//fmt.Printf("game %v: search moves (so far %v moves)", game, cntMoves)
	//printBoard(board)

	//if cntMoves < 15 {
	//	fmt.Printf("---> level %v reached \n", cntMoves)
	//	fmt.Printf("---> game %v\n", game)
	//}

	var moves [][4]int8

	moves = findAllPossibleMoves(board)
	//fmt.Printf(" -> found %v moves\n", len(moves))

	if len(moves) == 0 {
		//ident(cntMoves)
		//fmt.Printf("no more moves found for game %v:\n", game)
		//printBoard(board)

		remaining := countRemaining(board)

		if remaining == 1 {
			fmt.Printf(">>>>>>>>>>>>> FOUND A WINNING GAME <<<<<<<<<<<<<<<<\n")
			fmt.Printf(">>>>>>>>>>>>> after %v games    <<<<<<<<<<<<<<<<\n", numberGames)
			printBoard(board)
		}

		if remaining < bestResult {
			bestResult = remaining
			fmt.Printf(" ===== new bestResult: %v\n", bestResult)
		}

		numberGames++
		if numberGames%1000000 == 0 {
			fmt.Printf("no more moves found for after %v million games (%s / %vk games per seconds) %v threads\n",
				numberGames/1000000,
				time.Since(start),
				((numberGames / 1000) / ((time.Since(start).Milliseconds() / 1000) + 1)),
				numberThreads)
		}
	}

	for cntM, element := range moves {
		var newBoard [7][7]int8

		// copy into new array because I don't know how to do it better
		var i, j int8
		for i = 0; i < 7; i++ {
			for j = 0; j < 7; j++ {
				newBoard[i][j] = board[i][j]
			}
		}

		//fmt.Printf("copied board for game %v:\n", game)
		//printBoard(board)
		//fmt.Printf("newBoard\n")
		//printBoard(newBoard)
		makeMove(&newBoard, element[0], element[1], element[2], element[3])
		//fmt.Printf("after makeMove for game %v:\n", game)
		//printBoard(board)
		//fmt.Printf("newBoard\n")
		//printBoard(newBoard)
		if cntMoves < 2 {
			numberThreads++
			wg.Add(1)
			go makeAllPossibleMoves(&newBoard, cntMoves+1, game+"."+strconv.Itoa(cntM), true)
		} else {
			makeAllPossibleMoves(&newBoard, cntMoves+1, game+"."+strconv.Itoa(cntM), false)
		}

	}

	if isThread {
		wg.Done()
	}
}

func playOneGame(board [7][7]int8) {
	//fmt.Println("========== search searching with board: ")
	//printBoard(&board)

	var fromI, fromJ, toI, toJ int8

	fromI, fromJ, toI, toJ = findAMove(&board)

	if fromI >= 0 {
		makeMove(&board, fromI, fromJ, toI, toJ)
		printBoard(&board)
		playOneGame(board)
	} else {
		fmt.Println("no more moves found")
		printBoard(&board)
	}
}

func makeMove(board *[7][7]int8, fromI int8, fromJ int8, toI int8, toJ int8) {
	board[fromI][fromJ] = 0
	board[toI][toJ] = 1
	// remove the middle guy
	if fromI == toI {
		if fromJ > toJ {
			board[fromI][fromJ-1] = 0
		} else {
			board[fromI][fromJ+1] = 0
		}
	}
	if fromJ == toJ {
		if fromI > toI {
			board[fromI-1][fromJ] = 0
		} else {
			board[fromI+1][fromJ] = 0
		}
	}

	//fmt.Printf("made move from %v/%v to %v/%v\n", fromI, fromJ, toI, toJ)
	numberMoves++
}

// returns fromY, fromX, toY, toX
func findAMove(board *[7][7]int8) (int8, int8, int8, int8) {

	var i, j int8
	for i = 0; i < 7; i++ {
		for j = 0; j < 7; j++ {
			if board[i][j] == 0 {

				//we now must find 2 consecute pieces to make a move
				// so lets check each direction
				// i is row, j is column

				// first north
				if i >= 2 && board[i-1][j] > 0 && board[i-2][j] > 0 {
					return i - 2, j, i, j
				}
				// then south
				if i <= 4 && board[i+1][j] > 0 && board[i+2][j] > 0 {
					return i + 2, j, i, j
				}
				// then west
				if j >= 2 && board[i][j-1] > 0 && board[i][j-2] > 0 {
					return i, j - 2, i, j
				}
				//then east
				if j <= 4 && board[i][j+1] > 0 && board[i][j+2] > 0 {
					return i, j + 2, i, j
				}
			}
		}
	}
	return -1, -1, -1, -1
}

// returns a slice of fromY, fromX, toY, toX
func findAllPossibleMoves(board *[7][7]int8) [][4]int8 {

	var moves [][4]int8

	var i, j int8
	for i = 0; i < 7; i++ {
		for j = 0; j < 7; j++ {
			if board[i][j] == 0 {

				//we now must find 2 consecute pieces to make a move
				// so lets check each direction
				// i is row, j is column

				// firstreturn -1, -1, -1, -1 north
				if i >= 2 && board[i-1][j] > 0 && board[i-2][j] > 0 {
					moves = append(moves, [4]int8{i - 2, j, i, j})
				}
				// then south
				if i <= 4 && board[i+1][j] > 0 && board[i+2][j] > 0 {
					moves = append(moves, [4]int8{i + 2, j, i, j})
				}
				// then west
				if j >= 2 && board[i][j-1] > 0 && board[i][j-2] > 0 {
					moves = append(moves, [4]int8{i, j - 2, i, j})
				}
				//then east
				if j <= 4 && board[i][j+1] > 0 && board[i][j+2] > 0 {
					moves = append(moves, [4]int8{i, j + 2, i, j})
				}
			}
		}
	}
	return moves
}

func isWon(board *[7][7]int8) bool {

	var i, j int8
	var sum int

	for i = 0; i < 7; i++ {
		for j = 0; j < 7; j++ {
			if board[i][j] > 0 {
				sum++
			}
		}
	}

	if sum == 1 {
		return true
	} else {
		return false
	}
}

func countRemaining(board *[7][7]int8) int8 {

	var i, j int8
	var sum int8

	for i = 0; i < 7; i++ {
		for j = 0; j < 7; j++ {
			if board[i][j] > 0 {
				sum++
			}
		}
	}

	return sum
}

func printBoard(board *[7][7]int8) {

	var i, j int8

	for i = 0; i < 7; i++ {
		for j = 0; j < 7; j++ {
			if board[i][j] < 0 {
				fmt.Printf("-")
			} else if board[i][j] > 0 {
				fmt.Printf("X")
			} else {
				fmt.Printf(" ")
			}

		}
		fmt.Println()
	}

}

func main() {
	fmt.Println("Starting Solitaire")
	start = time.Now()
	//printBoard(&board)
	makeAllPossibleMoves(&board, 0, "1", false)
	wg.Wait()
	fmt.Printf("finished after %v games and %v moves\n", numberGames, numberMoves)
}
