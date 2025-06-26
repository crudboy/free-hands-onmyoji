package window

import (
	"fmt"
	"free-hands-onmyoji/pkg/logger"

	"github.com/go-vgo/robotgo"
)

func ActiveWindow(appName string, index int) error {
	if index < 0 {
		return fmt.Errorf("索引不能小于0")
	}
	if appName == "" {
		return fmt.Errorf("应用程序名称不能为空")
	}

	// 查找进程 ID
	ids, err := robotgo.FindIds(appName)
	if err != nil {
		return fmt.Errorf("无法查找进程 ID: %v", err)
	}
	logger.Info("查找应用程序 %s 的进程 ID: %v 如果无法唤起说明PID无效请重启相关应用", appName, ids)

	length := len(ids)

	// 检查是否找到至少一个 ID
	if length == 0 {
		return fmt.Errorf("未找到应用程序 %s 的进程 ID", appName)
	}
	if length <= index {
		return fmt.Errorf("索引 %d 超出进程 ID 列表范围，找到的进程 ID 数量为 %d", index, length)
	}
	// 激活第一个找到的进程
	robotgo.ActivePid(ids[index])

	logger.Info("成功激活进程 %d", ids[index])
	return nil
}
func AlertNotify(title string, message string) {
	robotgo.Alert(title, message)

}
