package main

import (
	"flag"
	"fmt"
	"github.com/robfig/cron/v3"
	"gopkg.in/yaml.v3"
	"ikuai-ddns/api"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var confPath = flag.String("c", "./config.yml", "配置文件路径")

var conf struct {
	IkuaiURL  string     `yaml:"ikuai-url"`
	Username  string     `yaml:"username"`
	Password  string     `yaml:"password"`
	UpdateApi []string   `yaml:"update-api"`
	Hostname  [][]string `yaml:"hostname"`
	Token     []string   `yaml:"Token"`
	Cron      string     `yaml:"cron"`
}

func main() {
	flag.Parse()

	err := readConf(*confPath)
	if err != nil {
		log.Println("读取配置文件失败：", err)
		return
	}

	update() //绑定ip到dynv6

	if conf.Cron == "" {
		return
	}

	c := cron.New()

	_, err = c.AddFunc(conf.Cron, update)
	if err != nil {
		log.Println("启动计划任务失败：", err)
		return
	} else {
		log.Println("已启动计划任务")
	}
	c.Start()

	{
		osSignals := make(chan os.Signal, 1)
		signal.Notify(osSignals, os.Interrupt, os.Kill, syscall.SIGTERM)
		<-osSignals
	}
}

func readConf(filename string) error {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(buf, &conf)
	if err != nil {
		return fmt.Errorf("in file %q: %v", filename, err)
	}
	return nil
}

func update() {
	err := readConf(*confPath)
	if err != nil {
		log.Println("更新配置文件失败：", err)
		return
	}

	baseurl := conf.IkuaiURL

	iKuai := api.NewIKuai(baseurl)

	err = iKuai.Login(conf.Username, conf.Password)
	if err != nil {
		log.Println("登陆失败：", err)
		return
	} else {
		log.Println("登录成功")
	}

	err = iKuai.ShowEtherInfoByComment(conf.UpdateApi, conf.Hostname, conf.Token)

	if err != nil {
		log.Println("运行失败：", err)
		return
	} else {
		log.Println("运行成功")
	}
}
