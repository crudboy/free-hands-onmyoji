package k28

import (
	"fmt"
	"free-hands-onmyoji/pkg/enums"
	"free-hands-onmyoji/pkg/logger"
	"free-hands-onmyoji/pkg/onmyoji"
	"free-hands-onmyoji/pkg/onmyoji/window"
	"free-hands-onmyoji/pkg/statemachine"
	"free-hands-onmyoji/pkg/types"
	"free-hands-onmyoji/pkg/utils"
	"time"
)

/*
	LevelCompletionDetector 过关检测
*/

type LevelCompletionDetector struct {
	ImgTemplate   onmyoji.ImgInfo // 模板图像信息
	window.Window                 // 嵌入公共字段
	runThreshold  int             // 运行次数阈值
	conf          onmyoji.K28Config
}

func newLevelCompletionDetectorTask(config onmyoji.Config, window window.Window, info onmyoji.ImgInfo) *LevelCompletionDetector {
	return &LevelCompletionDetector{
		ImgTemplate: info,
		Window:      window,
		conf:        config.K28,
	}
}
func (t *LevelCompletionDetector) Name() enums.TaskType {
	return enums.JieSuan
}
func (t *LevelCompletionDetector) Execute(controller statemachine.TaskController) error {
	t.runThreshold++
	time.Sleep(time.Duration(t.conf.LevelCompletionWaitTime) * time.Millisecond) // 等待600毫秒，确保界面稳定
	// 使用公共方法计算模板位置并添加随机偏移点击
	clicked, err := t.ClickAtTemplatePositionWithRandomOffset(t.ImgTemplate.Image, 0.8)
	if err != nil {
		return fmt.Errorf("模板图像匹配错误: %v", err)
	}
	if t.runThreshold > t.conf.LevelCompletionThreshold {
		logger.Warn("结算任务执行次数超过阈值 %d，跳过结算点击操作，进入寻怪任务", t.conf.LevelCompletionThreshold)
		controller.Next(enums.XunGuai) // 切换到寻怪任务
		t.runThreshold = 0             // 重置运行次数阈值
		return nil
	}
	if clicked {
		logger.Info("结算界面点击成功")
		t.runThreshold = 0 // 重置运行次数阈值
		t.Count++          // 增加执行次数
		logger.Debug("结算任务执行成功，当前执行次数: %d", t.Count)
		boss, err := controller.GetAttribute(types.Boss)
		hasBoss := boss != nil && err == nil //是否已经寻到Boss
		if hasBoss {
			t.Count = 0 // 如果已经寻到Boss，重置执行次数
		}

		//由于boss也是有结算画面的所以这里是4当第二次进入后count默认就是1了
		if t.Count >= t.conf.AfterMinionAttemptsBossCheck && !hasBoss {
			logger.Info("结算任务执行超过 %d，尝试检测Boss", t.conf.AfterMinionAttemptsBossCheck)
			randomInt := utils.RandomInt(500, 857)
			time.Sleep(time.Duration(randomInt) * time.Millisecond) // 等待随机时间
			logger.Info("检测Boss任务开始执行")
			controller.Next(enums.Boss) // 切换到Boss任务
			t.runThreshold = 0          // 重置运行次数阈值
			return nil
		} else {
			time.Sleep(200 * time.Millisecond) // 等待切换完成
			logger.Info("结算任务执行完成，切换到寻怪任务")
			controller.Next(enums.XunGuai) // 切换到寻怪任务
			t.runThreshold = 0             // 重置运行次数阈值
			return nil
		}

	}
	logger.Info("未找到结算界面或相似度不足，跳过点击操作")
	return nil

}
