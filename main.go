package main

import (
    "github.com/joho/godotenv"
    "log"
    "fmt"
    "os"
    "os/exec"
    "bytes"
    "bufio"
    "io"
    //"strings"
)

type uciGameCommand struct {
    position string
    depth int
}

type Config struct {
    binPath string
}


func newConfig() (*Config, error) {
    binPath := os.Getenv("STOCKFISH")

    if binPath == "" {
        return nil, fmt.Errorf("the the path for stockfish binary is not set")
    }

    config := &Config{
        binPath: binPath,
    }

    return config, nil
}

func processUCICommand(scanner *bufio.Scanner) string {
    if !scanner.Scan() {
        return ""
    }

    return scanner.Text()
}

func readInputs(r io.Reader) chan string {
    inputs := make(chan string)

    go func() {
        defer close(inputs)

        reader := bufio.NewReader(r)

        for {
            line, err := reader.ReadBytes('\n')

            if err != nil {
                panic(err)
            }

            inputs <- string(line)
        }
    }()

    return inputs
}

func readInputs2(pipe io.WriteCloser) {
    reader := bufio.NewReader(os.Stdin)

    go func() {
        defer pipe.Close()
        for {
            line, err := reader.ReadBytes('\n')

            if err != nil {
                panic(err)
            }

            fmt.Println(line)

            pipe.Write(line)

        }
    }()
}

func runStockfish() error {
    cmd := exec.Command("stockfish")

    r, err := cmd.StdoutPipe()

    if err != nil {
        return err
    }
    
    pipe, err := cmd.StdinPipe()

    if err != nil {
        return err
    }

    reader := bufio.NewReader(r)

    if err := cmd.Start(); err != nil {
        return err
    }

    //inputs := readInputs(os.Stdin)

    go readInputs2(pipe)

    msg, _ := reader.ReadBytes('\n')

    fmt.Println(string(bytes.TrimRight(msg, "\n")))

    //builder := strings.Builder{}

    for {
        msg, err := reader.ReadBytes('\n')

        if err != nil {
            fmt.Println(err)
            break
        }

        fmt.Println(string(bytes.TrimRight(msg, "\n")))

    }

    if err := cmd.Wait(); err != nil {
        return err
    }

    return nil
}

func main() {
    if err := godotenv.Load(); err != nil {
        log.Fatal(err)
    }
    
    _, err := newConfig()

    if err != nil {
        log.Fatal(err)
    }

    if err := runStockfish(); err != nil {
        log.Panic(err)
    }

}
