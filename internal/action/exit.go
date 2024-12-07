package action

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

type exit struct {
	log *logrus.Logger
}

func NewExit(l *logrus.Logger) Action {
	return &exit{log: l}
}

func (e *exit) Execute() error {
	color.Green("<== Bye :)")
	return fmt.Errorf("exit")
}
