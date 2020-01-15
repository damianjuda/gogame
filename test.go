package main

import (
	"fmt"
	"strconv"
	"strings"
	"github.com/damianjuda/gogame/game"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	rows    = kingpin.Arg("rows", "Board rows no").Required().Int()
	columns = kingpin.Arg("columns", "Board columns no").Required().Int()
	split   = kingpin.Arg("split", "Split by").String()
)

func main() {
	kingpin.Parse()
	fmt.Printf("Rows %d Columns %d Split %s \n", *rows, *columns, *split)
	b := game.EmptyBoard(*rows, *columns)
	for _, sub := range strings.Split(*split, ".") {
		phrase := strings.Split(sub, ",")
		if len(phrase) > 1 {
			markRow, errRow := strconv.Atoi(phrase[0])
			markColumn, errColumn := strconv.Atoi(phrase[1])
			if errRow != nil {
				panic(errRow)
			} else if errColumn != nil {
				panic(errColumn)
			} else {
				fmt.Printf("Mark %d,%d\n", markRow, markColumn)
				b.Mark(markRow, markColumn)
			}
		}
	}
	result := game.Play(b)
	fmt.Printf("Result %f\n", result)
}
