package k28

import (
	"free-hands-onmyoji/internal/logger"
	"free-hands-onmyoji/internal/onmyoji"
	"free-hands-onmyoji/internal/onmyoji/window"
	"free-hands-onmyoji/internal/statemachine"
	"free-hands-onmyoji/internal/tasks"
)

// ExploreDetector 探索按钮检测任务
type ExploreDetector struct {
	ImgTemplate   onmyoji.ImgInfo
	window.Window // 嵌入公共字段
}

func newExploreDetectorTask(window window.Window, info onmyoji.ImgInfo) *ExploreDetector {
	return &ExploreDetector{
		ImgTemplate: info,
		Window:      window,
	}
}
func (t *ExploreDetector) Name() tasks.TaskType {
	return tasks.JinRu
}
func (t *ExploreDetector) Execute(controller statemachine.TaskController) error {
	// 使用公共方法计算模板位置并点击
	clicked, err := t.ClickAtTemplatePosition(t.ImgTemplate.Image, 0.8)

	if err != nil {
		return err
	}

	if clicked {
		logger.Info("进入任务点击成功")
		controller.Next(tasks.XunGuai)
		return nil
	}
	logger.Info("进入任务点击失败，尝试切换到章节点击任务")
	controller.Next(tasks.ZhangJie) // 切换到章节任务

	return nil

}
