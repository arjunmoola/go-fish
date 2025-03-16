package main

import (
    "bufio"
    "strings"
    "fmt"
)

type Cmd interface {
    String() string
    ProcessResponse(e *Engine) error
}

type uciCmd struct {}

func (u uciCmd) String() string {
    return "uci"
}

func (u uciCmd) ProcessResponse(e *Engine) error {
    scanner := bufio.NewScanner(e.out)
    results := strings.Builder{}
    for scanner.Scan() {
        text := scanner.Text()

        results.WriteString(text + "\n")

        if strings.HasPrefix(text, "uciok") {
            break
        }
    }
    e.uciResults = results.String()
    return nil
}

type uciIsReady struct {}

func (u uciIsReady) String() string {
    return "isready"
}

func (u uciIsReady) ProcessResponse(e *Engine) error {
    scanner := bufio.NewScanner(e.out)

    if scanner.Scan() {
        e.uciOk = scanner.Text()
    }

    return nil
}

type uciNewGame struct {}

func (u uciNewGame) String() string {
    return "ucinewgame"
}

func (u uciNewGame) ProcessResponse(e *Engine) error {
    return nil
}

type uciQuit struct {}

func (u uciQuit) String() string {
    return "quit"
}

func (u uciQuit) ProcessResponse(e *Engine) error {
    return nil
}

type uciGo struct {
    depth int
}

func (u uciGo) String() string {
    return fmt.Sprintf("go depth %d\n", u.depth)
}

func (u uciGo) ProcessResponse(e *Engine) error {
    scanner := bufio.NewScanner(e.out)
    results := strings.Builder{}

    for scanner.Scan() {
        text := scanner.Text()

        results.WriteString(text + "\n")

        if strings.HasPrefix(text, "bestmove") {
            break
        }
    }
    e.searchResults = results.String()
    return nil
}

type uciPosition struct {
    fen string
    moves []string
}

func (u uciPosition) String() string {
    builder := strings.Builder{}
    builder.WriteString("position ")
    builder.WriteString(fmt.Sprintf("position %s ", u.fen))
    builder.WriteString("moves ")
    for _, mv := range u.moves {
        builder.WriteString(mv + " ")
    }
    builder.WriteString("\n")
    return builder.String()
}

func (u uciPosition) ProcessResponse(e *Engine) error {
    return nil
}
