package window

import (
	"fmt"
	"free-hands-onmyoji/pkg/logger"

	"github.com/andybrewer/mack"
	"github.com/vcaesar/screenshot"
)

// GetWindowPosition 获取特定应用窗口的位置和大小信息
// 参数：
//   - appName: 应用程序名称，例如 "BlueStacks"
//
// 返回：
//   - WindowInfo: 包含窗口位置和大小的信息
//   - error: 如果获取失败，返回错误
var DisplayID int = -1 // 默认值为-1，表示使用主显示器

func GetWindowPosition(appName string) (Window, error) {
	displayId := 0
	if DisplayID != -1 {
		displayId = DisplayID
		return GetWindowPositionOnSecondDisplay(appName, displayId)
	}
	return GetWindowPositionByMaster(appName)
}

func GetWindowPositionByMaster(appName string) (Window, error) {

	script := fmt.Sprintf(`
tell application "System Events"
 tell application process "%s"
  get {position, size} of window 1
 end tell
end tell
`, appName)

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

// GetWindowPositionOnSecondDisplay 获取特定应用窗口相对于第二显示器的位置和大小信息
// 参数：
//   - appName: 应用程序名称，例如 "BlueStacks"
//
// 返回：
//   - Window: 包含窗口相对于第二显示器位置和大小的信息
//   - error: 如果获取失败，返回错误
func GetWindowPositionOnSecondDisplay(appName string, displayId int) (Window, error) {
	// 首先获取窗口在主屏幕上的位置
	mainScreenPosition, err := GetWindowPositionByMaster(appName)
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
func SetScaleWindow(appName string, width, height int) {
	mack.Tell("System Events", fmt.Sprintf(`tell application "System Events"
 tell application process "%s"
  set size of window 1 to {%d, %d}
 end tell
end tell`, appName, width, height))
}
