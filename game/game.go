package game

import "fmt"

type SolutionHash string // reprezentacja gwarantujaca unikalnosc kolejnych krokow, zeby sie nie zapetlic


type Step interface {
	Steps() []Step // kazdy krok generuje kolejne kroki
	Hash() SolutionHash // kazdy krok posiada unikalna reprezentacje
	IsSolution() bool // kazdy krok moze byc rozwiazaniem gry
}

type registry map[SolutionHash]Step // mapa krokow, nie jest synchronizowana, ale z racji ze nie mozna usuwac elementow to wydaje sie to bezpieczne
func gameRound(frontier chan Step, visited *registry, tokens chan token, results chan<- Step) {
	// w kolko szukaj rozwiazania
	for {
		token := <-tokens // dziwny mechanizm wykrywania, ze nie ma dalszych krokow do analizy i trzeba oglosic porazke - mozliwe ze niepoprawny
		if len(tokens) == 0 && len(frontier) == 0 {
			results <- nil
		} else {
			test := <-frontier // wez krok do analizy
			tokens <- token // oddaj token
			(*visited)[test.Hash()] = test // oznacz krok jako juz analizowany
			if test.IsSolution() {
				results <- test // jesli krok jest rozwiazaniem zwroc fo jako sukces
				break
			}
			// jesli aktualny krok nie jest rozwiazaniem to zobacz mozliwe do wykonania kroki
			steps := test.Steps()
			for index := range steps {
				step := steps[index]
				if (*visited)[step.Hash()] != nil {
					// zignoruj juz znane kroki
				} else {
					frontier <- step // wrzuc nowy krok do analizy jesli jest nowy
				}
			}
		}
	}
}

type token int

func Play(initial Step) Step {
	visited := &registry{} // mapa krokow, juz minionych (wystarczyl by zbior hashy, bo teraz zuzywa za duzo pamieci)
	frontier := make(chan Step) // synchronizowany kanał kroków do analizy
	results := make(chan Step) // kanał na wyniki (de facto pierwszy)
	jobsNo := 5 // liczba 'watkow'
	tokens := make(chan token, jobsNo) // kanal pomocniczy do wykrywania, ze juz nie ma krokow do analizy
	for index := 0; index < jobsNo; index++ {
		tokens <- token(index)
		go gameRound(frontier, visited, tokens, results) // odpal watki
	}
	frontier <- initial // wrzuc stan poczatkowy
	result := <-results // czekaj na wynik (krok - success, nil - nie ma rozwiazania)
	return result
}
