package k28

import (
	"fmt"
	"free-hands-onmyoji/pkg/enums"
	"free-hands-onmyoji/pkg/logger"
	"free-hands-onmyoji/pkg/onmyoji"
	"free-hands-onmyoji/pkg/onmyoji/entity"
	"free-hands-onmyoji/pkg/onmyoji/window"
	"free-hands-onmyoji/pkg/statemachine"
	"free-hands-onmyoji/pkg/types"
	"free-hands-onmyoji/pkg/utils"
	"time"
)

// MonsterDetector 小怪检测任务
type MonsterDetector struct {
	ImgTemplate   entity.ImgInfo
	window.Window // 嵌入公共字段
	conf          onmyoji.K28Config
}

func newMonsterDetectorTask(config onmyoji.Config, window window.Window, info entity.ImgInfo) *MonsterDetector {
	return &MonsterDetector{
		ImgTemplate: info,
		Window:      window,
		conf:        config.K28,
	}
}

func (t *MonsterDetector) Name() enums.TaskType {
	return enums.XunGuai
}

func (t *MonsterDetector) Execute(controller statemachine.TaskController) error {
	t.Count++ // 增加执行次数
	// 使用公共方法计算模板位置并添加随机偏移点击
	_, _, num, found, err := t.CalculateTemplatePosition(t.ImgTemplate.Image)
	if err != nil {
		return fmt.Errorf("模板图像匹配错误: %v", err)
	}

	// 如果匹配成功且相似度足够高，可以执行点击操作
	if found && num > 0.8 { // 设置相似度阈值
		// 使用公共方法点击
		clicked, err := t.ClickAtTemplatePositionWithRandomOffset(t.ImgTemplate.Image, 0.8)
		if err != nil {
			return err
		}

		if clicked {
			t.Count = 0
			//点击成功 但是有可能会跑掉所以需要再次尝试匹配一次
			logger.Info("小怪匹配成功 防止小怪跑掉，尝试再次匹配小怪")
			_, _, num, found, err = t.CalculateTemplatePosition(t.ImgTemplate.Image)
			if err != nil {
				return fmt.Errorf("模板图像匹配错误: %v", err)
			}
			if found && num > 0.8 { // 设置相似度阈值
				// 使用公共方法点击
				_, err := t.ClickAtTemplatePositionWithRandomOffset(t.ImgTemplate.Image, 0.8)
				if err != nil {
					return err
				}
			}
			logger.Info("小怪点击成功，切换到结算任务")
			controller.Next(enums.JieSuan) // 切换到结算任务
			return nil
		}
	} else {
		logger.Warn("寻找怪物模板匹配相似度过低: %.3f，跳过点击操作", num)
	}
	boss, err := controller.GetAttribute(types.Boss)
	hasBoss := boss != nil && err == nil //是否已经寻到Boss

	if hasBoss {
		t.Count = 0
		logger.Info("Boss已经寻找到，重置计数器进入寻找宝箱任务")
		time.Sleep(500 * time.Millisecond) // 等待500毫秒，确保界面稳定
		controller.Next(enums.BaoXiang)    // 切换到宝箱任务
		return nil
	}
	if t.Count >= t.conf.XunguaiFindThreshold { // 如果执行次数超过阈值
		//说明什么没有怪物了 需要向前移动 Boss任务会在结算后尝试匹配的所以这里并不需要尝试匹配boss
		logger.Info("由于没有识别到怪物 '%s' 执行次数已达 %d 次，切换到移动任务", t.Name(), t.Count)
		randomMoveCount := utils.RandomInt(1, 3)
		logger.Info("随机生成的移动次数: %d", randomMoveCount)
		controller.SetAttribute(types.MoveCount, randomMoveCount) // 重置计数器
		controller.Next(enums.Move)
		t.Count = 0 // 重置计数器
	}

	return nil
}
