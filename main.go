package main

import (
    //"github.com/joho/godotenv"
    "log"
    "fmt"
    "os"
    "bufio"
)

func parseInput(input string) Cmd {
    var cmd Cmd
    if input == "uci" {
        cmd = uciCmd{}
    } else if input == "isready" {
        cmd = uciIsReady{}
    } else if input == "quit" {
        cmd = uciQuit{}
    } else if input == "go" {
        cmd = uciGo{
            depth: 3,
        }
    }

    return cmd
}

func main() {
    engine, err := NewEngine("stockfish")
    
    if err != nil {
        log.Fatal(err)
    }

    scanner := bufio.NewScanner(os.Stdin)

    var quit bool

    for scanner.Scan() {
        input := scanner.Text()

        cmd := parseInput(input)

        if cmd == nil {
            continue
        }

        if err := engine.Run(cmd); err != nil {
            log.Panic(err)
        }

       var output string

        switch cmd.(type) {
        case uciCmd:
            output = engine.UciResults()
        case uciIsReady:
            output = engine.UciOk()
        case uciGo:
            output = engine.SearchResults()
        case uciQuit:
            output = "engine has quit"
            quit = true
        default:
            continue
        }

        fmt.Println(output)

        if quit {
            break
        }
    }
}

