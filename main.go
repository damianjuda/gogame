package main

import (
	"fmt"
	"github.com/damianjuda/gogame/game/game"
	"gopkg.in/alecthomas/kingpin.v2"
	"strconv"
	"strings"
)

var (
	rows    = kingpin.Arg("rows", "Board rows no").Required().Int()
	columns = kingpin.Arg("columns", "Board columns no").Required().Int()
	split   = kingpin.Arg("split", "Split by").String()
)

func main() {
	play()
}

func play() {
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
	game.Play(b)
}
