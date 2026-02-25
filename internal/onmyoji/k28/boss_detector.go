package k28

import (
	"fmt"
	"free-hands-onmyoji/internal/logger"
	"free-hands-onmyoji/internal/onmyoji"
	"free-hands-onmyoji/internal/onmyoji/window"
	"free-hands-onmyoji/internal/statemachine"
	"free-hands-onmyoji/internal/tasks"
)

// BossDetector Boss检测任务
type BossDetector struct {
	ImgTemplate   onmyoji.ImgInfo // 模板图片信息
	window.Window                 // 嵌入公共字段
}

func newBossDetectorTask(window window.Window, info onmyoji.ImgInfo) *BossDetector {
	return &BossDetector{
		ImgTemplate: info,
		Window:      window,
	}
}
func (t *BossDetector) Name() tasks.TaskType {
	return tasks.Boss
}
func (t *BossDetector) Execute(controller statemachine.TaskController) error {
	// 使用公共方法计算模板位置并添加随机偏移点击
	clicked, err := t.ClickAtTemplatePositionWithRandomOffset(t.ImgTemplate.Image, 0.8, 300)
	if err != nil {
		return fmt.Errorf("模板图像匹配错误: %v", err)
	}

	if clicked {
		controller.SetAttribute(tasks.BossState, true) // 设置点击完成标志
		controller.Next(tasks.JieSuan)                 // 切换到结算任务
	} else {
		controller.Next(tasks.XunGuai) // 切换到寻怪任务
		logger.Info("Boss没有匹配到切换到寻怪任务")
	}
	return nil

}
