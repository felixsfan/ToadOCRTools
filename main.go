package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"suvvm.work/ToadOCRTools/config"
	"suvvm.work/ToadOCRTools/dal/db"
	"suvvm.work/ToadOCRTools/method"
)

var (
	dbConfig = "./conf/db_config.yaml"
	sdkConfig = "./conf/al_sdk_config.yaml"
)

// InitConfig 初始化配置信息
func InitConfig() {
	str, err := os.Getwd() // 获取相对路径
	if err != nil {
		panic(fmt.Sprintf("filepath failed, err=%v", err))
	}
	dbFileName, err := filepath.Abs(filepath.Join(str, dbConfig)) // 获取db配置文件路径
	if err != nil {
		panic(fmt.Sprintf("filepath failed, err=%v", err))
	}
	conf := config.Init(dbFileName)                    // 读取db配置文件
	if err = db.InitDB(&conf.DBConfig); err != nil { // 初始化db链接
		panic(fmt.Sprintf("init db conn err=%v", err))
	}
	sdkFileName, err := filepath.Abs(filepath.Join(str, sdkConfig)) // 获取db配置文件路径
	if err != nil {
		panic(fmt.Sprintf("filepath failed, err=%v", err))
	}
	conf = config.Init(sdkFileName)
}

func main() {
	InitConfig()
	//log.Printf("SDK:%v", config.AppConfig.SdkConfig)
	method.InitEmail()
	r := gin.Default()
	r.Use(Cors())
	register(r)
	log.Printf("run toad ocr api service...")
	if err := r.Run(":18889"); err != nil {
		log.Printf("run api service fail!")
	}
}

// Cors 跨域
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			//接收客户端发送的origin （重要！）
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			//服务器支持的所有跨域请求的方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			//允许跨域设置可以返回其他子段，可以自定义字段
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session")
			// 允许浏览器（客户端）可以解析的头部 （重要）
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			//设置缓存时间
			c.Header("Access-Control-Max-Age", "172800")
			//允许客户端传递校验信息比如 cookie (重要)
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		//允许类型校验
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic info is: %v", err)
			}
		}()
		c.Next()
	}
}
