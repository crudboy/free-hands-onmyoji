package capture

import (
	"free-hands-onmyoji/internal/onmyoji"
	"free-hands-onmyoji/internal/onmyoji/window"
	"free-hands-onmyoji/internal/statemachine"
	"free-hands-onmyoji/internal/tasks"
)

type CaptureAndSave struct {
	window.Window
}

func newCaptureAndSave(window window.Window, info onmyoji.ImgInfo, config onmyoji.Config) *CaptureAndSave {
	return &CaptureAndSave{
		Window: window,
	}
}

func (t *CaptureAndSave) Name() tasks.TaskType {
	return tasks.CaptureDisplay
}
func (t *CaptureAndSave) Execute(controller statemachine.TaskController) error {
	t.Window.CaptureAndSave()
	return nil
}
