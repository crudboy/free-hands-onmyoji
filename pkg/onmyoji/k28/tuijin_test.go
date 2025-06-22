package k28

import (
	"fmt"
	"free-hands-onmyoji/pkg/onmyoji/window"
	"free-hands-onmyoji/pkg/statemachine"
	"testing"

	"github.com/andybrewer/mack"
)

// MockTaskController 模拟TaskController接口
type MockTaskController struct {
	NextCalled      bool
	NextIndexCalled bool
	NextKey         string
}

func (m *MockTaskController) Next(key string) error {
	m.NextCalled = true
	m.NextKey = key
	return nil
}

func (m *MockTaskController) NextIndex() {
	m.NextIndexCalled = true
}

// 由于robotgo.CaptureImg以及robotgo.MoveClick在测试环境中不应该被真正执行，
// 因此这里使用临时的替代函数，并在测试后恢复原来的函数
func TestTuiJing_Execute(t *testing.T) {

	script := fmt.Sprintf(`
tell application "System Events"
 tell application process "%s"
  get {position, size} of window 1
 end tell
end tell
`, "BlueStacks")
	tell, err := mack.Tell("System Events", script)
	if err != nil {
		panic(err)
	}
	var x, y, w, h int
	fmt.Sscanf(tell, "%d, %d, %d, %d", &x, &y, &w, &h)
	// 创建TuiJing实例
	tuiJing := &TuiJing{
		Window: window.Window{
			WindowX: x,
			WindowY: y,
			WindowW: w,
			WindowH: h,
			OffsetX: 0,
			OffsetY: 0,
		},
	}
	machine := statemachine.NewStateMachine()
	machine.AddTask(tuiJing)
	for i := 0; i < 5; i++ {
		machine.Run()
	}

}
