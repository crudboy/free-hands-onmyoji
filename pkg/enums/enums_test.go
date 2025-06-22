package enums

import "testing"

func Test_PrintType(t *testing.T) {
	// 打印所有任务类型
	for _, taskType := range []TaskType{TuiJing, XunGuai, JieSuan, Boss} {
		println(string(taskType))
	}
}
