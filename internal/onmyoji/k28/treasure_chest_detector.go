package k28

import (
	"fmt"
	"free-hands-onmyoji/internal/logger"
	"free-hands-onmyoji/internal/onmyoji"
	"free-hands-onmyoji/internal/onmyoji/window"
	"free-hands-onmyoji/internal/statemachine"
	"free-hands-onmyoji/internal/tasks"
	"time"
)

// TreasureChestDetector 宝箱检测任务
type TreasureChestDetector struct {
	ImgTemplate   onmyoji.ImgInfo // 模板图片信息
	window.Window                 // 嵌入公共字段
	conf          onmyoji.K28Config
}

func newTreasureChestDetectorTask(conf onmyoji.Config, window window.Window, info onmyoji.ImgInfo) *TreasureChestDetector {
	return &TreasureChestDetector{
		ImgTemplate: info,
		Window:      window,
		conf:        conf.K28,
	}
}

func (t *TreasureChestDetector) Name() tasks.TaskType {
	return tasks.BaoXiang
}

func (t *TreasureChestDetector) Execute(controller statemachine.TaskController) error {
	// 使用公共方法计算模板位置并添加随机偏移点击
	clicked, err := t.ClickAtTemplatePositionWithRandomOffset(t.ImgTemplate.Image, 0.8, 300)
	if err != nil {
		return fmt.Errorf("模板图像匹配错误: %v", err)
	}

	if clicked {
		logger.Info("宝箱点击成功")
		time.Sleep(200 * time.Millisecond) // 等待点击操作完成
		// 点击地板
		logger.Info("宝箱点击后，点击地板")
		t.ClickFloor(0, 30)
		time.Sleep(200 * time.Millisecond) // 等待地板点击完成
	} else {
		//宝箱寻找完成
		time.Sleep(time.Duration(t.conf.ChestWaitTime) * time.Millisecond) // 等待切换完成
		controller.ClearAttributes()                                       // 清除所有任务属性
		controller.Next(tasks.JinRu)                                       // 切换到进入任务
		logger.Info("宝箱没有匹配到，切换到进入任务")
	}
	return nil

}
