package onmyoji

import (
	"fmt"
	"image"
)

var modes = map[string]string{
	"guren":         "业原火",
	"k28":           "困28",
	"breaker":       "突破",
	"mitama":        "御灵",
	"limitedEvents": "限定活动",
}

type ImgInfo struct {
	Path    string
	ImgMaxX int
	ImgMaxY int
	Image   image.Image
}

func ValidateModeExists(mode string) (string, error) {
	value, exists := modes[mode]
	if !exists {
		return "", fmt.Errorf("模式不存在: %s", mode)
	}
	return value, nil
}
func GetModeNames() map[string]string {
	return modes
}
