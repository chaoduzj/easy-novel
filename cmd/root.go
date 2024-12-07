package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/767829413/easy-novel/internal/config"
	"github.com/767829413/easy-novel/internal/novel"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cfgFile string
var log = logrus.New()
var rootCmd = &cobra.Command{
	Use:   "easy-novel",
	Short: "一个下载网络小说的简单工具",
	Long:  `主要是通过搜索互联网上全本的小说,根据自己的喜好获取爱看的小说!`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		// Setup logger
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
		log.SetOutput(os.Stdout)

		// Set log level
		log.SetLevel(logrus.ErrorLevel)

		// Load config here
		err = config.LoadConfig(cfgFile)
		if err != nil {
			return fmt.Errorf("error loading config: %w", err)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetConf()
		log.Info("Welcome to easy-novel!")
		log.WithFields(logrus.Fields{
			"sourceID":     cfg.Base.SourceID,
			"downloadPath": cfg.Base.DownloadPath,
		}).Info("Configuration loaded")

		// Set log level based on configuration
		logLevel, err := logrus.ParseLevel(cfg.Base.LogLevel)
		if err != nil {
			log.WithError(err).Warn("Invalid log level in config, defaulting to Info")
			logLevel = logrus.ErrorLevel
		}
		log.SetLevel(logLevel)

		// Create a context that we can cancel
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Handle OS signals
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-sigCh
			log.Info("Received interrupt signal. Shutting down...")
			cancel()
		}()

		// Call the main logic entry point
		if err := novel.Run(ctx, log); err != nil {
			log.WithError(err).Error("Failed to run easy-novel")
			os.Exit(1)
		}
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().
		StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.easy-novel.yaml)")

}
