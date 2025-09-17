package main

import (
	"flag"
	"fmt"
	"free-hands-onmyoji/pkg/events"
	"free-hands-onmyoji/pkg/logger"
	"free-hands-onmyoji/pkg/onmyoji"
	"free-hands-onmyoji/pkg/onmyoji/breaker"
	"free-hands-onmyoji/pkg/onmyoji/general"
	"free-hands-onmyoji/pkg/onmyoji/k28"
	"free-hands-onmyoji/pkg/onmyoji/window"
	"free-hands-onmyoji/pkg/statemachine"
	"os"
	"time"

	"github.com/go-vgo/robotgo"
)

func main() {

	flag.Usage = func() {
		fmt.Printf("用法: %s -task <任务类型>\n", os.Args[0])
		fmt.Println("任务类型:")
		modes := onmyoji.GetModeNames()
		for key, value := range modes {
			fmt.Printf("  %s - %s 任务\n", key, value)
		}

		fmt.Println("默认任务类型为 k28")
		fmt.Println("display:")
		fmt.Println("  -1      - 主显示器")
		fmt.Println("  1       - 扩展显示器")
		fmt.Println("timeout:")
		fmt.Println("  0       - 不限制运行时间（默认）")
		fmt.Println("  >0      - 运行指定分钟数后自动退出")
		fmt.Println("closeBlueStacks:")
		fmt.Println("  false   - 定时退出时不关闭BlueStacks（默认）")
		fmt.Println("  true    - 定时退出时自动关闭BlueStacks")
		fmt.Println("注意：手动退出（Command+Shift+O）始终不会关闭BlueStacks")
		fmt.Println("示例: ./free-hands-onmyoji -task breaker -display 1 -timeout 60 -closeBlueStacks")
	}

	taskType := flag.String("task", "breaker", "指定任务类型: k28 或 breaker")
	displayID := flag.Int("display", -1, "指定显示器ID，默认为-1（主显示器）")
	timeout := flag.Int("timeout", 0, "程序运行时间限制（分钟），0表示不限制")
	closeBlueStacks := flag.Bool("closeBlueStacks", false, "程序退出时是否自动关闭BlueStacks")
	flag.Parse()
	value, err := onmyoji.ValidateModeExists(*taskType)
	if err != nil {
		fmt.Printf("无效的任务类型: %s\n", *taskType)
		os.Exit(1)
	}
	window.DisplayID = *displayID
	robotgo.DisplayID = *displayID // 设置robotgo的显示器ID
	// 初始化日志系统
	logger.Init()
	// 加载配置文件
	config, error := onmyoji.LoadDefaultConfig()
	if error != nil {
		panic(error)
	}
	// 创建一个通道用于控制程序退出
	exitChan := make(chan events.ExitSignal)

	// 在后台监听键盘事件
	go events.ListenForExitKey(exitChan)

	// 启动定时退出功能（如果设置了超时时间）
	events.StartTimeoutExit(*timeout, exitChan)
	// 获取游戏窗口的位置和大小
	logger.Info("正在获取游戏窗口位置和大小...")
	activeError := window.ActiveWindow("BlueStacks", 0)
	if activeError != nil {
		logger.Fatal("无法激活游戏窗口: %v", activeError)
		os.Exit(1)
	}

	windowInfo, err := window.GetWindowPosition("BlueStacks")
	logger.Info("获取到游戏窗口信息: %v", windowInfo)
	if err != nil {
		logger.Fatal("获取窗口信息失败: %v", err)
	}

	sm := statemachine.NewStateMachine()
	// Initialize tasks with the window position and size
	if window.DisplayID != -1 {
		// 如果是第二块显示器，设置偏移量 调试发现在扩展屏幕中会出现偏移问题 但是在主显示器中不会出现偏移问题
		windowInfo.OffsetX = 15
	}

	registrator := onmyoji.NewRegistrator(sm, windowInfo, config)
	taskName := *taskType
	switch taskName {
	case "k28":

		registrator.Registration(new(k28.Registrator))
	case "breaker":
		registrator.Registration(new(breaker.Registrator))
	default:
		registrator.Registration(general.Registrator{
			Path: "./" + taskName + "/",
		})
	}
	logger.Info("注册任务: %s 成功！", value)
	logger.Info("状态机开始运行...")
	logger.Info("游戏窗口已激活，开始任务执行...")
	logger.Info("当前任务: %s", sm.GetCurrentTask().Name())
	logger.Info("----------------------------------------")
	logger.Info("按下 Command+Shift+O 组合键可以停止程序运行（不会关闭BlueStacks）")
	if *timeout > 0 {
		logger.Info("程序已设置 %d 分钟后自动退出", *timeout)
		if *closeBlueStacks {
			logger.Info("定时退出时将自动关闭BlueStacks模拟器")
		} else {
			logger.Info("定时退出时不会关闭BlueStacks模拟器")
		}
	}
	time.Sleep(2 * time.Second) // 等待2秒，确保状态机初始化完成

runLoop:
	for {
		select {
		case exitSignal := <-exitChan:
			// 调用退出函数，根据退出类型和配置决定是否关闭BlueStacks
			events.ExitWithBlueStacksClose(exitSignal, *closeBlueStacks)
			break runLoop
		default:
			sm.Run()
			time.Sleep(100 * time.Millisecond) // 每秒执行一次，方便观察
		}
	}

	logger.Info("程序已退出")
}
