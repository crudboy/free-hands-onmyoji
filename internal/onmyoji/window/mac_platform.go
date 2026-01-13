//go:build darwin
// +build darwin

package window

import (
	"fmt"
	"free-hands-onmyoji/internal/logger"
	"time"

	"github.com/andybrewer/mack"
	"github.com/go-vgo/robotgo"
)

// GetWindowPosition macOS 实现：获取特定应用窗口的位置和大小信息
func (m *MacPlatform) GetWindowPosition() (Window, error) {
	displayId := 0
	if DisplayID != -1 {
		displayId = DisplayID
		return m.GetWindowPositionOnSecondDisplay(displayId)
	}
	return m.getWindowPositionByMaster()
}

// getWindowPositionByMaster macOS 实现：获取主屏幕窗口位置
func (m *MacPlatform) getWindowPositionByMaster() (Window, error) {
	script := fmt.Sprintf(`
tell application "System Events"
 tell application process "%s"
  get {position, size} of window 1
 end tell
end tell
`, m.appName)

	tell, err := mack.Tell("System Events", script)
	if err != nil {
		return Window{}, fmt.Errorf("无法获取窗口信息: %v", err)
	}

	var x, y, w, h int
	_, err = fmt.Sscanf(tell, "%d, %d, %d, %d", &x, &y, &w, &h)
	if err != nil {
		return Window{}, fmt.Errorf("解析窗口位置信息失败: %v", err)
	}

	logger.Info("获取窗口位置和大小: 位置(%d,%d), 大小(%d,%d)", x, y, w, h)

	return Window{
		WindowX: x,
		WindowY: y,
		WindowW: w,
		WindowH: h,
		CaptureArea: CaptureArea{
			X: x,
			Y: y,
			W: w,
			H: h,
		},
	}, nil
}

// GetWindowPositionOnSecondDisplay macOS 实现：获取第二显示器上的窗口位置
func (m *MacPlatform) GetWindowPositionOnSecondDisplay(displayId int) (Window, error) {
	// 首先获取窗口在主屏幕上的位置
	mainScreenPosition, err := m.getWindowPositionByMaster()
	logger.Debug("获取主屏幕窗口位置: %+v", mainScreenPosition)
	if err != nil {
		return Window{}, fmt.Errorf("获取主屏幕窗口位置失败: %v", err)
	}

	// 获取所有显示器的信息
	displays, err := getDisplaysInfo()
	if err != nil {
		return Window{}, fmt.Errorf("获取显示器信息失败: %v", err)
	}

	if len(displays) < 2 {
		return Window{}, fmt.Errorf("未检测到第二显示器")
	}

	// 计算相对于第二显示器的位置
	// 找到窗口所在的显示器
	secondDisplay := displays[displayId] // 默认使用第二个显示器

	// 检查窗口是否在第二个显示器上
	// 如果窗口X坐标大于等于第二个显示器的X起始坐标，说明窗口在第二个显示器上
	if mainScreenPosition.WindowX >= secondDisplay.X &&
		mainScreenPosition.WindowY >= secondDisplay.Y &&
		mainScreenPosition.WindowX < secondDisplay.X+secondDisplay.W &&
		mainScreenPosition.WindowY < secondDisplay.Y+secondDisplay.H {
		// 窗口已经在第二个显示器上，计算相对坐标
		logger.Debug("窗口位于第二个显示器上")
	}

	// 计算相对坐标
	xOnSecondDisplay := mainScreenPosition.WindowX - secondDisplay.X
	yOnSecondDisplay := mainScreenPosition.WindowY - secondDisplay.Y

	logger.Debug("窗口在第二显示器上的位置: (%d,%d), 大小(%d,%d)",
		xOnSecondDisplay, yOnSecondDisplay, mainScreenPosition.WindowW, mainScreenPosition.WindowH)
	mainScreenPosition.CaptureArea = CaptureArea{
		X: xOnSecondDisplay,
		Y: yOnSecondDisplay,
		W: mainScreenPosition.WindowW,
		H: mainScreenPosition.WindowH,
	}
	return mainScreenPosition, nil
}

// SetScaleWindow macOS 实现：设置窗口大小
func (m *MacPlatform) SetScaleWindow(width, height int) {
	mack.Tell("System Events", fmt.Sprintf(`tell application "System Events"
 tell application process "%s"
  set size of window 1 to {%d, %d}
 end tell
end tell`, m.appName, width, height))
}
func (m *MacPlatform) CloseApplication() error {
	logger.Info("正在关闭BlueStacks模拟器...")

	// 查找BlueStacks进程
	ids, err := robotgo.FindIds("BlueStacks")
	if err != nil {
		logger.Info("无法查找BlueStacks进程: %v", err)
		return fmt.Errorf("无法查找BlueStacks进程: %v", err)
	}

	if len(ids) == 0 {
		logger.Info("未找到BlueStacks进程，可能已经关闭")
		return nil
	}

	logger.Info("找到BlueStacks进程 ID: %v", ids)

	// 尝试优雅关闭所有BlueStacks进程
	for _, pid := range ids {
		logger.Info("正在关闭进程 %d", pid)

		// 先尝试激活窗口然后发送关闭信号
		robotgo.ActivePid(pid)
		time.Sleep(500 * time.Millisecond)

		// 发送Alt+F4组合键来优雅关闭
		robotgo.KeyTap("f4", "alt")
		time.Sleep(1000 * time.Millisecond)
	}

	// 等待进程关闭
	time.Sleep(2 * time.Second)

	// 检查是否还有残留进程
	remainingIds, err := robotgo.FindIds("BlueStacks")
	if err == nil && len(remainingIds) > 0 {
		logger.Info("检测到残留进程，强制结束: %v", remainingIds)
		// 如果还有进程，强制结束
		for _, pid := range remainingIds {
			robotgo.Kill(pid)
		}
	}

	logger.Info("BlueStacks已成功关闭")
	return nil

}
