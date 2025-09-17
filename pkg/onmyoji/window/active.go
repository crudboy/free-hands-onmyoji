package window

import (
	"fmt"
	"free-hands-onmyoji/pkg/logger"
	"time"

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
	//调整窗口大小
	time.Sleep(2000 * time.Millisecond) // 等待窗口激活
	SetScaleWindow(appName, 805, 485)
	return nil
}

func AlertNotify(title string, message string) {
	robotgo.Alert(title, message)

}

// CloseBlueStacks 关闭BlueStacks模拟器
func CloseBlueStacks() error {
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
