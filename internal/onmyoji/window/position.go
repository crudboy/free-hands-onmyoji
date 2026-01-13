package window

import (
	"free-hands-onmyoji/internal/logger"

	"github.com/vcaesar/screenshot"
)

// DisplayID 显示器ID，默认值为-1，表示使用主显示器
var DisplayID int = -1

// DisplayInfo 存储显示器的位置和大小信息
type DisplayInfo struct {
	X int
	Y int
	W int
	H int
}

// getDisplaysInfo 获取所有连接的显示器信息
// 返回:
//   - []DisplayInfo: 所有显示器的位置和大小信息
//   - error: 如果获取失败，返回错误
func getDisplaysInfo() ([]DisplayInfo, error) {
	numDisplays := screenshot.NumActiveDisplays()
	logger.Debug("检测到 %d 个显示器", numDisplays)
	displays := make([]DisplayInfo, numDisplays)
	for i := 0; i < numDisplays; i++ {
		bounds := screenshot.GetDisplayBounds(i)
		displays[i] = DisplayInfo{
			X: bounds.Min.X,
			Y: bounds.Min.Y,
			W: bounds.Dx(),
			H: bounds.Dy(),
		}
		logger.Debug("显示器 %d: 位置(%d,%d), 大小(%d,%d), 右下角(%d,%d)",
			i+1, bounds.Min.X, bounds.Min.Y, bounds.Dx(), bounds.Dy(),
			bounds.Min.X+bounds.Dx(), bounds.Min.Y+bounds.Dy())
	}
	return displays, nil
}

// GetWindowPosition 获取特定应用窗口的位置和大小信息
// 参数：
//   - appName: 应用程序名称，例如 "BlueStacks"
//
// 返回：
//   - Window: 包含窗口位置和大小的信息
//   - error: 如果获取失败，返回错误
func GetWindowPosition(appName string) (Window, error) {
	return GetPlatform().GetWindowPosition()
}

// SetScaleWindow 设置窗口的大小
// 参数：
//   - appName: 应用程序名称
//   - width: 窗口宽度
//   - height: 窗口高度
func SetScaleWindow(appName string, width, height int) {
	GetPlatform().SetScaleWindow(width, height)
}
