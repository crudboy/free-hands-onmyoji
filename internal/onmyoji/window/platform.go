package window

import (
	"fmt"
	"free-hands-onmyoji/internal/logger"
	"time"

	"github.com/go-vgo/robotgo"
)

// Platform 定义跨平台窗口操作接口
type Platform interface {
	// GetWindowPosition 获取特定应用窗口的位置和大小信息
	// 参数：appName - 应用程序名称
	// 返回：Window 窗口信息，error 错误信息
	GetWindowPosition() (Window, error)

	// GetWindowPositionOnSecondDisplay 获取特定应用窗口相对于第二显示器的位置和大小信息
	// 参数：appName - 应用程序名称，displayId - 显示器ID
	// 返回：Window 窗口信息，error 错误信息
	GetWindowPositionOnSecondDisplay(displayId int) (Window, error)

	// SetScaleWindow 设置窗口的大小
	// 参数：appName - 应用程序名称，width - 宽度，height - 高度
	SetScaleWindow(width, height int)
	ActiveWindow(index int) error
	CloseApplication() error
	GetAppName() string
}
type BasePlatform struct {
	appName string
}

func (p *BasePlatform) GetAppName() string {
	return p.appName
}
func (p *BasePlatform) ActiveWindow(index int) error {

	if index < 0 {
		return fmt.Errorf("索引不能小于0")
	}
	if p.appName == "" {
		return fmt.Errorf("应用程序名称不能为空")
	}

	// 查找进程 ID
	ids, err := robotgo.FindIds(p.appName)
	if err != nil {
		return fmt.Errorf("无法查找进程 ID: %v", err)
	}
	logger.Info("查找应用程序 %s 的进程 ID: %v 如果无法唤起说明PID无效请重启相关应用", p.appName, ids)
	length := len(ids)

	// 检查是否找到至少一个 ID
	if length == 0 {
		return fmt.Errorf("未找到应用程序 %s 的进程 ID", p.appName)
	}
	if length <= index {
		return fmt.Errorf("索引 %d 超出进程 ID 列表范围，找到的进程 ID 数量为 %d", index, length)
	}

	// 激活指定索引的进程
	robotgo.ActivePid(ids[index])
	logger.Info("成功激活进程 %d", ids[index])

	// 等待窗口激活
	time.Sleep(2000 * time.Millisecond)
	SetScaleWindow(p.appName, 805, 485)

	return nil

}

// MacPlatform macOS 平台实现
type MacPlatform struct {
	BasePlatform
}

// WindowsPlatform Windows 平台实现
type WindowsPlatform struct {
	BasePlatform
}

var platformImpl Platform

// GetPlatform 获取当前平台的实现
func GetPlatform() Platform {
	return platformImpl
}

// InitPlatform 初始化平台实现
func InitPlatform(os string) {
	switch os {
	case "mac":
		platformImpl = &MacPlatform{
			BasePlatform: BasePlatform{
				appName: "BlueStacks",
			},
		}
	case "windows":
		platformImpl = &WindowsPlatform{
			BasePlatform: BasePlatform{
				appName: "BlueStacks",
			},
		}
	default:
		platformImpl = &MacPlatform{
			BasePlatform: BasePlatform{
				appName: "BlueStacks",
			},
		} // 默认使用 macOS 实现
	}
}
