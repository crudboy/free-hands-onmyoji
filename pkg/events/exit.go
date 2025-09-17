package events

import (
	"free-hands-onmyoji/pkg/logger"
	"free-hands-onmyoji/pkg/onmyoji/window"
	"time"

	hook "github.com/robotn/gohook"
)

// ExitType 退出类型
type ExitType int

const (
	// ManualExit 手动退出（不关闭BlueStacks）
	ManualExit ExitType = iota
	// TimeoutExit 定时退出（根据配置决定是否关闭BlueStacks）
	TimeoutExit
)

// ExitSignal 退出信号
type ExitSignal struct {
	Type ExitType
}

// ListenForExitKey 监听组合键 Command+Shift+O 用于退出程序
// 当组合键被检测到时，会通过退出通道发送一个手动退出信号
func ListenForExitKey(exitChan chan ExitSignal) {
	// 注册键盘事件
	evChan := hook.Start()
	defer hook.End()

	logger.Info("键盘监听已启动，正在等待组合键...")

	// 跟踪按键状态
	keysPressed := make(map[uint16]bool)

	for ev := range evChan {
		// 键盘按下事件
		if ev.Kind == hook.KeyDown || ev.Kind == hook.KeyHold {
			// 记录按键状态
			keysPressed[ev.Rawcode] = true

			// 检查组合键 - 测试多种可能的键码组合
			// Command键通常是55或56
			// Shift键通常是57或60
			// O键的Rawcode通常是31
			cmdKey := keysPressed[55]
			shiftKey := keysPressed[56]
			// 检查O键 (多种可能的检测方式)
			oKey := ev.Rawcode == 31

			// 如果当前按键是O，且Command和Shift都被按下
			if oKey && cmdKey && shiftKey {
				logger.Info("检测到 Command+Shift+O 组合键！")
				logger.Info("程序即将退出...")
				exitChan <- ExitSignal{Type: ManualExit}
				return
			}
		}

		// 键盘释放事件
		if ev.Kind == hook.KeyUp {
			if ev.Keycode > 0 {
				// 清除按键状态
				delete(keysPressed, ev.Keycode)
			}
		}
	}
}

// StartTimeoutExit 启动定时退出功能
// timeout: 超时时间（分钟），0表示不限制时间
// exitChan: 退出通道
func StartTimeoutExit(timeout int, exitChan chan ExitSignal) {
	if timeout <= 0 {
		return // 不启动定时器
	}

	duration := time.Duration(timeout) * time.Minute
	logger.Info("程序将在 %d 分钟后自动退出", timeout)

	go func() {
		timer := time.NewTimer(duration)
		defer timer.Stop()

		<-timer.C
		logger.Info("运行时间已达 %d 分钟，程序自动退出", timeout)
		exitChan <- ExitSignal{Type: TimeoutExit}
	}()
}

// ExitWithBlueStacksClose 根据退出类型和配置决定是否关闭BlueStacks
// exitSignal: 退出信号，包含退出类型
// closeBlueStacks: 配置参数，是否关闭BlueStacks模拟器
func ExitWithBlueStacksClose(exitSignal ExitSignal, closeBlueStacks bool) {
	switch exitSignal.Type {
	case ManualExit:
		logger.Info("手动退出程序，保持BlueStacks运行...")
	case TimeoutExit:
		if closeBlueStacks {
			logger.Info("定时退出程序并关闭BlueStacks...")
			err := window.CloseBlueStacks()
			if err != nil {
				logger.Error("关闭BlueStacks时出现错误: %v", err)
			}
		} else {
			logger.Info("定时退出程序，保持BlueStacks运行...")
		}
	default:
		logger.Info("正在退出程序...")
	}
}
