package functions

import (
	"github.com/767829413/easy-novel/pkg/utils"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

type checkUpdate struct {
	log   *logrus.Logger
	delay int
}

func NewCheckUpdate(l *logrus.Logger, delay int) App {
	return &checkUpdate{log: l, delay: delay}
}

func (c *checkUpdate) Execute() error {
	c.log.Info("Checking for updates")
	utils.GetColorIns(color.BgHiYellow).Println("检查更新...")
	// 实现检查更新的逻辑
	return nil
}
