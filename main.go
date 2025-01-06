package main

import (
	"PerformanceMonitor/pkg/model"
	"PerformanceMonitor/pkg/utils"
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var warnCount int

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic:", err)
		}
	}()
	viper.SetConfigFile("config.toml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	err = viper.Unmarshal(&model.Config)
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	c := cron.New()
	_, err = c.AddFunc(model.Config.CollectionIntervalCorn, func() {
		percent, err := cpu.Percent(time.Second, false)
		if err != nil {
			panic(fmt.Errorf("Read CPU fail: %s \n", err))
		}
		if percent[0] >= model.Config.WarnIndex {
			warnCount++
		}
		if warnCount == model.Config.AlarmCount {
			for i := 0; i < 3; i++ {
				res, msg := utils.SendMsg(utils.WebHookInfo{
					MsgType: "text",
					Text: utils.WebHookInfo2{
						Content: fmt.Sprintf("%s(%s) 已连续 %d%s CPU占用率超过 70%% 请尽快处理！", model.Config.ServerAlias, model.Config.ServerIp, model.Config.AlarmCount, model.Config.CollectionIntervalUnit),
					},
				})
				if res {
					break
				} else {
					fmt.Println("发送告警信息失败：", msg)
				}
			}
			warnCount = 0
		}
	})
	if err != nil {
		panic(fmt.Errorf("Fatal error add CronFunc: %s \n", err))
	}
	c.Start()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	c.Stop()
}
