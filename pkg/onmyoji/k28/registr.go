package k28

import (
	"fmt"
	"free-hands-onmyoji/pkg/enums"
	"free-hands-onmyoji/pkg/logger"
	"free-hands-onmyoji/pkg/onmyoji/entity"
	"free-hands-onmyoji/pkg/onmyoji/window"
	"free-hands-onmyoji/pkg/statemachine"
	"free-hands-onmyoji/pkg/utils"
	"image"
	"os"
	"path/filepath"
)

var templateMap = make(map[string]entity.ImgInfo)

func initImages() {
	// 如果模板图片已经加载，则不需要重新加载
	if len(templateMap) > 0 {
		return
	}

	imgPath := "./k28/"
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
		func(r *os.File) {
			defer r.Close()

			if r.Name() == filepath.Join(imgPath, ".DS_Store") {
				return
			}
			im, _, err := image.DecodeConfig(r)
			if err != nil {
				logger.Error("图片解码错误 %s: %v", imgFile.Name(), err)
				return // 跳过错误图片而不是终止程序
			}
			// 不要文件后缀
			fileName := imgFile.Name()[:len(imgFile.Name())-len(filepath.Ext(imgFile.Name()))]

			templateMap[fileName] = entity.ImgInfo{
				Path:    imgPath + imgFile.Name(),
				ImgMaxX: im.Width,
				ImgMaxY: im.Height,
				Image:   utils.ReadPic(imgPath + imgFile.Name()),
			}
		}(reader)
	}

	// 检查是否成功加载了所有必要的图片
	requiredImages := []string{
		string(enums.ZhangJie),
		string(enums.JinRu),
		string(enums.XunGuai),
		string(enums.JieSuan),
		string(enums.Boss),
		string(enums.BaoXiang),
	}

	for _, imgName := range requiredImages {
		if _, exists := templateMap[imgName]; !exists {
			panic(fmt.Errorf("模板图片 '%s' 未找到，请检查图片目录: %s", imgName, imgPath))
		}
	}
}
func Registration(machine *statemachine.StateMachine, info entity.WindowInfo) {
	initImages() // 初始化模板图片
	_registration(machine, zhangjie(info, templateMap[string(enums.ZhangJie)]))
	_registration(machine, jinru(info, templateMap[string(enums.JinRu)]))
	_registration(machine, xunguai(info, templateMap[string(enums.XunGuai)]))
	_registration(machine, tuijing(info))
	_registration(machine, jieSuan(info, templateMap[string(enums.JieSuan)]))
	_registration(machine, boss(info, templateMap[string(enums.Boss)]))
	_registration(machine, baoxiang(info, templateMap[string(enums.BaoXiang)]))
}

func _registration(machine *statemachine.StateMachine, task statemachine.NamedTask) {
	machine.AddTask(task)
}
func xunguai(info entity.WindowInfo, templateImg entity.ImgInfo) *XunGuai {
	var xunguai = &XunGuai{
		TemplateImg: templateImg,
		Window: window.Window{
			WindowX: info.WindowX,
			WindowY: info.WindowY,
			WindowH: info.WindowH,
			WindowW: info.WindowW,
		},
	}
	return xunguai
}
func tuijing(info entity.WindowInfo) *TuiJing {
	var tuijing = &TuiJing{
		Window: window.Window{
			WindowX: info.WindowX,
			WindowY: info.WindowY,
			WindowH: info.WindowH,
			WindowW: info.WindowW,
			OffsetX: 0, // 可选偏移调整
			OffsetY: 0, // 可选偏移调整
		},
	}
	return tuijing
}
func jieSuan(info entity.WindowInfo, templateImg entity.ImgInfo) *JieSuan {
	var jieSuan = &JieSuan{
		TemplateImg: templateImg,
		Window: window.Window{
			WindowX: info.WindowX,
			WindowY: info.WindowY,
			WindowH: info.WindowH,
			WindowW: info.WindowW,
		},
	}
	return jieSuan
}
func boss(info entity.WindowInfo, templateImg entity.ImgInfo) *Boss {
	var boss = &Boss{
		TemplateImg: templateImg,
		Window: window.Window{
			WindowX: info.WindowX,
			WindowY: info.WindowY,
			WindowH: info.WindowH,
			WindowW: info.WindowW,
		},
	}
	return boss

}
func baoxiang(info entity.WindowInfo, templateImg entity.ImgInfo) *BaoXiang {
	var baoxiang = &BaoXiang{
		TemplateImg: templateImg,
		Window: window.Window{
			WindowX: info.WindowX,
			WindowY: info.WindowY,
			WindowH: info.WindowH,
			WindowW: info.WindowW,
		},
	}
	return baoxiang

}
func zhangjie(info entity.WindowInfo, templateImg entity.ImgInfo) *ZhangJie {
	var zhangjie = &ZhangJie{
		TemplateImg: templateImg,
		Window: window.Window{
			WindowX: info.WindowX,
			WindowY: info.WindowY,
			WindowH: info.WindowH,
			WindowW: info.WindowW,
		},
	}
	return zhangjie
}
func jinru(info entity.WindowInfo, templateImg entity.ImgInfo) *JinRu {
	var jinru = &JinRu{
		TemplateImg: templateImg,
		Window: window.Window{
			WindowX: info.WindowX,
			WindowY: info.WindowY,
			WindowH: info.WindowH,
			WindowW: info.WindowW,
		},
	}
	return jinru
}
