package window

import (
	"github.com/go-vgo/robotgo"
)

func AlertNotify(title string, message string) {
	robotgo.Alert(title, message)

}
