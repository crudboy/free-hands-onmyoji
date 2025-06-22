package events

import (
	"free-hands-onmyoji/pkg/logger"

	hook "github.com/robotn/gohook"
)

// ListenForExitKey 监听组合键 Command+Shift+O 用于退出程序
// 当组合键被检测到时，会通过退出通道发送一个信号
func ListenForExitKey(exitChan chan bool) {
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
				exitChan <- true
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
