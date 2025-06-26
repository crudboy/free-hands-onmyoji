package main

import (
	"flag"
	"fmt"
	"free-hands-onmyoji/pkg/events"
	"free-hands-onmyoji/pkg/logger"
	"free-hands-onmyoji/pkg/onmyoji"
	"free-hands-onmyoji/pkg/onmyoji/breaker"
	"free-hands-onmyoji/pkg/onmyoji/k28"
	"free-hands-onmyoji/pkg/onmyoji/window"
	"free-hands-onmyoji/pkg/statemachine"
	"os"
	"time"
)

func main() {

	flag.Usage = func() {
		fmt.Printf("用法: %s -task <任务类型>\n", os.Args[0])
		fmt.Println("任务类型:")
		fmt.Println("  k28     - K28任务")
		fmt.Println("  breaker - 突破任务")
		fmt.Println("默认任务类型为 k28")
		fmt.Println("示例: ./free-hands-onmyoji -task breaker")
	}
	taskType := flag.String("task", "breaker", "指定任务类型: k28 或 breaker")
	displayID := flag.Int("display", -1, "指定显示器ID，默认为-1（主显示器）")
	window.DisplayID = *displayID // 设置全局显示器ID
	flag.Parse()
	// 初始化日志系统
	logger.Init()

	// 加载配置文件
	config, error := onmyoji.LoadDefaultConfig()
	if error != nil {
		panic(error)
	}
	// 创建一个通道用于控制程序退出
	exitChan := make(chan bool)

	// 在后台监听键盘事件
	go events.ListenForExitKey(exitChan)
	// 获取游戏窗口的位置和大小
	logger.Info("正在获取游戏窗口位置和大小...")
	windowInfo, _, err := window.GetWindowPosition("BlueStacks")
	if err != nil {
		logger.Fatal("获取窗口信息失败: %v", err)
	}

	sm := statemachine.NewStateMachine()
	// Initialize tasks with the window position and size
	registrator := onmyoji.NewRegistrator(sm, windowInfo, config)
	taskName := *taskType
	switch taskName {
	case "k28":
		logger.Info("注册K28任务...")
		registrator.Registration(new(k28.Registrator))
	case "breaker":
		logger.Info("注册突破任务...")
		registrator.Registration(new(breaker.Registrator))
	default:
		logger.Fatal("未知任务类型: %s", taskName)
		os.Exit(1)
	}
	logger.Info("状态机开始运行...")
	window.ActiveWindow("BlueStacks", 0)
	logger.Info("游戏窗口已激活，开始任务执行...")
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
