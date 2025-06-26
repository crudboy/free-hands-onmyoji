package onmyoji

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

const defaultConfig = `
# 阴阳师自动化配置文件 默认配置
[k28]
# 结算阀值：多少次结算后进入下一个状态
level_completion_threshold = 7

# 刷几次小怪后尝试检查Boss的次数
after_minion_attempts_boss_check = 3
# 匹配章节次数匹配不成功会进入下一个环节
zhangjie_find_threshold = 2
# 匹配小怪匹配不成功会进入移动环节
xunguai_find_threshold = 2
chest_wait_time = 500 # 宝箱寻找完成后等待页面切换时间，单位毫秒
level_completion_wait_time = 600 # 结算等待时间，单位毫秒
`
const configPath = "./config.toml"

// Config 整体配置结构
type Config struct {
	K28     K28Config     `toml:"k28"`
	Breaker BreakerConfig `toml:"breaker"` // 添加Breaker配置
}

// K28Config 御魂副本相关配置
type K28Config struct {
	LevelCompletionThreshold     int `toml:"level_completion_threshold"`       // 结算阀值
	AfterMinionAttemptsBossCheck int `toml:"after_minion_attempts_boss_check"` // 刷几次小怪后尝试检查Boss的次数
	ZhangjieFindThreshold        int `toml:"zhangjie_find_threshold"`          // 章节任务寻找阈值
	XunguaiFindThreshold         int `toml:"xunguai_find_threshold"`           // 寻怪任务寻找阈值
	ChestWaitTime                int `toml:"chest_wait_time"`                  // 宝箱寻找完成后等待页面切换时间，单位毫秒
	LevelCompletionWaitTime      int `toml:"level_completion_wait_time"`       // 结算等待时间，单位毫秒
}
type BreakerConfig struct {
}

// NewDefaultConfig 创建一个带有默认值的Config实例
func NewDefaultConfig() Config {
	return Config{
		K28: K28Config{
			LevelCompletionThreshold:     7,   // 默认结算阀值为7
			AfterMinionAttemptsBossCheck: 3,   // 默认刷3次小怪后检查Boss
			ZhangjieFindThreshold:        2,   // 默认章节任务寻找阈值为2
			XunguaiFindThreshold:         2,   // 默认寻怪任务寻找阈值为2
			ChestWaitTime:                500, // 默认宝箱等待时间为500毫秒
			LevelCompletionWaitTime:      600, // 默认结算等待时间为600毫秒

		},
	}
}

// LoadConfig 从指定路径的TOML文件加载配置
// 如果文件不存在，将自动创建带有默认配置的文件
// 返回配置对象和可能的错误
// 即使发生错误，也会返回有效配置（默认配置或部分加载的配置）
// 调用者可以选择：1) 检查error并处理错误 2) 忽略error直接使用返回的配置
func LoadConfig(configPath string) (Config, error) {
	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 文件不存在，创建默认配置文件
		defaultCfg := NewDefaultConfig()
		err := os.WriteFile(configPath, []byte(defaultConfig), 0644)
		if err != nil {
			// 创建失败时返回默认配置和错误
			return defaultCfg, fmt.Errorf("创建默认配置文件失败: %v", err)
		}
		// 创建成功返回默认配置，无错误
		return defaultCfg, nil
	}

	// 读取配置文件
	var config Config

	// 解析TOML文件
	_, err := toml.DecodeFile(configPath, &config)
	if err != nil {
		// 解析失败时返回默认配置和错误
		return NewDefaultConfig(), fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 加载成功，返回加载的配置，无错误
	return config, nil
}

// LoadDefaultConfig 加载默认路径的配置文件，如果不存在则创建
// 返回配置对象和可能的错误
// 即使发生错误，也会返回有效配置（默认配置或部分加载的配置）
// 调用者可以选择：1) 检查error并处理错误 2) 忽略error直接使用返回的配置
func LoadDefaultConfig() (Config, error) {
	return LoadConfig(configPath)
}
