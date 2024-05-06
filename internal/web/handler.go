package web

import (
	"github.com/alehua/store-service/internal/pkg/ginx"
	"github.com/alehua/store-service/internal/pkg/logger"
	"github.com/alehua/store-service/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

type StoreHandler struct {
}

var log = logger.NewZapLogger()

func (s *StoreHandler) RegisterRoutes(engine *gin.Engine) {
	g := engine.Group("/file")

	g.GET("/test", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})

	g.POST("/upload", s.Upload)

	g.GET("/download/:file", s.Download)

	g.GET("/index", s.Index)
}

func (s *StoreHandler) Upload(ctx *gin.Context) {
	rFile, err := ctx.FormFile("file")
	if err != nil {
		log.Error("读取file字段失败", logger.Error(err))
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 1,
			Msg:  "读取file字段失败",
			Data: err.Error(),
		})
	} else {
		// 开始文件上传
		err = service.Save(rFile)
		if err != nil {
			ctx.JSON(http.StatusOK, ginx.Result{
				Code: 2,
				Msg:  "文件保存数据库失败",
				Data: err.Error(),
			})
		} else {
			ctx.String(http.StatusOK, "保存成功")
		}
	}

}

func (s *StoreHandler) Download(ctx *gin.Context) {
	fileName := ctx.Param("file")
	targetPath := service.DownLoad(fileName)
	_, err := os.Open(targetPath)
	if err != nil {
		ctx.JSON(http.StatusOK, &ginx.Result{
			Code: 3,
			Msg:  "文件读取失败",
			Data: err.Error(),
		})
		return
	}
	ctx.Header("Content-Type", "application/octet-stream")
	//强制浏览器下载
	ctx.Header("Content-Disposition", "attachment; filename="+fileName)
	ctx.Header("Content-Transfer-Encoding", "binary")

	ctx.File(targetPath)
}

func (s *StoreHandler) Index(ctx *gin.Context) {

}
