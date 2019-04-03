package main

import (
	"fmt"
	"strings"

	"github.com/chzyer/readline"
)

func main() {
	l, err := readline.NewEx(&readline.Config{
		Prompt:              ">>",
		HistoryFile:         "read.tmp",
		AutoComplete:        nil,
		InterruptPrompt:     "^C",
		EOFPrompt:           "exit",
		HistorySearchFold:   false,
		FuncFilterInputRune: nil,
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()
	read := true
	for read {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		}
		line = strings.TrimSpace(line)
		fmt.Println("%[", line, "]")
		switch {
		case strings.HasPrefix(line, "exit"):
			fmt.Println("Bye")
			read = false
			break
		}
	}
}
