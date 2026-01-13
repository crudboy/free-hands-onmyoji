//go:build windows
// +build windows

package window

import (
	"fmt"

	"github.com/go-vgo/robotgo"
	"github.com/lxn/win"
)

// GetWindowPosition Windows 实现：获取特定应用窗口的位置和大小信息
func (w *WindowsPlatform) GetWindowPosition() (Window, error) {
	displayId := 0
	if DisplayID != -1 {
		displayId = DisplayID
		return w.GetWindowPositionOnSecondDisplay(displayId)
	}
	return w.getWindowPositionByMaster()
}

// getWindowPositionByMaster Windows 实现：获取主屏幕窗口位置
func (w *WindowsPlatform) getWindowPositionByMaster() (Window, error) {
	fpid, err := robotgo.FindIds(w.appName)
	if err != nil {
		return Window{}, fmt.Errorf("查找应用 '%s' 失败: %v", w.appName, err)
	}

	if len(fpid) == 0 {
		return Window{}, fmt.Errorf("未找到运行中的应用: %s", w.appName)
	}
	x, y, width, height := robotgo.GetBounds(fpid[0])
	// 3. 验证返回值的有效性（未找到窗口时可能返回 0 或错误坐标）
	if width == 0 && height == 0 {
		return Window{}, fmt.Errorf("应用 '%s' 正在运行，但无法获取其窗口位置", w.appName)
	}
	return Window{
		WindowX: x,
		WindowY: y,
		WindowW: width,
		WindowH: height,
		CaptureArea: CaptureArea{
			X: x,
			Y: y,
			W: width,
			H: height,
		},
	}, nil
}

// GetWindowPositionOnSecondDisplay Windows 实现：获取第二显示器上的窗口位置
func (w *WindowsPlatform) GetWindowPositionOnSecondDisplay(displayId int) (Window, error) {
	// TODO: Windows 平台实现
	return Window{}, fmt.Errorf("Windows 平台暂未实现多显示器窗口位置检测")
}

func (w *WindowsPlatform) CloseApplication() error {
	fpid, err := robotgo.FindIds(w.appName)
	if err != nil {
		return fmt.Errorf("查找应用 '%s' 失败: %v", w.appName, err)
	}

	if len(fpid) == 0 {
		return fmt.Errorf("未找到运行中的应用: %s", w.appName)
	}
	robotgo.Kill(fpid[0])
	return nil
}

// SetScaleWindow Windows 实现：设置窗口大小
func (w *WindowsPlatform) SetScaleWindow(width, height int) {
	fpid, err := robotgo.FindIds(w.appName)
	if err != nil {
		fmt.Printf("查找应用 '%s' 失败: %v\n", w.appName, err)
		return
	}

	if len(fpid) == 0 {
		fmt.Printf("未找到运行中的应用: %s\n", w.appName)
		return
	}
	win.MoveWindow(win.HWND(fpid[0]), 100, 100, 800, 600, true)
}
