package functions

import (
	"fmt"

	"github.com/767829413/easy-novel/pkg/utils"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

type exit struct {
	log *logrus.Logger
}

func NewExit(l *logrus.Logger) App {
	return &exit{log: l}
}

func (e *exit) Execute() error {
	utils.GetColorIns(color.BgHiGreen).Println("<== Bye :)")
	return fmt.Errorf("exit")
}
