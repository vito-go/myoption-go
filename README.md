# myoption

产品介绍：
本产品以上证指数或股票价格为标的，用于投注上证指数的价格。用户可选择下注金币数量（10-2000金币），以及不同场次（2分钟、3分钟、5分钟、10分钟、20分钟、30分钟、60分钟、全天）。投注方向有两个选项：看涨和看跌。

举例说明：当前时间为10:01:00，实时价格指数为3230.64。用户选择5分钟场次，看涨，并下注10金币。

到10:06:00，如果价格高于3230.64（如3233.56），用户盈利10金币；如果价格小于等于3230.64，用户亏损10金币。"

# 体验地址

- 网页在线: https://vitogo.tpddns.cn:9131/web
- 安卓下载链接: http://vitogo.tpddns.cn:9130/web/myoption.apk

# Quick Start

## 初始化数据库(需要安装postgresql)

```bash
# 初始化https证书, 生成的证书在configs目录的keys目录下
$ make initkey
# 初始化数据库
$ make initdb
```

## 编译静态文件

```bash
mkdir -p www/test/
git clone https://github.com/vito-go/myoption.git
cd myoption
make build
# 链接静态文件到本项目web目录: 例如 ~/go/srv/myoption/www/test/web
ln -s $PWD/build/web  ~/go/srv/myoption/www/test/web
```

## 2. 启动服务

```bash
$ make start
```

## 3. 在浏览器中打开(信任证书，即可访问。证书可在configs/keys中找到，或者自己生成)

```bash
https://localhost:9131/web
```

安卓下载链接: http://localhost:9130/web/myoption.apk

## 预览

<div>
<img  src="./images/home.png" style="width: 32%">
<img  src="./images/my.png" style="width: 32%">
<img  src="./images/transaction_detail.png" style="width: 32%">
</div>
<div>
<img  src="./images/wallet_detail.png" style="width: 32%">
<img  src="./images/title_order.png" style="width: 32%">
<img  src="./images/chart.png" style="width: 32%">
<img  src="./images/about.png" style="width: 32%">
</div>
