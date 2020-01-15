package game

import (
	"sync"
	"fmt"
)

type SolutionHash int // reprezentacja gwarantujaca unikalnosc kolejnych krokow, zeby sie nie zapetlic

type Step interface {
	Steps() []Step // kazdy krok generuje kolejne kroki
	Hash() int // kazdy krok posiada unikalna reprezentacje
	IsSolution() bool // kazdy krok moze byc rozwiazaniem gry
}

func gameRound(index int, frontier chan Step, visited *sync.Map, tokens chan token, results chan<- Step) {
	// w kolko szukaj rozwiazania
	for {
		token := <-tokens // dziwny mechanizm wykrywania, ze nie ma dalszych krokow do analizy i trzeba oglosic porazke - mozliwe ze niepoprawny
		if len(tokens) == 0 && len(frontier) == 0 {
			fmt.Printf("Stack (id=%d) (token=%d) (todo=%d) (workers=%d) \n", index, token, len(frontier), len(tokens))
			tokens <- token // oddaj token
			results <- nil
		} else {
			test := <-frontier // wez krok do analizy
			tokens <- token // oddaj token
			(*visited).Store(test.Hash(), true)// oznacz krok jako juz analizowany
			if test.IsSolution() {
				results <- test // jesli krok jest rozwiazaniem zwroc fo jako sukces
				break
			}
			// jesli aktualny krok nie jest rozwiazaniem to zobacz mozliwe do wykonania kroki
			fmt.Printf("Step (id=%d) (token=%d) (todo=%d) (workers=%d) \n", index, token, len(frontier), len(tokens))
			steps := test.Steps()
			for stepIndex := range steps {
				step := steps[stepIndex]
				hash := step.Hash()
				_, ok := (*visited).Load(hash)
				if ok {
					// zignoruj juz znane kroki
				} else {
					(*visited).Store(step.Hash(), true)// oznacz krok jako juz analizowany
					frontier <- step // wrzuc nowy krok do analizy jesli jest nowy
				}
			}
		}
	}
}

type token int

func Play(initial Step) Step {
	visited := &sync.Map{} // mapa krokow, juz minionych (wystarczyl by zbior hashy, bo teraz zuzywa za duzo pamieci)
	frontier := make(chan Step, 1000000) // synchronizowany kanał kroków do analizy
	results := make(chan Step) // kanał na wyniki (de facto pierwszy)
	jobsNo := 5 // liczba 'watkow'
	tokens := make(chan token, jobsNo) // kanal pomocniczy do wykrywania, ze juz nie ma krokow do analizy
	for index := 0; index < jobsNo; index++ {
		tokens <- token(index)
		go gameRound(index, frontier, visited, tokens, results) // odpal watki
	}
	frontier <- initial // wrzuc stan poczatkowy
	for index := 0; index < jobsNo; index++ {
		result := <-results // czekaj na wynik (krok - success, nil - nie ma rozwiazania)
		if result != nil {
			fmt.Printf("Success\n")
			return result
		}
	}
	fmt.Printf("Failure\n")
	return nil
}
