package main

import (
    "io"
    "os/exec"
    "sync"
    "fmt"
)

type Engine struct {
    mu sync.RWMutex
    cmd *exec.Cmd
    in *io.PipeWriter
    out *io.PipeReader
    uciResults string
    uciOk string
    searchResults string
}

func NewEngine(p string) (*Engine, error) {
    path, err := exec.LookPath(p)

    if err != nil {
        return nil, err
    }

    rin, win := io.Pipe()
    rOut, wOut := io.Pipe()

    cmd := exec.Command(path)
    cmd.Stdin = rin
    cmd.Stdout = wOut

    engine := &Engine{
        cmd: cmd,
        in: win,
        out: rOut,
    }

    _ = engine.cmd.Start()

    return engine, nil
}

func (e *Engine) Close() error {
    if err := e.in.Close(); err != nil {
        return err
    }
    if err := e.out.Close(); err != nil {
        return err
    }
    return e.cmd.Process.Kill()
}

func (e *Engine) processCommandLocked(cmd Cmd) error {
    e.mu.Lock()
    defer e.mu.Unlock()
    return e.processCommand(cmd)
}

func (e *Engine) processCommand(cmd Cmd) error {
    if _, err := fmt.Fprintln(e.in, cmd.String()); err != nil {
        return err
    }

    if err := cmd.ProcessResponse(e); err != nil {
        return err
    }

    return nil
}

func (e *Engine) Run(cmd Cmd) error {
    switch cmd.(type) {
    case uciQuit:
        return e.Close()
    default:
        return e.processCommandLocked(cmd)
    }
}

func (e *Engine) UciResults() string {
    e.mu.RLock()
    defer e.mu.RUnlock()
    return e.uciResults
}

func (e *Engine) UciOk() string {
    e.mu.RLock()
    defer e.mu.RUnlock()
    return e.uciOk
}

func (e *Engine) SearchResults() string {
    e.mu.RLock()
    defer e.mu.RUnlock()
    return e.searchResults
}
