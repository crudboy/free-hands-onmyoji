package utils

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/vcaesar/gcv"
)

// 把样本图片变成image.Image
func ReadPic(path string) image.Image {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	return img
}

func RandFromRangeInt64(min int64, max int64) int64 {
	// 使用带有时间种子的新随机源替代 rand.Seed
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Int63n(max-min) + min
}

// NormalFloat64 正态分布公式
func NormalFloat64(x int64, miu int64, sigma int64) float64 {
	randomNormal := 1 / (math.Sqrt(2*math.Pi) * float64(sigma)) * math.Pow(math.E, -math.Pow(float64(x-miu), 2)/(2*math.Pow(float64(sigma), 2)))
	//注意下是x-miu，我看网上好多写的是miu-miu是不对的
	return randomNormal
}

// RandomNormalInt64 正态分布随机数生产器：min:最小值，max:最大值，miu:期望值（均值），sigma:方差
func RandomNormalInt64(min int64, max int64, miu int64, sigma int64) (bool, int64) {
	if min >= max {
		return false, 0
	}
	if miu < min {
		miu = min
	}
	if miu > max {
		miu = max
	}
	var x int64
	var y, dScope float64
	for {
		x = RandFromRangeInt64(min, max)
		y = NormalFloat64(x, miu, sigma) * 100000
		dScope = float64(RandFromRangeInt64(0, int64(NormalFloat64(miu, miu, sigma)*100000)))
		//注意下传的是两个miu
		if dScope <= y {
			break
		}
	}
	return true, x
}

func Jpg2RGBA(img image.Image) *image.RGBA {
	tmp := image.NewRGBA(img.Bounds())

	draw.Draw(tmp, img.Bounds(), img, img.Bounds().Min, draw.Src)
	return tmp
}
func FindTempPos(temp, img image.Image) (int, int, float32) {
	//把image.Image统一转换成image.RGBA，不然会断言失败
	_, num, _, pos := gcv.FindImg(Jpg2RGBA(temp), Jpg2RGBA(img))
	return pos.X, pos.Y, num
}
func SaveImg(img image.Image, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create image file: %w", err)
	}
	defer f.Close()
	err = png.Encode(f, img)
	return nil
}

func RandomInt(i int, i2 int) int {
	if i >= i2 {
		return i
	}
	// 使用带有时间种子的新随机源替代 rand.Seed
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(i2-i+1) + i
}
