package main

import (
	"fmt"
	"os"

	"github.com/wbhildeb/nytimesgames/wordle/hint"
)

func main() {
	targets := map[string]string{
		"natea":  "Suck  my  nuts,  Winnipegger",
		"amir":   "Suck  my  tiny  nuts,  prince  ali",
		"rafid":  "I  miss  working  with  ya",
		"robert": "I \nhave \ntoo \nmuch \ntime \non \nmy \nhands",
	}

	_ = targets

	if len(os.Args) < 3 {
		fmt.Println("Usage: wordle <target> <guess>")
		fmt.Println("Available targets:")
		for name := range targets {
			fmt.Printf("  - %s\n", name)
		}
		os.Exit(1)
	}

	targetName := os.Args[1]
	guess := os.Args[2]

	target, ok := targets[targetName]
	if !ok {
		fmt.Printf("Unknown target: %s\n", targetName)
		fmt.Println("Available targets:")
		for name := range targets {
			fmt.Printf("  - %s\n", name)
		}
		os.Exit(1)
	}

	fmt.Println(
		hint.FormattedHint(guess, target),
	)
}
