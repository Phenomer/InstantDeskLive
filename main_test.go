// filepath: main_test.go
package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

// テスト用の一時的なJSONファイルを作成
func createTempConfig(t *testing.T, processes []ProcessConfig) string {
	t.Helper()
	tmpfile, err := os.CreateTemp("", "process*.json")
	if err != nil {
		t.Fatalf("テスト用の一時ファイル作成に失敗: %v", err)
	}

	data, err := json.Marshal(processes)
	if err != nil {
		t.Fatalf("JSONエンコードに失敗: %v", err)
	}

	if _, err := tmpfile.Write(data); err != nil {
		t.Fatalf("一時ファイルへの書き込みに失敗: %v", err)
	}

	if err := tmpfile.Close(); err != nil {
		t.Fatalf("一時ファイルのクローズに失敗: %v", err)
	}

	return tmpfile.Name()
}

// loadConfig のテスト
func TestLoadConfig(t *testing.T) {
	// 正常系のテスト
	validProcs := []ProcessConfig{
		{Name: "proc1", Command: "echo", Args: []string{"hello"}},
		{Name: "proc2", Command: "echo", Args: []string{"world"}},
	}

	configPath := createTempConfig(t, validProcs)
	defer func() {
		if err := os.Remove(configPath); err != nil {
			t.Errorf("一時ファイルの削除に失敗: %v", err)
		}
	}()

	procs, err := loadConfig(configPath)
	if err != nil {
		t.Fatalf("ロードに失敗: %v", err)
	}

	if len(procs) != 2 {
		t.Errorf("ロードされたプロセス数が想定と異なる: got %d, want 2", len(procs))
	}

	// ファイルが見つからない場合のテスト
	_, err = loadConfig("non_existent_file.json")
	if err == nil {
		t.Error("存在しないファイルでエラーが発生していない")
	}

	// 不正なJSONのテスト
	invalidJSON := filepath.Join(t.TempDir(), "invalid.json")
	if err := os.WriteFile(invalidJSON, []byte("invalid json"), 0644); err != nil {
		t.Fatalf("不正なJSONファイルの作成に失敗: %v", err)
	}

	_, err = loadConfig(invalidJSON)
	if err == nil {
		t.Error("不正なJSONでエラーが発生していない")
	}
}

// startProcesses と terminateProcesses のテスト
func TestStartAndTerminateProcesses(t *testing.T) {
	// スリープコマンドを実行する設定
	var sleepCmd string
	var sleepArgs []string

	if runtime.GOOS == "windows" {
		sleepCmd = "timeout"
		sleepArgs = []string{"/T", "10", "/NOBREAK"}
	} else {
		sleepCmd = "sleep"
		sleepArgs = []string{"10"}
	}

	processes := []ProcessConfig{
		{Name: "sleeper1", Command: sleepCmd, Args: sleepArgs},
		{Name: "sleeper2", Command: sleepCmd, Args: sleepArgs},
	}

	// 短い時間で実行終了するコマンドをチェック
	procs, err := startProcesses(processes)
	if err != nil {
		t.Fatalf("プロセス起動に失敗: %v", err)
	}

	// すべてのプロセスが起動しているか確認
	if len(procs) != 2 {
		t.Errorf("起動したプロセス数が想定と異なる: got %d, want 2", len(procs))
	}

	// プロセスが実行中であることを確認
	for _, p := range procs {
		if p.cmd.Process == nil {
			t.Errorf("プロセス %s が正しく起動していない", p.name)
		}
	}

	// プロセスを終了
	terminateProcesses(procs)
	time.Sleep(500 * time.Millisecond)
}

// monitorProcesses のテスト
func TestMonitorProcesses(t *testing.T) {
	// 短時間で終了するコマンドを使用
	var echoCmd string
	var echoArgs []string

	if runtime.GOOS == "windows" {
		echoCmd = "cmd"
		echoArgs = []string{"/C", "echo", "test"}
	} else {
		echoCmd = "echo"
		echoArgs = []string{"test"}
	}

	cmd := exec.Command(echoCmd, echoArgs...)
	if err := cmd.Start(); err != nil {
		t.Fatalf("テスト用コマンドの起動に失敗: %v", err)
	}

	procs := []procInfo{
		{name: "echo-test", cmd: cmd},
	}

	done := make(chan string, len(procs))
	monitorProcesses(procs, done)

	// プロセスが終了し、通知が来ることを確認
	select {
	case procName := <-done:
		if procName != "echo-test" {
			t.Errorf("モニター対象外のプロセス名を受信: got %s, want echo-test", procName)
		}
	case <-time.After(3 * time.Second):
		t.Error("プロセス終了通知のタイムアウト")
	}
}

// 不正なコマンドで startProcesses が失敗するテスト
func TestStartProcessesWithInvalidCommand(t *testing.T) {
	processes := []ProcessConfig{
		{Name: "invalid", Command: "non-existent-command", Args: []string{}},
	}

	_, err := startProcesses(processes)
	if err == nil {
		t.Error("存在しないコマンドでエラーが発生していない")
	}
}
