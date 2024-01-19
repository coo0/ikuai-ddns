package main

import (
	"flag"
	"fmt"
	"github.com/coo0/ikuai-ddns/api"
	"github.com/robfig/cron/v3"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"os/signal"
	"path/filepath"
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
	CronDel   string     `yaml:"crondel"`
}

func main() {
	flag.Parse()

	err := readConf(*confPath)
	if err != nil {
		log.Println("读取配置文件失败：", err)
		return
	}
	//判断目录是或否存在，不存在就创建
	err = ensureDirectoryExist()
	if err != nil {
		log.Println("生成ip文件夹失败：", err)
		return
	}
	//deleteFilesInDirectory()

	update() //绑定ip到dynv6

	if conf.Cron == "" {
		return
	}

	c := cron.New()
	_, err = c.AddFunc(conf.CronDel, deleteFilesInDirectory)
	if err != nil {
		log.Println("启动定时清空文件夹计划任务失败：", err)
		return
	} else {
		log.Println("已启动定时清空文件夹计划任务")
	}
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

func ensureDirectoryExist() error {
	// 检查目录是否存在
	path := "./ip_tmp"
	_, err := os.Stat(path)
	if err == nil {
		// 目录已存在，无需操作
		return nil
	}
	if os.IsNotExist(err) {
		// 目录不存在，尝试创建
		err = os.MkdirAll(path, 0755)
		if err != nil {
			log.Printf("创建目录失败: %s\n", err)
		}
		return nil
	}
	// 其他错误情况，直接返回错误信息
	return err
}

// 删除文件夹下的所有文件
func deleteFilesInDirectory() {
	// 打开目录
	dirPath := "./ip_tmp"
	d, err := os.Open(dirPath)
	if err != nil {
		log.Printf("打开目录失败: %s\n", err)
	}
	defer d.Close()

	// 遍历目录中的所有元素
	files, err := d.Readdir(-1)
	if err != nil {
		log.Printf("读取目录内容失败: %s", err)
	}

	// 遍历并删除每个元素
	for _, fileInfo := range files {
		fullPath := filepath.Join(dirPath, fileInfo.Name())
		err = os.Remove(fullPath)
		if err != nil {
			log.Printf("删除文件失败: %s", err)
		}
	}

}
