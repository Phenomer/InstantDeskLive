package main

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "os"
    "os/exec"
    "os/signal"
    "syscall"
    "time"
)

type ProcessConfig struct {
    Name    string   `json:"name"`
    Command string   `json:"command"`
    Args    []string `json:"args"`
}

type procInfo struct {
    name string
    cmd  *exec.Cmd
}

func main() {
    processes, err := loadConfig("process.json")
    if err != nil {
        fmt.Println(err)
        return
    }

    if len(processes) < 2 {
        fmt.Println("サブプロセスは2つ以上必要です。")
        return
    }

    procs, err := startProcesses(processes)
    if err != nil {
        fmt.Println(err)
        return
    }

    done := make(chan string, len(procs))
    monitorProcesses(procs, done)

    // signal.NotifyContext を使い、キャンセル可能なコンテキストを作成
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer stop()

    select {
    case proc := <-done:
        fmt.Printf("%s が終了しました。全てのサブプロセスを終了します。\n", proc)
    case <-ctx.Done():
        fmt.Printf("シグナルを受信しました [%v]。全てのサブプロセスを終了します。\n", ctx.Err())
    }

    terminateProcesses(procs)

    // サブプロセス終了確認のための待機
    time.Sleep(10 * time.Second)
}

func loadConfig(filename string) ([]ProcessConfig, error) {
    data, err := os.ReadFile(filename)
    if err != nil {
        return nil, fmt.Errorf("process.json の読み込みに失敗しました: %w", err)
    }
    var processes []ProcessConfig
    if err := json.Unmarshal(data, &processes); err != nil {
        return nil, fmt.Errorf("JSON の解析に失敗しました: %w", err)
    }
    return processes, nil
}

func startProcesses(processes []ProcessConfig) ([]procInfo, error) {
    var procs []procInfo
    for _, p := range processes {
        cmd := exec.Command(p.Command, p.Args...)
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr

        if err := cmd.Start(); err != nil {
            return nil, fmt.Errorf("%s の起動に失敗しました: %w", p.Name, err)
        }
        fmt.Printf("%s を起動しました (PID: %d)\n", p.Name, cmd.Process.Pid)
        procs = append(procs, procInfo{name: p.Name, cmd: cmd})
    }
    return procs, nil
}

func monitorProcesses(procs []procInfo, done chan<- string) {
    for _, p := range procs {
        go func(pi procInfo) {
            if err := pi.cmd.Wait(); err != nil && !errors.Is(err, exec.ErrNotFound) {
                fmt.Printf("%s 終了時にエラー: %v\n", pi.name, err)
            }
            done <- pi.name
        }(p)
    }
}

func terminateProcesses(procs []procInfo) {
    for _, p := range procs {
        killProcess(p.cmd, p.name)
    }
}

func killProcess(cmd *exec.Cmd, name string) {
    if cmd.Process != nil {
        if err := cmd.Process.Kill(); err != nil {
            fmt.Printf("%s (PID: %d) の終了に失敗しました: %v\n", name, cmd.Process.Pid, err)
        } else {
            fmt.Printf("%s (PID: %d) を終了しました。\n", name, cmd.Process.Pid)
        }
    }
}
