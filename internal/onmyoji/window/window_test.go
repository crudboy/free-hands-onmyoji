//go:build windows
// +build windows

package window

import (
	"free-hands-onmyoji/internal/logger"
	"free-hands-onmyoji/internal/utils"
	"image/png"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-vgo/robotgo"
)

func TestWindow(t *testing.T) {
	logger.Init()
	err := utils.EnableProcessDPIAware()
	if err != nil {
		t.Fatalf("EnableProcessDPIAware failed: %v", err)
	}
	hwnd, err := utils.FindHwndByTitle("阴阳师-MuMu模拟器专版")
	if err != nil {
		t.Fatalf("FindHwndByTitle failed: %v", err)
	}
	err = utils.ActivateWindow(hwnd)
	time.Sleep(2 * time.Second) // 等待窗口激活
	if err != nil {
		t.Fatalf("ActivateWindow failed: %v", err)
	}
	t.Logf("Found window handle: %v", hwnd)
	x, y, w, h, err := utils.GetWindowBounds(hwnd)
	if err != nil {
		t.Fatalf("GetWindowBounds failed: %v", err)
	}
	t.Logf("Window bounds: x=%d, y=%d, width=%d, height=%d", x, y, w, h)

	// 使用 robotgo 截图
	img, err := robotgo.CaptureImg(x, y, w, h)
	if err != nil {
		t.Fatalf("CaptureImg failed: %v", err)
	}

	// 创建 capture 目录
	if err := os.MkdirAll("capture", 0755); err != nil {
		t.Fatalf("创建capture目录失败: %v", err)
	}

	// 生成文件名，包含时间戳
	filename := filepath.Join("capture", "window_"+time.Now().Format("20060102_150405")+".png")
	file, err := os.Create(filename)
	if err != nil {
		t.Fatalf("创建文件失败: %v", err)
	}
	defer file.Close()

	// 编码并保存为 PNG
	if err := png.Encode(file, img); err != nil {
		t.Fatalf("保存图片失败: %v", err)
	}

	t.Logf("窗口截图已保存到: %s", filename)
}
