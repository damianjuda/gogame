package game

import "fmt"

type SolutionHash string

type Step interface {
	Steps() []Step
	Hash() SolutionHash
	IsSolution() bool
}

type registry map[SolutionHash]Step
type status int

const (
	SUCCESS    status = 0
	FAILURE    status = 1
	NEXT_ROUND status = 2
)

type result struct {
	status status
	step   Step
}

func gameRound(frontier chan Step, visited *registry, tokens chan token, results chan<- Step) {
	for {
		token := <-tokens
		if len(tokens) == 0 && len(frontier) == 0 {
			results <- nil
		} else {
			fmt.Printf("Go %d\n", len(tokens))
			test := <-frontier
			tokens <- token
			(*visited)[test.Hash()] = test
			if test.IsSolution() {
				results <- test
				break
			}
			steps := test.Steps()
			for index := range steps {
				step := steps[index]
				if (*visited)[step.Hash()] != nil {
					//ignore
					fmt.Printf("Ignore\n")
				} else {
					frontier <- step
				}
			}
		}
	}
}

type token int

func Play(initial Step) Step {
	visited := &registry{}
	frontier := make(chan Step)
	results := make(chan Step)
	jobsNo := 5
	tokens := make(chan token, jobsNo)
	for index := 0; index < jobsNo; index++ {
		tokens <- token(index)
		go gameRound(frontier, visited, tokens, results)
	}
	frontier <- initial
	result := <-results
	return result
}
