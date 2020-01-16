package game

import (
	"fmt"
	"github.com/Workiva/go-datastructures/queue"
	"runtime"
	"sync"
	"time"
)

type SolutionHash int // reprezentacja gwarantujaca unikalnosc kolejnych krokow, zeby sie nie zapetlic

type Step interface {
	queue.Item
	Steps() []Step    // kazdy krok generuje kolejne kroki
	Hash() string     // kazdy krok posiada unikalna reprezentacje
	IsSolution() bool // kazdy krok moze byc rozwiazaniem gry
	Reject() bool
}

func gameRound(index int, frontier *queue.PriorityQueue, visited *sync.Map, tokens chan token, results chan<- Step) {
	// w kolko szukaj rozwiazania
	counter := 0
	for {
		counter += 1
		token := <-tokens // dziwny mechanizm wykrywania, ze nie ma dalszych krokow do analizy i trzeba oglosic porazke - mozliwe ze niepoprawny
		if frontier.Empty() {
			if len(tokens) == 0 {
				fmt.Printf("Stack (id=%d) (token=%d) (todo=%d) (workers=%d) \n", index, token, frontier.Len(), len(tokens))
				results <- nil
			} else {
				fmt.Printf("Sleep (id=%d) (token=%d) (todo=%d) (workers=%d) \n", index, token, frontier.Len(), len(tokens))
				time.Sleep(1000 * time.Millisecond)
				tokens <- token // oddaj token
			}
		} else {
			tokens <- token               // oddaj token
			items, err := frontier.Get(1) // wez krok do analizy
			if err != nil {
				panic(err)
			}
			for _, item := range items {
				test := item.(Step)
				fmt.Printf("Hash %s\n", test.Hash())
				(*visited).Store(test.Hash(), true) // oznacz krok jako juz analizowany
				if test.IsSolution() {
					results <- test // jesli krok jest rozwiazaniem zwroc fo jako sukces
				} else if !test.Reject() {
					// jesli aktualny krok nie jest rozwiazaniem to zobacz mozliwe do wykonania kroki
					if counter%1000 == 0 {
						fmt.Printf("Step (id=%d) (token=%d) (todo=%d) (workers=%d) (counter=%d)\n", index, token, frontier.Len(), len(tokens), counter)
					}
					steps := test.Steps()
					for stepIndex := range steps {
						step := steps[stepIndex]
						hash := step.Hash()
						_, ok := (*visited).Load(hash)
						if ok {
							// zignoruj juz znane kroki
						} else {
							(*visited).Store(step.Hash(), true) // oznacz krok jako juz analizowany
							frontier.Put(step)                  // wrzuc nowy krok do analizy jesli jest nowy
						}
					}
				}
			}
		}
	}
}

type token int

func Play(initial Step) Step {
	visited := &sync.Map{}                            // mapa krokow, juz minionych (wystarczyl by zbior hashy, bo teraz zuzywa za duzo pamieci)
	frontier := queue.NewPriorityQueue(1000000, true) // synchronizowany kanał kroków do analizy
	results := make(chan Step)                        // kanał na wyniki (de facto pierwszy)
	jobsNo := runtime.NumCPU()
	fmt.Printf("Jobs no %d\n", jobsNo)
	runtime.GOMAXPROCS(jobsNo)
	tokens := make(chan token, jobsNo) // kanal pomocniczy do wykrywania, ze juz nie ma krokow do analizy
	for index := 0; index < jobsNo; index++ {
		tokens <- token(index)
		go gameRound(index, frontier, visited, tokens, results) // odpal watki
	}
	frontier.Put(initial) // wrzuc stan poczatkowy
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
