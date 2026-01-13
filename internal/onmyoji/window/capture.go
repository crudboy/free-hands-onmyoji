package window

import (
	"fmt"
	"free-hands-onmyoji/internal/logger"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"time"

	"github.com/go-vgo/robotgo"
)

// CaptureAndSave 截取窗口区域并保存到文件
// 返回保存的文件路径和错误信息
func (tc *Window) CaptureAndSave() (string, error) {
	var (
		img image.Image
		err error
	)

	// 根据显示器ID选择截图方式
	if DisplayID != -1 {
		logger.Info("截图第二块屏幕 (%d, %d, %d, %d)", tc.CaptureArea.X, tc.CaptureArea.Y, tc.CaptureArea.W, tc.CaptureArea.H)
		capture := tc.CaptureArea
		robotgo.DisplayID = 1 // 设置显示器ID为1
		img, err = robotgo.CaptureImg(capture.X, capture.Y, capture.W, capture.H)
		if err != nil {
			return "", fmt.Errorf("获取第二块显示器窗口截图失败: %v", err)
		}
	} else {
		img, err = robotgo.CaptureImg(tc.WindowX, tc.WindowY, tc.WindowW, tc.WindowH)
		if err != nil {
			return "", fmt.Errorf("截图失败: %v", err)
		}
	}

	// 创建 capture 文件夹（如果不存在）
	captureDir := "capture"
	if err := os.MkdirAll(captureDir, 0755); err != nil {
		return "", fmt.Errorf("创建capture文件夹失败: %v", err)
	}

	// 生成带时间戳的文件名
	timestamp := time.Now().Format("20060102_150405")
	filename := filepath.Join(captureDir, fmt.Sprintf("screenshot_%s.png", timestamp))

	// 创建文件
	file, err := os.Create(filename)
	if err != nil {
		return "", fmt.Errorf("创建文件失败: %v", err)
	}
	defer file.Close()

	// 保存图片为PNG格式
	if err := png.Encode(file, img); err != nil {
		return "", fmt.Errorf("保存图片失败: %v", err)
	}

	logger.Info("截图已保存: %s", filename)
	return filename, nil
}
