package breaker

import (
	"free-hands-onmyoji/pkg/enums"
	"free-hands-onmyoji/pkg/logger"
	"free-hands-onmyoji/pkg/onmyoji/entity"
	"free-hands-onmyoji/pkg/onmyoji/window"
	"free-hands-onmyoji/pkg/statemachine"
	"time"
)

type RewardDetector struct {
	ImgTemplate entity.ImgInfo // 模板图片信息
	window.Window
}

func newRewardDetector(window window.Window, info entity.ImgInfo) *RewardDetector {
	return &RewardDetector{
		ImgTemplate: info,
		Window:      window,
	}
}
func (t *RewardDetector) Name() enums.TaskType {
	return enums.BreakerReward
}
func (t *RewardDetector) Execute(controller statemachine.TaskController) error {
	// 使用公共方法计算模板位置并添加随机偏移点击
	clicked, err := t.ClickAtTemplatePositionWithRandomOffset(t.ImgTemplate.Image, 0.8)
	if err != nil {
		return err
	}
	if clicked {
		logger.Info("奖励检测成功，点击奖励按钮后等待2秒")
		time.Sleep(2 * time.Second)
	}

	controller.Next(enums.BreakerPlayer) // 切换到玩家检测状态
	return nil
}
