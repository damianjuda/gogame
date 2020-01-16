package game

import (
	"fmt"
	"github.com/Workiva/go-datastructures/queue"
	"math/rand"
	"time"
)

type dir int

const (
	UD dir = 0 // pionowe domino
	LR dir = 1 // poziome domino
)

type domino struct {
	index int    // indeks na 1 wymiarowej tablicy (reprezentującej 2 wymiarową) prawy górny róg domino
	dir   dir    // poziome albo pionowe domino
	board *Board // referencja do całej tablicy
}

func NewDomino(x int, y int, board *Board) *domino {
	if x+board.columns == y {
		return &domino{y, UD, board}
	} else if y == x+board.columns {
		return &domino{x, UD, board}
	} else if y == x+1 {
		return &domino{x, LR, board}
	} else if x == y+1 {
		return &domino{y, LR, board}
	} else {
		panic("It's invalid domino")
	}
}

func (domino *domino) indicies() (int, int) {
	if domino.dir == UD {
		return domino.index, domino.index + domino.board.columns
	}
	return domino.index, domino.index + 1
}

type Board struct {
	rows            int // liczba wierszy kratownicy
	columns         int // liczba kolumn kratownicy
	clustersNoCache int
	board           []bool // tablica reprezentująca puste i pełne miejsca na kratownicy
	dominos         []*domino
}

func (board *Board) index(row int, column int) int {
	return row*board.columns + column // indeks w 1 wymiarowej tablicy emulującej 2 wymiarową
}

func (board *Board) mindex(index int) (int, int) {
	return index / board.columns, index % board.columns
}

func (board *Board) valid(row int, column int) bool {
	return row < board.rows && column < board.columns && row >= 0 && column >= 0 // czy pozycja (wiersz,kolumna) nalezy do kratownicy
}

func (board *Board) neighbours(index int) []int {
	row, column := board.mindex(index)
	left := []int{row, column - 1}
	right := []int{row, column + 1}
	up := []int{row - 1, column}
	down := []int{row + 1, column}
	neighbours := [][]int{left, right, up, down}
	validNeighbours := make([]int, 0)
	for _, neighbour := range neighbours {
		if board.valid(neighbour[0], neighbour[1]) {
			validNeighbour := board.index(neighbour[0], neighbour[1])
			validNeighbours = append(validNeighbours, validNeighbour)
		}
	}
	return validNeighbours
}

func (board *Board) cluster(index int) int {
	cluster := 0
	todo := make([]int, 10)
	item := index
	todo = append(todo, item)
	for len(todo) > 0 {
		item, todo = todo[0], todo[1:]
		if board.board[item] == false {
			board.board[item] = true
			cluster += 1
			todo = append(todo, board.neighbours(item)...)
		}
	}
	return cluster
}

func (board *Board) Mark(row int, column int) {
	board.board[board.index(row, column)] = true // zaznacz pole jako zajęte
}

func EmptyBoard(rows int, columns int) *Board {
	// pusta kratonica
	return &Board{
		rows,
		columns,
		-2,
		make([]bool, rows*columns, rows*columns),
		make([]*domino, 0),
	}
}

func copyBoard(base *Board) *Board {
	// skopiuj kratownice
	board := EmptyBoard(base.rows, base.columns)
	for index, value := range base.board {
		board.board[index] = value
	}
	for _, value := range base.dominos {
		board.dominos = append(board.dominos, value)
	}
	return board
}

func (b *Board) moves() []*domino {
	// wszystkie nowe domina ktore mozna dolozyc do kratownicy
	dominos := []*domino{}
	for index := range b.board {
		// dla kazdego pola jesli jest puste (zakladac, ze jest prawy gornym rogiem nowego klocka)
		if !b.board[index] {
			// sprawdz czy klocek sie miesci na kratownicy i jest na niego miejsce
			right := index + 1
			if b.valid(index/b.columns, index%b.columns+1) && !b.board[right] {
				dominos = append(dominos, &domino{index, LR, b})
			}
			down := index + b.columns
			if b.valid(index/b.columns+1, index%b.columns) && !b.board[down] {
				dominos = append(dominos, &domino{index, UD, b})
			}
		}
	}
	return dominos
}

func (board *Board) full() float64 {
	counter := 0.0
	for index := range board.board {
		if board.board[index] {
			counter += 1
		}
	}
	return counter
}

func (board *Board) Hash() string {
	// Niezbyt madra funkcja generujaca unikalnego hasha (stringa) dla kazdych 2 roznych kratownic
	hash := fmt.Sprintf("#%d ", len(board.dominos))
	for _, domino := range board.dominos {
		hash += fmt.Sprintf("(%d%d)", domino.dir, domino.index)
	}
	return hash
}

func (board *Board) IsSolution() bool {
	// czy wszystkie pola na kratownicy sa zajete
	return board.clustersNo() == 0
}

func (board *Board) Reject() bool {
	return board.clustersNo() == -1
}

func (board *Board) clustersNo() int {
	if (*board).clustersNoCache == -2 {
		test := copyBoard(board)
		clustersNo := 0
		for index := range test.board {
			clusterSize := test.cluster(index)
			if clusterSize%2 == 1 {
				clustersNo = -1
				break
			} else if clusterSize > 0 {
				clustersNo += 1
			}
		}
		(*board).clustersNoCache = clustersNo
	}
	return (*board).clustersNoCache
}

func (board *Board) Compare(other queue.Item) int {
	otherBoard := other.(*Board)
	if board.full() > otherBoard.full() {
		return -1
	} else if board.full() < otherBoard.full() {
		return 1
	} else {
		return 0
	}
}

func (board *Board) normalize(domino *domino) bool {
	x, y := domino.indicies()
	board.board[x] = true
	board.board[y] = true
	board.dominos = append(board.dominos, domino)
	for index, value := range board.board {
		if value == false {
			moves := make([]int, 4)
			for _, neighbour := range board.neighbours(index) {
				if board.board[neighbour] == false {
					moves = append(moves, neighbour)
				}
			}
			if len(moves) == 0 {
				fmt.Printf("Empty move")
				return false
			} else if len(moves) == 1 {
				board.dominos = append(board.dominos, NewDomino(index, moves[0], board))
				board.board[index] = true
				board.board[moves[0]] = true
			}
		}
	}
	return true
}

func (board *Board) Steps() []Step {
	dominos := board.moves() // dla danej kratownicy , znajdz mozliwe ruchy (nowe domina)
	boards := make([]Step, 0, len(dominos))
	for index := range dominos {
		// wykonaj mozliwe ruchy, tworzac nowe kratownice
		newboard := copyBoard(board)
		domino := dominos[index]
		if newboard.normalize(domino) {
			boards = append(boards, newboard)
		}
	}
	return boards
}

var runes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func (board *Board) Display() {
	rand.Seed(time.Now().UnixNano())
	display := make(map[int]rune)
	for _, domino := range board.dominos {
		randomLetter := rune(runes[rand.Intn(len(runes))])
		a, b := domino.indicies()
		display[a] = randomLetter
		display[b] = randomLetter
	}
	fmt.Printf("\n")
	for index := range board.board {
		if index%board.columns == 0 {
			fmt.Printf("\n")
		}
		el, ok := display[index]
		if !ok {
			el = ' '
		}
		fmt.Printf("%c", el)
	}
	fmt.Printf("\n")
}
