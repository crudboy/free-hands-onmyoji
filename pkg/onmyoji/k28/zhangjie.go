package k28

import (
	"fmt"
	"free-hands-onmyoji/pkg/enums"
	"free-hands-onmyoji/pkg/logger"
	"free-hands-onmyoji/pkg/onmyoji"
	"free-hands-onmyoji/pkg/onmyoji/entity"
	"free-hands-onmyoji/pkg/onmyoji/window"
	"free-hands-onmyoji/pkg/statemachine"
)

type ZhangJie struct {
	TemplateImg entity.ImgInfo // 模板图片信息
	window.Window
	conf onmyoji.K28Config
}

func newZhangJieTask(config onmyoji.Config, window window.Window, info entity.ImgInfo) *ZhangJie {
	return &ZhangJie{
		TemplateImg: info,
		Window:      window,
		conf:        config.K28,
	}
}
func (t *ZhangJie) Name() enums.TaskType {
	return enums.ZhangJie
}

func (t *ZhangJie) Execute(controller statemachine.TaskController) error {
	t.Count++ // 增加执行次数

	// 使用公共方法计算模板位置并添加随机偏移点击
	clicked, err := t.ClickAtTemplatePositionWithRandomOffset(t.TemplateImg.Image, 0.8)
	if err != nil {
		return fmt.Errorf("模板图像匹配错误: %v", err)
	}

	if clicked {
		logger.Info("章节点击成功")
		t.Count = 0                  // 重置执行次数
		controller.Next(enums.JinRu) // 切换到进入任务
		return nil
	}
	if t.Count > t.conf.ZhangjieFindThreshold {
		logger.Warn("选择章节任务执行失败超过阈值 %d，尝试执行进入任务", t.conf.ZhangjieFindThreshold)
		t.Count = 0 // 重置执行次数
		logger.Info("章节没有匹配到，切换到进入任务")
		controller.Next(enums.JinRu) // 切换到进入任务
	}

	return nil

}
