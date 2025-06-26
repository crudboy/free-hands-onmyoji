package breaker

import (
	"fmt"
	"free-hands-onmyoji/pkg/enums"
	"free-hands-onmyoji/pkg/logger"
	"free-hands-onmyoji/pkg/onmyoji"
	"free-hands-onmyoji/pkg/onmyoji/entity"
	"free-hands-onmyoji/pkg/onmyoji/window"
	"free-hands-onmyoji/pkg/statemachine"
	"free-hands-onmyoji/pkg/utils"
	"image"
	"os"
	"path/filepath"
	"strings"
)

type Registrator struct {
}

func (r Registrator) Registration(machine *statemachine.StateMachine, info entity.WindowInfo, config onmyoji.Config, imgMap map[string]entity.ImgInfo) error {
	w := window.Window{
		WindowX: info.WindowX,
		WindowY: info.WindowY,
		WindowH: info.WindowH,
		WindowW: info.WindowW,
	}

	imgList := convertImgList(imgMap)
	onmyoji.Registration(machine, newPlayerDetector(w, imgList))
	onmyoji.Registration(machine, newAttackDetector(w, imgMap[string(enums.BreakerAttack)], config))
	onmyoji.Registration(machine, newRewardDetector(w, imgMap[string(enums.BreakerReward)]))
	onmyoji.Registration(machine, newBreakerFailDetector(w, imgMap[string(enums.BreakerFail)]))
	onmyoji.Registration(machine, newBreakerWinDetector(w, imgMap[string(enums.BreakerWin)]))

	return nil
}

func convertImgList(imgMap map[string]entity.ImgInfo) []entity.ImgInfo {
	imgList := make([]entity.ImgInfo, 0, len(imgMap))
	for imgName, imgInfo := range imgMap {
		if strings.HasPrefix(imgName, "breaker_player_") {
			imgList = append(imgList, imgInfo)
		}
		// 注册每个图片
	}
	return imgList
}
func (r Registrator) LoadImageTemplates() (map[string]entity.ImgInfo, error) {

	logger.Info("加载突破任务模板图片")

	// 初始化模板图片map
	imgMap := make(map[string]entity.ImgInfo)

	imgPath := "./breaker/"
	files, err := os.ReadDir(imgPath)
	if err != nil {
		panic(fmt.Errorf("读取图片目录失败: %v", err))
	}
	for _, imgFile := range files {
		reader, err := os.Open(filepath.Join(imgPath, imgFile.Name()))
		if err != nil {
			continue
		}

		// 使用闭包确保每次迭代结束时关闭文件
		func(file *os.File) {
			defer file.Close()

			if file.Name() == filepath.Join(imgPath, ".DS_Store") {
				return
			}
			im, _, err := image.DecodeConfig(file)
			if err != nil {
				logger.Error("图片解码错误 %s: %v", imgFile.Name(), err)
				return // 跳过错误图片而不是终止程序
			}
			// 不要文件后缀
			fileName := imgFile.Name()[:len(imgFile.Name())-len(filepath.Ext(imgFile.Name()))]

			imgMap[fileName] = entity.ImgInfo{
				Path:    imgPath + imgFile.Name(),
				ImgMaxX: im.Width,
				ImgMaxY: im.Height,
				Image:   utils.ReadPic(imgPath + imgFile.Name()),
			}
		}(reader)
	}

	// 检查是否成功加载了所有必要的图片
	requiredImages := []func() (string, bool){
		func() (string, bool) {
			_, exists := imgMap[string(enums.BreakerAttack)]
			return string(enums.BreakerAttack), exists
		},
		func() (string, bool) {
			_, exists := imgMap[string(enums.BreakerReward)]
			return string(enums.BreakerReward), exists
		},
		func() (string, bool) {
			_, exists := imgMap[string(enums.BreakerFail)]
			return string(enums.BreakerFail), exists
		},
		func() (string, bool) {
			_, exists := imgMap[string(enums.BreakerWin)]
			return string(enums.BreakerWin), exists
		},
		func() (string, bool) {
			for imgName := range imgMap {
				if strings.HasPrefix(imgName, "breaker_player_") {
					return imgName, true
				}
			}
			return "breaker_player_", false
		},
	}

	for _, imgCheck := range requiredImages {
		if imgName, exists := imgCheck(); !exists {
			panic(fmt.Errorf("模板图片 '%s' 未找到，请检查图片目录: %s", imgName, imgPath))
		}
	}
	return imgMap, nil
}
