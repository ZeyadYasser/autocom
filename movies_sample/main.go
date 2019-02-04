package main

import (
	"fmt"
	"os"
	"bufio"

	"github.com/zeyadyasser/autocom/engine/skip"
)

func setMovies(E *skip.SkipEngine) {
	file, err := os.Open("movies_sample/titles.txt")
    if err != nil {
        panic(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
		title := scanner.Text()
		E.Set(title, nil)
    }
    if err := scanner.Err(); err != nil {
        panic(err)
    }
}

func main() {
	opts := skip.Options{
		MaxLevels: 5,
		ToLower: true,
		SkipBegin: true,
	}
	E := skip.NewSkipEngine(opts, nil)

	setMovies(E)
	E.Remove("Mad Max: Fury Road")

	scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
		text := scanner.Text()
		top, _ := E.TopN(text, 4)
		for k := range top {
			fmt.Println(k)
		}
    }
}
