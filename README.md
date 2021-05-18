# ToadOCRTools

ToadOCRTools是ToadOCR的Http AIP

## 一、需求

### 1、需求背景

1. **解决什么问题**

   当前存在可供使用的ToadOCR RPC服务集群，并开放了可供调用的RPC接口（详见[接口IDL](https://github.com/suvvm/ToadOCRPreprocessor/blob/master/rpc/idl/toad_ocr_preprocessor.proto)），为了方便演示图像处理效果，需要提供拥有GUI的客户端。

2. **可覆盖场景**

   - 答辩演示
   - 用户个人操作与商用

3. **为什么这么做**

   - 若调用ToadOCR服务当前只能通过RPC一种方式，需要开放其他的调用手段
   - 答辩演示只使用[doc文档站]()与rpc服务集群效果不佳

   - 如上方IDL所示RPC接口的调用需要认证AppID与AppSecret，当前没有构造AppID与AppSecret认证信息的方式，故需要提供构造工具

4. **潜在迭代方向**

   - RPC现有接口IDL可能发生变动
   - RPC服务可能新增接口

### 2、需求内容

![img1](https://raw.githubusercontent.com/suvvm/ToadOCRTools/master/resources/images/toadhttptools.png)

- 选择具体身份认证方式，提供设计方案
- 服务端
  - 提供整体架构概述
  - 提供接口idl
- GUI客户端
  - 提供GUI页面构造图组
    - 整体页面关系结构图
    - 单个页面功能概括（实现功能/调用接口）
    - 单个页面总体布局设计
  - log 内容与存储方式
  - 数据埋点

## 二、框架选型与存储设计

- http服务端使用gin
- RPC使用gRPC
- 持久化数据将以KV的形式缓存在etcd集群中，并在DB中留档

### 数据字典

#### app_infos

| 字段   | 数据类型     | 是否为空 | 默认值         | 描述                         |
| ------ | ------------ | -------- | -------------- | ---------------------------- |
| id     | int          | not null | auto increment | appID， 主键，自增 起始95501 |
| secret | varchar(200) | not null |                | appSecret                    |
| email  | varchar(50)  | not null | "empty"        | 邮箱                         |
| p_num  | varchar(11)  | not null | "empty"        | 手机号                       |

## 三、接口信息

### 1、接口总览

```go
r.Any("/toad_ocr/ping", handler.Pong)

r.POST("/toad_ocr/process", handler.Process)
r.POST("/toad_ocr/process/v2", handler.ProcessV2)

r.POST("/toad_ocr/send/sms", handler.Sms)
r.POST("/toad_ocr/send/email", handler.Email)

r.POST("/toad_ocr/application", handler.ApplicationAdd)
r.DELETE("/toad_ocr/application", handler.ApplicationDel)
r.GET("/toad_ocr/application", handler.ApplicationGet)
```

### 2、验证信息发送

- 发送手机验证码

  ```yaml
  Path: /toad_ocr/send/sms
  Method: POST
  Headers:
  	Content-Type:application/json
  UrlParam:
  BodyParam:
  	p_num(string)		# 目标手机号
  	code(string)		# 要发送的验证码
  response:
  	code(int)				# 错误码，无错误为0
  	message(string) # 错误信息，无错误为success
  ```

- 发送邮箱验证码

  ```yaml
  Path: /toad_ocr/send/email
  Method: POST
  Headers:
  	Content-Type:application/json
  UrlParam:
  BodyParam:
  	email(string)		# 目标邮箱
  	code(string)		# 要发送的验证码
  response:
  	code(int)				# 错误码，无错误为0
  	message(string) # 错误信息，无错误为success
  ```

### 3、身份认证接口

- 添加认证信息（需要客户端与服务端双验证）

  ```yaml
  Path: /toad_ocr/application
  Method: POST
  Headers:
  	Content-Type:application/json
  UrlParam:
  BodyParam:
  	p_num(string)		# 目标手机
  	email(string)		# 目标邮箱（手机和邮箱可只提供一个）
  	user_verify_code(string)   # 用户输入的验证码
  	client_verify_code(string) # 客户端生成的验证码
  	#（客户端校验完成后服务端再次校验，user_verify_code和client_verify_code一致就行）
  response:
  	code(int)				  # 错误码，无错误为0
  	message(string)   # 错误信息，无错误为success
  	app_info(object)	# 认证信息
  		id(int)					# AppID
  		secret(string)  # AppSecret
  		email(string)   # 邮箱
  		p_num(string)   # 手机号
  ```

- 查询认证信息（需要客户端验证）

  ```yaml
  Path: /toad_ocr/application
  Method: GET
  Headers:
  	Content-Type:application/json
  UrlParam:
  	p_num(string)		# 目标手机
  	email(string)		# 目标邮箱（手机和邮箱可只提供一个）
  BodyParam:
  response:
  	code(int)				  # 错误码，无错误为0
  	message(string)   # 错误信息，无错误为success
  	app_info(object)	# 认证信息
  		id(int)					# AppID
  		secret(string)  # AppSecret
  		email(string)   # 邮箱
  		p_num(string)   # 手机号
  ```

- 注销认证信息（需要客户端与服务端双验证）

  ```yaml
  Path: /toad_ocr/application
  Method: DELETE
  Headers:
  	Content-Type:application/json
  UrlParam:
  BodyParam:
  	p_num(string)		# 目标手机
  	email(string)		# 目标邮箱（手机和邮箱可只提供一个）
  	user_verify_code(string)   # 用户输入的验证码
  	client_verify_code(string) # 客户端生成的验证码
  	#（客户端校验完成后服务端再次校验，user_verify_code和client_verify_code一致就行）
  response:
  	code(int)				  # 错误码，无错误为0
  	message(string)   # 错误信息，无错误为success
  	app_info(object)	# 认证信息
  		id(int)					# AppID
  		secret(string)  # AppSecret
  		email(string)   # 邮箱
  		p_num(string)   # 手机号
  ```

### 3、OCR处理接口

- Process V2

  ```yaml
  Path: /toad_ocr/process/v2
  Method: POST
  Headers:
  	Content-Type:multipart/form-data
  	Basic-Token:{{md5 hash str}}
  	# Basic-Token 为用户AppSecret+请求数据大小拼接为字符串使用MD5加密得到的哈希字符串，服务端会用请求长度Content-Length+根据body中的AppID获取AppSecret生成校验串并与Basic-Token进行对比
  UrlParam:
  BodyParam:
  	net_flag(string)	  # 网络标志 有snn/cnn可供选择 
  	file(string)		  	# 待测图片
  	app_id(string)     	# 认证信息 AppID
  response:
  	code(int)					  # 错误码，无错误为0
    label(array/string) # 无错误时为返回预测结果数组，有错误时为具体错误信息
    message(string)		  # 错误信息，无错误为success
  }
  ```

- Process V1（即将废弃）

  ```yaml
  Path: /toad_ocr/process/v2
  Method: POST
  Headers:
  	Content-Type:multipart/form-data
  UrlParam:
  BodyParam:
  	net_flag(string)	  # 网络标志 有snn/cnn可供选择 
  	file(string)		  	# 待测图片
  response:
  	code(int)					  # 错误码，无错误为0
    label(array/string) # 无错误时为返回预测结果数组，有错误时为具体错误信息
    message(string)		  # 错误信息，无错误为success
  }
  ```

  

