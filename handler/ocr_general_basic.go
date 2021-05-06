package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"suvvm.work/ToadOCRTools/common"
	"suvvm.work/ToadOCRTools/rpc"
)

var supportImageExtNames = []string{".jpg", ".jpeg", ".png", ".ico", ".bmp", ".gif"}


// Pong ping-pong测试接口
//
// 入参
//	ctx *gin.Context	// 上下文参数
func Pong(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": common.HandlerSuccess,
		"msg": "pong",
	})
}

func Process(ctx *gin.Context) {
	file, fileHeader, err := ctx.Request.FormFile("file")	// 读取上传文件
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":     -1,
			"message": "文件上传失败，http 400",
			"labels":  nil,
		})
		return
	}
	if fileHeader.Size == 0 {	// 判断文件大小是否合法
		log.Printf("recv screenshot data size is zero")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":     -1,
			"message": "文件上传失败，接受到的文件大小为0",
			"labels":  nil,
		})
		return
	} else if fileHeader.Size > common.ImageSizeLimit {
		log.Printf("recv screenshot data size is over")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":     -1,
			"message": "文件上传失败，上传文件大小超出限制",
			"labels":  nil,
		})
		return
	}
	filename := fileHeader.Filename
	if strings.Contains(filename, "." + string(os.PathSeparator)){	// 判断文件名是否合法
		log.Printf("recv illegal file ")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":     -1,
			"message": "文件命名非法，含有路径分隔符",
			"labels":  nil,
		})
		return
	}
	extname := strings.ToLower(path.Ext(filename))	// 判断扩展名是否合法
	if !isImage(extname) {
		log.Printf("recv file is not a picture")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":     -1,
			"message": "文件扩展名非法",
			"labels":  nil,
		})
		return
	}
	data := make([]byte, fileHeader.Size)
	if _, err := file.Read(data); err != nil {	// 读取文件内容至字节数组中
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":     -1,
			"message": "文件内容读取失败",
			"labels":  nil,
		})
		return
	}
	netFlag := ctx.PostForm("net_flag")
	labels, err := rpc.Process(netFlag, data)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":     -1,
			"message": "rpc 调用失败",
			"labels":  err,
		})
		return
	}
	log.Printf("success")
	ctx.JSON(200, gin.H{
		"code":     0,
		"message":  "success",
		"label":   labels,
	})
}

func isImage(extName string) bool {
	for i := 0; i < len(supportImageExtNames); i++ {
		if supportImageExtNames[i] == extName {
			return true
		}
	}
	return false
}
