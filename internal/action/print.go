package action

import (
	"fmt"

	"github.com/767829413/easy-novel/internal/config"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

type printConf struct {
	log *logrus.Logger
}

func NewPrintConf(l *logrus.Logger) Action {
	return &printConf{log: l}
}

func (p *printConf) Execute() error {
	fmt.Println(config.GetConf().ToJSON())
	return nil
}

type printHint struct {
	log     *logrus.Logger
	version string
}

func NewPrintHint(l *logrus.Logger, version string) Action {
	return &printHint{log: l}
}

func (p *printHint) Execute() error {
	cfg := config.GetConf()
	color.Blue("easy-novel %s （本项目开源免费）", p.version)
	fmt.Println("开源地址：https://github.com/767829413/easy-novel")
	fmt.Printf("当前书源：%d\n", cfg.Base.SourceID)
	color.Cyan("导出格式：%s", cfg.Base.Extname)
	color.Yellow("请务必阅读 readme.txt")
	return nil
}
