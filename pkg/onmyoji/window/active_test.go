package window

import (
	"free-hands-onmyoji/pkg/logger"
	"free-hands-onmyoji/pkg/utils"
	"testing"
	"time"

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
	err = robotgo.ActivePid(ids[0])
	if err != nil {
		t.Fatalf("Failed to activate process with PID %d: %v", ids[0], err)
	}

	t.Logf("Successfully activated process with PID %d", ids[0])
}
func TestActiveWindow(t *testing.T) {
	// AlertNotify("测试", "这是一个测试通知")
	robotgo.DisplayID = 1 // 设置显示器ID为1
	img, err := robotgo.Capture(2146, 447, 873, 523)
	if err != nil {
		t.Fatalf("Failed to capture screen: %v", err)
	}
	bounds := img.Bounds()
	t.Logf("Captured image bounds: %v", bounds)
}
func TestCaptureScreen(t *testing.T) {
	logger.Init()
	entity, err := GetWindowPositionOnSecondDisplay("BlueStacks", 1)
	if err != nil {
		t.Fatalf("获取窗口位置失败: %v", err)
	}
	time.Sleep(2 * time.Second) // 等待窗口稳定
	capture := entity.CaptureArea
	t.Logf("窗口位置: %+v", entity)
	robotgo.DisplayID = 1 // 设置显示器ID为1
	bit := robotgo.CaptureScreen(capture.X, capture.Y, capture.W, capture.H)
	img := robotgo.ToImage(bit)
	utils.SaveImg(img, "test_capture_screen.png")

}
