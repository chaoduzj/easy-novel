package functions

import (
	"encoding/json"
	"time"

	"github.com/767829413/easy-novel/internal/version"
	"github.com/767829413/easy-novel/pkg/utils"
	"github.com/Masterminds/semver/v3"
	"github.com/fatih/color"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/sirupsen/logrus"
)

const (
	RELEASE_URL = "https://api.github.com/repos/767829413/easy-novel/releases"
)

type checkUpdate struct {
	log          *logrus.Logger
	timeoutMills int
}

func NewCheckUpdate(l *logrus.Logger, timeoutMills int) App {
	return &checkUpdate{log: l, timeoutMills: timeoutMills}
}

func (c *checkUpdate) Execute() error {
	c.log.Info("Checking for updates")
	utils.GetColorIns(color.BgHiYellow).Println("检查更新...")
	// 实现检查更新的逻辑
	collector := colly.NewCollector(
		colly.Async(true),
	)
	extensions.RandomUserAgent(collector)
	collector.SetRequestTimeout(time.Duration(c.timeoutMills) * time.Millisecond)

	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 1,
	})
	var latestVersion, latestUrl string
	collector.OnResponse(func(r *colly.Response) {
		var releases []map[string]interface{}
		err := json.Unmarshal(r.Body, &releases)
		if err != nil || len(releases) == 0 {
			utils.GetColorIns(color.FgHiRed).Println("<== 无法获取版本信息")
			return
		}

		latest := releases[0]
		currentVersion := version.Version
		latestVersion = latest["tag_name"].(string)
		latestUrl = latest["html_url"].(string)

		v1, err := semver.NewVersion(currentVersion)
		if err != nil {
			c.log.Errorf("Error parsing current version: %s", err)
			return
		}

		v2, err := semver.NewVersion(latestVersion)
		if err != nil {
			c.log.Errorf("Error parsing latest version: %s\n", err)
			return
		}

		if v2.GreaterThan(v1) {
			utils.GetColorIns(color.FgHiYellow).
				Printf("<== 发现新版本: %s (%s)\n", latestVersion, latestUrl)
		} else {
			utils.GetColorIns(color.FgHiGreen).
				Printf("<== %s 已是最新版本！(%s)\n", latestVersion, latestUrl)
		}
	})

	collector.OnError(func(r *colly.Response, err error) {
		utils.GetColorIns(color.FgHiRed).
			Printf("<== 检查失败，当前网络环境暂时无法访问 GitHub，请稍后再试 (%s)\n", err.Error())
	})

	collector.Visit(RELEASE_URL)
	collector.Wait()
	return nil
}
