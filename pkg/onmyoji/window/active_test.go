package window

import (
	"testing"

	"github.com/go-vgo/robotgo"
)

func TestMove_Execute(t *testing.T) {
	// 查找进程 ID
	ids, err := robotgo.FindIds("Bluestacks")
	if err != nil {
		t.Fatalf("Failed to find process IDs: %v", err)
	}

	// 检查是否找到至少一个 ID
	if len(ids) == 0 {
		t.Fatalf("No process IDs found for 'Bluestacks'")
	}
	pid := robotgo.GetPid()
	t.Logf("Current process ID: %d", pid)

	// 激活第一个找到的进程
	err = robotgo.ActivePid(ids[1])
	if err != nil {
		t.Fatalf("Failed to activate process with PID %d: %v", ids[1], err)
	}

	t.Logf("Successfully activated process with PID %d", ids[1])
}
