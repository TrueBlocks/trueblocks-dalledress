package msgs

import (
	"fmt"
)

func EmitStatus(message string) {
	EmitMessage(EventStatus, message)
}

func EmitError(source string, err error) {
	if err == nil {
		return
	}

	message := fmt.Sprintf("Error: %s - %s", source, err.Error())
	EmitMessage(EventError, message)
}

func EmitWarning(source string, warning string) {
	message := fmt.Sprintf("Warning: %s - %s", source, warning)
	EmitMessage(EventError, message)
}

func EmitManager(reason string) {
	EmitMessage(EventManager, reason)
}

func EmitAppInit() {
	EmitMessage(EventAppInit, "")
}

func EmitAppReady() {
	EmitMessage(EventAppReady, "")
}

func EmitViewChange(view string) {
	EmitMessage(EventViewChange, view)
}
