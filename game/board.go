package game

type dir int

const (
	UD dir = 0 // pionowe domino
	LR dir = 1 // poziome domino
)

type domino struct {
	index int // indeks na 1 wymiarowej tablicy (reprezentującej 2 wymiarową) prawy górny róg domino
	dir   dir // poziome albo pionowe domino
	board *Board // referencja do całej tablicy
}

func (domino *domino) indicies() (int, int) {
	if domino.dir == UD {
		return domino.index, domino.index + domino.board.columns
	}
	return domino.index, domino.index + 1
}

type Board struct {
	rows    int // liczba wierszy kratownicy
	columns int // liczba kolumn kratownicy
	board   []bool // tablica reprezentująca puste i pełne miejsca na kratownicy
}

func (board *Board) index(row int, column int) int {
	return row*board.columns + column // indeks w 1 wymiarowej tablicy emulującej 2 wymiarową
}

func (board *Board) valid(row int, column int) bool {
	return row < board.rows && column < board.columns // czy pozycja (wiersz,kolumna) nalezy do kratownicy
}

func (board *Board) Mark(row int, column int) {
	board.board[board.index(row, column)] = true      // zaznacz pole jako zajęte
}

func EmptyBoard(rows int, columns int) *Board {
	// pusta kratonica
	return &Board{
		rows,
		columns,
		make([]bool, rows*columns, rows*columns),
	}
}

func copyBoard(base Board) *Board {
	// skopiuj kratownice
	board := EmptyBoard(base.rows, base.columns)
	for index := range base.board {
		board.board[index] = base.board[index]
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

func (board *Board) Hash() SolutionHash {
	// Niezbyt madra funkcja generujaca unikalnego hasha (stringa) dla kazdych 2 roznych kratownic
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
	// czy wszystkie pola na kratownicy sa zajete
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
	dominos := board.moves() // dla danej kratownicy , znajdz mozliwe ruchy (nowe domina)
	boards := make([]Step, 0, len(dominos))
	for index := range dominos {
		// wykonaj mozliwe ruchy, tworzac nowe kratownice
		newboard := copyBoard(*board) 
		domino := dominos[index]
		x, y := domino.indicies()
		newboard.board[x] = true
		newboard.board[y] = true
		boards = append(boards, newboard)
	}
	return boards
}
