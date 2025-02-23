# wxbot - 微信聊天机器人
> 适用于微信（WeChat **3.9.8.15** | 3.9.7.29）
> 可在Windows PC微信 **设置** - **关于微信** - **版本信息** 中获取您当前的微信版本，如果您当前的微信版本不在上述可用的版本列表中，请至下方 **3、可用版本微信安装包获取** 选择最新版微信重新安装使用

**未经过大量测试，使用远程线程注入方式可能会被报毒（无毒，请放心使用！），也可以尝试使用例如x64dbg等方式进行注入，注入手段并不重要，只要将wxbot.dll注入到wechat.exe中即可**

## 免责声明
本仓库发布的内容，仅用于学习研究，请勿用于非法用途和商业用途！如因此产生任何法律纠纷，均与作者无关！

## 1、运行
bin目录下有如下两个文件（仅在windows 10 & windows server 2012 R2系统上进行测试）：
* inject.exe (bin/inject.exe)
* wxbot.dll (bin/wxbot.dll)

运行的时请保证这两个文件在同一目录下，直接运行inject.exe即可（运行注入器前请保证微信已登陆！）
默认wxbot.dll为最新版（3.9.8.15），低版本微信注入请选择对应版本的wxbot-xxxx.dll替换为wxbot.dll后注入即可
**运行成功时微信会弹出注入成功弹窗！（http server在弹窗确认后启动）**

## 2、使用
### 2.1、配置文件
> 配置文件路径为：[微信根目录]/wxbot.json
> 配置文件为json格式，默认不自动创建！
```json
{
    "addr": "0.0.0.0:8080",
    "sync-url": [
        {
            "timeout": 3000,
            "url": "http://localhost:8081/callback"
        }
    ]
}
```
* **addr: wxbot服务监听地址（固定为ip:port形式）**
* **sync-url：http回调地址列表（建议通过下面的/sync-url接口修改，不要手动修改）**
  * **url：回调url**
  * **timeout：回调超时时间**

### 2.2、路由列表
> 响应信息
> 固定为JSON格式响应：{"code": 200, data: xxxx, "message": "xxx"}
> * code：固定200
> * message: 成功为success，失败为faild
> * data: 根据请求接口不同数据不同，无特别描述时下面的请求接口返回字段全部为该data字段的子字段

**路由列表概览：**
* **功能类**
  * /userinfo      - 获取登陆用户信息
  * /contacts      - 获取通讯录信息（wxid从这个接口获取）
  * /sendtxtmsg    - 发送文本消息
  * /sendimgmsg    - 发送图片消息（支持json和form-data表单上传两种方式，json方式请将二进制数据使用base64编码后发送）
  * /sendfilemsg   - 发送文件消息（支持json和form-data表单上传两种方式，json方式请将二进制数据使用base64编码后发送）

* **回调注册类（目前仅用来获取微信实时消息 - 同步消息接口，同时支持WebSocket和http两种方式！）**
  * /ws            - 注册websocket回调（支持注册多个ws通道）
  * /sync-url      - http回调相关（支持注册多个http接口，注册请带上协议头：http/https，注册成功会持久化到配置文件中）

#### 2.2.1、功能类接口
##### 2.2.1.1、登陆信息
**协议信息**
GET /userinfo
**响应字段**
* custom_account: 微信号
* nickname： 微信昵称
* phone： 手机号
* phone_system： 手机系统
* profile_picture： 头像
* profile_picture_small： 小头像
* wxid

##### 2.2.1.2、通讯录
**协议信息**
GET /contacts
**响应字段**
* custom_account： 微信号
* nickname： 昵称
* note： 备注
* pinyin： 昵称拼音首字母大写
* pinyinAll： 昵称拼音全
* type1： 用户类别1
* type2： 用户类别2
* wxid

##### 2.2.1.3、发送文本消息
**协议信息**
POST /sendtxtmsg
**请求字段**
* wxid
* content： 发送消息内容

**响应示例**
{"code":200,"msg":"success"}

##### 2.2.1.4、发送图片消息
**协议信息**
POST /sendimgmsg
> 支持JSON和form-data表单两种方式提交

**请求头**
* **JSON：`Content-Type: application/json`**
* **form-data表单：`Content-Type: multipart/form-data`**

**请求字段**
* **JSON：**
    * wxid
    * image： 图片二进制数据base64编码后字符串
* **form-data表单**
    符合标准`form-data`数据格式，需要两个参数分别是`wxid`和`image`

##### 2.2.1.5、发送文件消息
**协议信息**
POST /sendfilemsg
> 支持JSON和form-data表单两种方式提交

**请求头**
* **JSON：`Content-Type: application/json`**
* **form-data表单：`Content-Type: multipart/form-data`**

**请求字段**
* **JSON：**
    * wxid
    * file： 文件二进制数据base64编码后字符串
    * filename： 文件名
* **form-data表单**
    符合标准`form-data`数据格式，需要两个参数分别是`wxid`和`file`

#### 2.2.2、回调注册类
> 目前仅用来同步微信消息

**响应字段**
* wxid
* content： 消息内容
* toUser： 消息接收人（一般为登陆用户wxid）
* msgid： 消息唯一标识
* originMsg： 原始消息（如：wxid:\nxxxxxxx）
* chatRoomSourceWxid： 如果为群聊消息，则为消息发送人wxid
* msgSource： 加密消息
* type： 消息类型
* displayMsg： 展示消息（一般群聊消息下可用）

##### 2.2.2.1、websocket协议消息
**协议信息**
GET ws://xxxxx/ws

> websocket没什么好说的，基本上第三方库都有直接可用的实现，协议升级后就是一条全双工通道，目前只用来接收同步微信的实时消息，不要发送消息到服务端，服务端不会响应。

##### 2.2.2.2、http协议
> 需要你自己起一个Http Server服务用来接收微信的实时消息，你自己的Http Server启动之后通过接口注册到wxbot即可

**注册接口**
POST /sync-url
**请求字段**
* url： 你自己启动的Http Server地址路由（**ip:port/[subpath]**）
* timeout： 超时时间（当有一条新消息通过wxbot发送到你的回调地址时的最长连接等待时间）

### 2.3、接口使用例子
```bash
# 获取登陆用户信息
curl 127.0.0.1:8080/userinfo

# 获取通讯录信息
curl 127.0.0.1:8080/contacts

# 发送文本消息
curl -XPOST -d'{"wxid": "47331170911@chatroom", "content": "测试内容\nHello World"}' 127.0.0.1:8080/sendtxtmsg

# 发送图片消息1（使用form-data表单方式提交）
curl -XPOST -F "wxid=47331170911@chatroom" -F "image=@/home/jwping/1.jpg" 127.0.0.1:8080/sendimgmsg
# 发送图片消息2（使用json方式提交）
curl -XPOST -H "Content-Type: application/json" -d'{"wxid": "47331170911@chatroom", "image": "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"}' 127.0.0.1:8080/sendimgmsg

# 发送文件消息1（使用form-data表单方式提交）
curl -XPOST -F "wxid=47331170911@chatroom" -F "file=@/home/jwping/1.txt" 127.0.0.1:8080/sendfilemsg
# 发送文件消息2（使用json方式提交）
curl -XPOST -H "Content-Type: application/json" -d'{"wxid": "47331170911@chatroom", "filename": "1.txt", "file": "aGVsbG8gd29ybGQh"}' 127.0.0.1:8080/sendfilemsg


# 注册ws回调
# 使用任意程序websocket客户端连接127.0.0.1:8080/ws

# 注册http回调（http协议头不能少！）
curl -XPOST -d'{"url": "http://127.0.0.1:8081/callback", "timeout": 6000}' 127.0.0.1:8080/sync-url

# 获取当前已注册的http回调
curl 127.0.0.1:8080/sync-url

# 删除一个已注册的http回调
curl -XDELETE -d'{"url": "http://127.0.0.1:8081/callback"}' 10.1.12.12:8080/sync-url


# 同步消息回调响应例子（回调消息为JSON格式）
# 下面例子为反序列化后输出
# WebSocket Client Response
{Wxid:34418372934@chatroom Content:你好 ToUser:wxid_gotub49l54fq29 Msgid:7438040783824576403 OriginMsg:wxid_o4jinvsgz6lp31:
你好 ChatRoomSourceWxid:wxid_o4jinvsgz6lp31 MsgSource:<msgsource>
        <pua>1</pua>
        <silence>1</silence>
        <membercount>193</membercount>
        <signature>v1_8TXCDRkh</signature>
        <tmp_node>
                <publisher-id></publisher-id>
        </tmp_node>
</msgsource>
 Type:1 DisplayMsg:}
```

## 3、wxbox.dll、注入器、可用版本微信安装包等获取
* **阿里网盘：**
https://www.aliyundrive.com/s/4eiNnE4hp4n
提取码: rt25

* **百度网盘：**
https://pan.baidu.com/s/1cmzXe8AxYvzXWW2WTVCdxQ?pwd=l671 
提取码：l671

## 4、交流
请添加微信：**Anshan_PL**，备注 **wxbot** 拉微信交流群
