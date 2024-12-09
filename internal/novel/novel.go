package novel

import (
	"context"
	"strings"

	"github.com/767829413/easy-novel/internal/functions"
	"github.com/767829413/easy-novel/pkg/utils"
	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

const Version = "v0.1"

// Run is the main entry point for the novel downloading logic
func Run(ctx context.Context, log *logrus.Logger) error {
	log.Info("Starting novel download process")
	options := []string{"1.下载小说", "2.检查更新", "3.查看配置文件", "4.使用须知", "5.结束程序"}
	actions := map[string]functions.App{
		"1": functions.NewDownload(log),
		"2": functions.NewCheckUpdate(log, 5000),
		"3": functions.NewPrintConf(log),
		"4": functions.NewPrintHint(log, Version),
		"5": functions.NewExit(log),
	}

	var completerItems []readline.PrefixCompleterInterface
	for _, option := range options {
		completerItems = append(completerItems, readline.PcItem(option))
	}

	completer := readline.NewPrefixCompleter(completerItems...)

	rl, err := readline.NewEx(&readline.Config{
		Prompt:       "按 Tab 键选择功能: ",
		AutoComplete: completer,
	})

	if err != nil {
		return err
	}
	defer rl.Close()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			line, err := rl.Readline()
			if err != nil {
				return err
			}
			cmd := strings.TrimSpace(line)
			index := strings.Split(cmd, ".")[0]
			action, found := actions[index]
			if !found {
				utils.GetColorIns(color.FgHiRed).Println("无效的选项，请重新选择")
				continue
			}

			err = action.Execute()
			if err != nil {
				if err.Error() == "exit" {
					return nil
				}
				utils.GetColorIns(color.FgHiRed).Printf("执行操作时发生错误: %v\n", err)
			}
		}
	}
}
