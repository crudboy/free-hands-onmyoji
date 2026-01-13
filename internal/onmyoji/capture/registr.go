package capture

import (
	"free-hands-onmyoji/internal/onmyoji"
	"free-hands-onmyoji/internal/onmyoji/window"
	"free-hands-onmyoji/internal/statemachine"
	"free-hands-onmyoji/internal/tasks"
)

type Registrator struct {
}

func (r Registrator) Registration(machine *statemachine.StateMachine, w window.Window, config onmyoji.Config, imgMap map[string]onmyoji.ImgInfo) error {

	onmyoji.Registration(machine, newCaptureAndSave(w, imgMap[string(tasks.CaptureDisplay)], config))

	return nil
}

func (r Registrator) LoadImageTemplates() (map[string]onmyoji.ImgInfo, error) {
	return nil, nil
}
