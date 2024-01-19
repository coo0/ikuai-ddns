package api

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const FuncNameEtherInfo = "homepage"

type EtherInfoData struct {
	SnapshootWan []WAN `json:"snapshoot_wan"`
}
type WAN struct {
	Id        int    `json:"id"`
	Errmsg    string `json:"errmsg"`
	Gateway   string `json:"gateway"`
	Interface string `json:"interface"`
	IpAddr    string `json:"ip_addr"`
}

func (i *IKuai) ShowEtherInfoByComment(url []string, hostname [][]string, token []string) error {
	param := struct {
		Type string `json:"TYPE"`
	}{
		Type: "ether_info,snapshoot",
	}
	req := CallReq{
		FuncName: FuncNameEtherInfo,
		Action:   "show",
		Param:    &param,
	}
	result := EtherInfoData{}
	resp := CallResp{Data: &result}
	err := postJson(i.client, i.baseurl+"/Action/call", &req, &resp)
	if err != nil {
		return err
	}
	if resp.Result != 30000 {
		return errors.New(resp.ErrMsg)

	}
	IpAddr := ""
	for _, wan := range result.SnapshootWan {
		if wan.IpAddr != "" {
			IpAddr = wan.IpAddr
			break
		}
	}
	log.Println("外网ip：" + IpAddr)
	filePath := "./ip_tmp/" + IpAddr
	exists, err := fileExists(filePath)
	if err != nil {
		log.Printf("检查过程中发生错误: %v", err)
	} else if exists {
		log.Println("文件存在")
	} else {
		file, err := os.Create(filePath)
		if err != nil {
			log.Println("创建文件时发生错误:", err)
			return err
		}
		defer file.Close()
		for i, u := range url {
			token := token[i]
			for _, hname := range hostname[i] {
				urlStr := ""
				//判断字符串中是否包含字符
				if strings.Contains(u, "dynv6.com") {
					urlStr = u + "?hostname=" + hname + "&token=" + token + "&ipv4=" + IpAddr
				} else if strings.Contains(u, "dedyn.io") {
					urlStr = u + "?hostname=" + hname + "&myipv4=" + IpAddr
				}
				if urlStr == "" {
					continue
				}
				fmt.Println(urlStr)
				client := &http.Client{}
				req, err := http.NewRequest("GET", urlStr, nil)
				if err != nil {
					fmt.Println(1)
					log.Fatal(err)
				}
				if strings.Contains(u, "dedyn.io") {
					req.Header.Set("Authorization", "Token "+token)
				}

				resp, err := client.Do(req)
				if err != nil {
					fmt.Println(2)
					log.Fatal(err)
				}
				defer resp.Body.Close()
				bodyText, err := io.ReadAll(resp.Body)
				if err != nil {
					fmt.Println(4)
					log.Println(err)
				} else {
					log.Println(hname + "绑定到" + IpAddr + ":" + string(bodyText))
				}
				time.Sleep(time.Second * 3)

			}
			time.Sleep(time.Second * 1)
		}
	}
	return nil
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil // 文件存在且无错误
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil // 文件不存在
	}
	return false, err // 其他错误（比如权限问题）
}
