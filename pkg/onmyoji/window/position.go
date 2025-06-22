package window

import (
	"fmt"
	"free-hands-onmyoji/pkg/logger"
	"free-hands-onmyoji/pkg/onmyoji/entity"

	"github.com/andybrewer/mack"
)

// GetWindowPosition 获取特定应用窗口的位置和大小信息
// 参数：
//   - appName: 应用程序名称，例如 "BlueStacks"
//
// 返回：
//   - entity.WindowInfo: 包含窗口位置和大小的信息
//   - error: 如果获取失败，返回错误
func GetWindowPosition(appName string) (entity.WindowInfo, error) {
	script := fmt.Sprintf(`
tell application "System Events"
 tell application process "%s"
  get {position, size} of window 1
 end tell
end tell
`, appName)

	tell, err := mack.Tell("System Events", script)
	if err != nil {
		return entity.WindowInfo{}, fmt.Errorf("无法获取窗口信息: %v", err)
	}

	var x, y, w, h int
	_, err = fmt.Sscanf(tell, "%d, %d, %d, %d", &x, &y, &w, &h)
	if err != nil {
		return entity.WindowInfo{}, fmt.Errorf("解析窗口位置信息失败: %v", err)
	}

	logger.Info("获取窗口位置和大小: 位置(%d,%d), 大小(%d,%d)", x, y, w, h)

	return entity.WindowInfo{
		WindowX: x,
		WindowY: y,
		WindowW: w,
		WindowH: h,
	}, nil
}
