package window

import (
	"fmt"
	"free-hands-onmyoji/pkg/logger"

	"github.com/go-vgo/robotgo"
)

func ActiveWindow(appName string, index int) (bool, error) {
	if index < 0 {
		return false, fmt.Errorf("索引不能小于0")
	}
	if appName == "" {
		return false, fmt.Errorf("应用程序名称不能为空")
	}

	// 查找进程 ID
	ids, err := robotgo.FindIds(appName)
	if err != nil {
		return false, fmt.Errorf("无法查找进程 ID: %v", err)
	}

	// 检查是否找到至少一个 ID
	if len(ids) == 0 {
		return false, fmt.Errorf("未找到应用程序 %s 的进程 ID", appName)
	}
	pid := robotgo.GetPid()
	logger.Info("当前进程 ID: %d", pid)

	// 激活第一个找到的进程
	err = robotgo.ActivePid(ids[index])
	if err != nil {
		return false, fmt.Errorf("无法激活进程 %d: %v", ids[index], err)
	}

	logger.Info("成功激活进程 %d", ids[index])
	return true, nil
}
func AlertNotify(title string, message string) {
	robotgo.Alert(title, message)

}
