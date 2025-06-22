package main

import (
	"free-hands-onmyoji/pkg/events"
	"free-hands-onmyoji/pkg/logger"
	"free-hands-onmyoji/pkg/onmyoji/k28"
	"free-hands-onmyoji/pkg/onmyoji/window"
	"free-hands-onmyoji/pkg/statemachine"
	"time"
)

func main() {
	// 初始化日志系统
	logger.Init()

	// 创建一个通道用于控制程序退出
	exitChan := make(chan bool)

	// 在后台监听键盘事件
	go events.ListenForExitKey(exitChan)

	// 获取游戏窗口的位置和大小
	windowInfo, err := window.GetWindowPosition("BlueStacks")
	if err != nil {
		logger.Fatal("获取窗口信息失败: %v", err)
	}

	sm := statemachine.NewStateMachine()
	// Initialize tasks with the window position and size
	k28.Registration(sm, windowInfo)
	logger.Info("状态机开始运行...")
	logger.Info("每个任务将根据自己的逻辑决定何时切换到下一个任务")
	logger.Info("当前任务: %s", sm.GetCurrentTask().Name())
	logger.Info("----------------------------------------")
	logger.Info("按下 Command+Shift+O 组合键可以停止程序运行")
	time.Sleep(2 * time.Second) // 等待2秒，确保状态机初始化完成

runLoop:
	for {
		select {
		case <-exitChan:
			logger.Info("正在退出程序...")
			break runLoop
		default:
			sm.Run()
			time.Sleep(1500 * time.Millisecond) // 每秒执行一次，方便观察
		}
	}

	logger.Info("程序已退出")
}
