package k28

import (
	"free-hands-onmyoji/pkg/enums"
	"free-hands-onmyoji/pkg/logger"
	"free-hands-onmyoji/pkg/onmyoji/entity"
	"free-hands-onmyoji/pkg/onmyoji/window"
	"free-hands-onmyoji/pkg/statemachine"
)

// ExploreDetector 探索按钮检测任务
type ExploreDetector struct {
	TemplateImg   entity.ImgInfo
	window.Window // 嵌入公共字段
}

func newExploreDetectorTask(window window.Window, info entity.ImgInfo) *ExploreDetector {
	return &ExploreDetector{
		TemplateImg: info,
		Window:      window,
	}
}
func (t *ExploreDetector) Name() enums.TaskType {
	return enums.JinRu
}
func (t *ExploreDetector) Execute(controller statemachine.TaskController) error {
	// 使用公共方法计算模板位置并点击
	clicked, err := t.ClickAtTemplatePosition(t.TemplateImg.Image, 0.8)

	if err != nil {
		return err
	}

	if clicked {
		logger.Info("进入任务点击成功")
		controller.Next(enums.XunGuai)
		return nil
	}
	logger.Info("进入任务点击失败，尝试切换到章节点击任务")
	controller.Next(enums.ZhangJie) // 切换到章节任务

	return nil

}
