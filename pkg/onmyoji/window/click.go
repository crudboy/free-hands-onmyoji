package window

import (
	"fmt"
	"free-hands-onmyoji/pkg/logger"
	"free-hands-onmyoji/pkg/utils"
	"image"
	"math/rand"
	"time"

	"github.com/go-vgo/robotgo"
)

// Window 定义任务的公共字段
type Window struct {
	WindowX int // 程序窗口在屏幕上的起始位置X（偏移量）
	WindowY int // 程序窗口在屏幕上的起始位置Y（偏移量）
	WindowH int // 截图区域的高度
	WindowW int // 截图区域的宽度
	OffsetX int // 点击位置的X轴偏移调整（可选）
	OffsetY int // 点击位置的Y轴偏移调整（可选）
	Count   int // 任务执行计数
}

// CalculateTemplatePosition 计算模板图片在屏幕中的位置
// 返回：
// - 屏幕坐标X
// - 屏幕坐标Y
// - 相似度
// - 是否找到图像
// - 错误信息
func (tc *Window) CalculateTemplatePosition(templateImage image.Image) (int, int, float32, bool, error) {
	// 窗口区域截图
	img, err := robotgo.CaptureImg(tc.WindowX, tc.WindowY, tc.WindowW, tc.WindowH)
	if err != nil {
		return 0, 0, 0, false, err
	}

	// 计算出模板图像在窗口区域内的坐标
	tempPosX, tempPosY, num := utils.FindTempPos(templateImage, img)

	if num == 0 {
		logger.Warn("未找到模板图像，无法执行任务")
		return 0, 0, num, false, nil
	}

	// 获取模板图像的尺寸，用于计算中心点
	templateBounds := templateImage.Bounds()
	templateWidth := templateBounds.Max.X - templateBounds.Min.X
	templateHeight := templateBounds.Max.Y - templateBounds.Min.Y

	// QQ截图的分辨率是实际的2倍，因此需要除以2
	tempPosX = tempPosX / 2
	tempPosY = tempPosY / 2
	templateWidth = templateWidth / 2
	templateHeight = templateHeight / 2

	// 计算模板图像中心点在窗口区域内的坐标
	centerX := tempPosX + templateWidth/2
	centerY := tempPosY + templateHeight/2

	// 计算模板图像中心点在屏幕上的绝对坐标，并应用偏移调整
	screenPosX := tc.WindowX + centerX + tc.OffsetX
	screenPosY := tc.WindowY + centerY + tc.OffsetY

	logger.Debug("截取区域: (%d, %d, %d, %d)", tc.WindowX, tc.WindowY, tc.WindowW, tc.WindowH)
	logger.Debug("原始模板匹配位置(2倍分辨率): (%d, %d), 相似度: %.3f", tempPosX*2, tempPosY*2, num)
	logger.Debug("调整后模板匹配位置: (%d, %d)", tempPosX, tempPosY)
	logger.Debug("调整后模板尺寸: %dx%d", templateWidth, templateHeight)
	logger.Debug("模板中心点(窗口内): centerX=%d, centerY=%d", centerX, centerY)
	logger.Debug("偏移调整: offsetX=%d, offsetY=%d", tc.OffsetX, tc.OffsetY)
	logger.Debug("最终屏幕坐标: screenPosX=%d, screenPosY=%d", screenPosX, screenPosY)

	return screenPosX, screenPosY, num, true, nil
}

// ClickAtTemplatePosition 计算模板图片位置并点击
// 如果找到模板图片，则点击并返回true，否则返回false
func (tc *Window) ClickAtTemplatePosition(templateImage image.Image, similarityThreshold float32) (bool, error) {
	if similarityThreshold <= 0.5 {
		return false, fmt.Errorf("相似度阈值过低，无法执行点击操作")
	}
	screenPosX, screenPosY, num, found, err := tc.CalculateTemplatePosition(templateImage)
	if err != nil {
		return false, err
	}

	if !found || num <= similarityThreshold {
		logger.Warn("模板匹配相似度过低: %.3f，跳过点击操作", num)
		return false, nil
	}

	logger.Info("点击位置: (%d, %d)，相似度: %.3f", screenPosX, screenPosY, num)
	robotgo.MoveClick(screenPosX, screenPosY)
	return true, nil
}

// ClickAtTemplatePositionWithRandomOffset 计算模板图片位置并添加随机偏移后点击
// 增加随机性，让点击更像人为操作
func (tc *Window) ClickAtTemplatePositionWithRandomOffset(templateImage image.Image, similarityThreshold float32) (bool, error) {
	screenPosX, screenPosY, num, found, err := tc.CalculateTemplatePosition(templateImage)
	if err != nil {
		return false, err
	}

	if !found || num <= similarityThreshold {
		logger.Warn("模板匹配相似度过低: %.3f，跳过点击操作", num)
		return false, nil
	}

	// 生成点击位置的随机偏移，增加真实性
	_, randomX := utils.RandomNormalInt64(int64(screenPosX-20), int64(screenPosX+20), int64(screenPosX), 10)
	_, randomY := utils.RandomNormalInt64(int64(screenPosY-20), int64(screenPosY+20), int64(screenPosY), 10)

	logger.Info("鼠标移动到: (%d, %d)", randomX, randomY)
	robotgo.Move(int(randomX), int(randomY))
	time.Sleep(200 * time.Millisecond) // 等待鼠标移动完成
	robotgo.Click("left")
	return true, nil
}

// ClickFloor 点击地板使角色向前移动
// rightOffset 允许指定向右的额外偏移量，传入0则使用默认位置
// 返回实际点击的坐标
func (tc *Window) ClickFloor(rightOffset int, downOffset int) (int, int) {
	// 实现Y轴向前移动（点击地板）
	// 计算点击位置 - 在屏幕中央偏下的位置点击，模拟点击地板
	centerX := tc.WindowX + tc.WindowW/2 + tc.OffsetX
	// Y轴点击位置，选择窗口偏下位置，通常是地板所在位置
	forwardY := tc.WindowY + int(float64(tc.WindowH)*0.7) + tc.OffsetY

	// 添加随机偏移，使移动更自然
	randOffsetX := rand.Intn(40) - 20 // -20 到 20 之间的随机值
	randOffsetY := rand.Intn(20) - 10 // -10 到 10 之间的随机值

	clickX := centerX + randOffsetX + rightOffset
	clickY := forwardY + randOffsetY + downOffset

	logger.Info("点击地板位置: (%d, %d) 使角色向前移动", clickX, clickY)
	robotgo.MoveClick(clickX, clickY, "left")

	return clickX, clickY
}
