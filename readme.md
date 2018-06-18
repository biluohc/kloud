# kloud 基于 Go 语言 和 PostgreSQL 实现的网络硬盘（云盘）

本系统采用前后端分离架构

## 后端提供 API
1. [Echo](github.com/labstack/echo)： Web 框架
2. [go-pg/pg](github.com/go-pg/pg)： ORM
3. [NanoId](github.com/matoous/go-nanoid)： UUID替代品 

## 前端提供用户界面
1. [Vue.js](https://github.com/vuejs/vue)：  前端框架
1. [iView](https://github.com/iview/iview)：组件库
2. [axios](https://github.com/axios/axios)： AJAX
3. [crypto-js](https://github.com/brix/crypto-js)：sha256
3. [clipboard.js](https://github.com/zenorocha/clipboard.js)：clipboard

## 数据库规划与文件存储
按 sha1 存储文件，使用数据表模拟用户目录。

## 功能与模块简介
2. 用户模块：登录、注册、注销、用户个人设置
1. 用户文件管理：文件目录浏览、新建目录、上传文件、下载文件、重命名、移动、删除到回收站、  
分享文件（可选密码、过期时间）、  
公开文件（可选 Referer、过期时间，可以作为图床使用） 
3. 用户回收站：恢复、彻底删除、清空
3. 用户分享文件管理：编辑分享（密码、过期时间）、删除
4. 用户公开文件管理：编辑公开（Referer、过期时间）、删除
5. 后台系统概况：日志查看等
5. 后台用户管理：封禁解封用户、编辑用户、新建用户等
5. 后台文件管理：查看、封禁解封文件

## 缺点
1. 前端水平糟糕，未适配移动端
2. 后台功能较为薄弱，有待加强
3. 未分离出配置文件，部署需要自己编译
3. 另外还有诸多小问题，主要是报错的友好度和前端界面的细节等。

## 未来
因为按带宽计算的服务器网络老贵，按流量又怕 DDOS 卖身, 觉意义不大, 遂暂停开发。

所以不建议使用，仅供学习参考。

## 部署 
1. 前端配置：`Afront/src/api.js` 的带协议服务器地址 `host`
2. 后端配置：`config.go` 里的 数据库配置（`dbConfig`）、服务器监听地址（`ServeAddr`）、带协议服务器地址（`Host`）、https则需要修改那几个 https 相关的变量
3. 数据库配置，导入 `pg.sql`
3. 将项目目录移动到 `GOPATH`下，在该`GOPATH下`新建 `src/github.com/biluohc/` 并且将本项目移动到下面（傻逼的模块与包管理）。
3. 下载项目后端相关依赖，注意 `golang/x`下的fq或者手动建立目录并`git clone` [github.com/golang](github.com/golang)对应的包镜像 后 `go build`
4. 在`Afront`目录下运行`yarn`下载前端依赖, `yarn build` 编译
5. 直接拷贝整个项目到目标机器，`env GO_LOG=2 nohup  ./kloud  &` 运行即可


### 参考的

1. [蓝眼云盘](https://github.com/eyebluecn/tank)
2. [vue-cnode](https://github.com/xjh22222228/vue-cnode)
2. 百度云盘

