package game

type dir int

const (
	UD dir = 0
	LR dir = 1
)

type domino struct {
	index int
	dir   dir
	board *Board
}

func (domino *domino) indicies() (int, int) {
	if domino.dir == UD {
		return domino.index, domino.index + domino.board.columns
	}
	return domino.index, domino.index + 1
}

type Board struct {
	rows    int
	columns int
	board   []bool
}

func (board *Board) index(row int, column int) int {
	return row*board.columns + column
}

func (board *Board) valid(row int, column int) bool {
	return row < board.rows && column < board.columns
}

func (board *Board) Mark(row int, column int) {
	board.board[board.index(row, column)] = true
}

func EmptyBoard(rows int, columns int) *Board {
	return &Board{
		rows,
		columns,
		make([]bool, rows*columns, rows*columns),
	}
}

func copyBoard(base Board) *Board {
	board := EmptyBoard(base.rows, base.columns)
	for index := range base.board {
		board.board[index] = base.board[index]
	}
	return board
}

func (b *Board) moves() []*domino {
	dominos := []*domino{}
	for index := range b.board {
		if !b.board[index] {
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

func (board *Board) Hash() SolutionHash {
	letters := ""
	for integer := range board.board {
		if board.board[integer] {
			letters += "T"
		} else {
			letters += "F"
		}
	}
	return SolutionHash(letters)
}

func (board *Board) IsSolution() bool {
	for index := range board.board {
		if board.board[index] {
			continue
		} else {
			return false
		}
	}
	return true
}

func (board *Board) Steps() []Step {
	dominos := board.moves()
	boards := make([]Step, 0, len(dominos))
	for index := range dominos {
		newboard := copyBoard(*board)
		domino := dominos[index]
		x, y := domino.indicies()
		newboard.board[x] = true
		newboard.board[y] = true
		boards = append(boards, newboard)
	}
	return boards
}
