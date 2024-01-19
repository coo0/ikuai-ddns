# iKuai ddns
爱快外网ip绑定到dynv6,dedyn

## 如何使用

[//]: # (1. 编写配置文件，命名为`config.yml`)

[//]: # (2. 进入爱快自带 docker 中，点击`镜像管理`->`添加`，选择`镜像库下载`，搜索`ztc1997/ikuai-bypass`，下载`TAG`为`latest`的镜像)

[//]: # (3. 点击`容器列表`->`添加`，`选择镜像文件`选择`ztc1997/ikuai-dynv6:latest`，打开`高级选项`，添加一个`挂载目录`，`源路径`填写放置配置文件的路径，`目标路径`填写`/etc/ikuai-bypass`，内存64M即可，其它根据需要自行填写)

[//]: # (4. 保存并启用)

### 配置文件模板

```yaml
## 爱快管理页面的 URL，结尾不要加 "/"，
ikuai-url: https://192.169.0.1
username: admin # 爱快用户名
password: pwd  # 爱快密码
update-api: [https://dynv6.com/api/update]
hostname: [[],[]]
Token: []
cron: 0 6 * * * # 执行更新的周期 crontab
crondel: 0 6 * * * # 执行清空更新的ip周期 crontab
```
