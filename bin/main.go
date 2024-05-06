package main

import (
	"fmt"
	"github.com/alehua/store-service/internal/pkg/ginx"
	"github.com/alehua/store-service/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
)

func main() {
	var addr int
	pflag.IntVar(&addr, "serverPort", 8000, "for port")
	var filePort int
	pflag.IntVar(&filePort, "filePort", 8001, "help message for port")
	// 解析命令行参数
	pflag.Parse()

	engine := gin.Default()
	engine.LoadHTMLGlob("templates/*")

	handler := web.StoreHandler{}
	handler.RegisterRoutes(engine)
	
	svc := &ginx.Server{
		Engine: engine,
		Addr:   fmt.Sprintf(":%d", addr),
	}

	eg := errgroup.Group{}
	eg.Go(func() error {
		return svc.Start()
	})
	eg.Go(func() error {
		return http.ListenAndServe(fmt.Sprintf(":%d", filePort),
			http.FileServer(http.Dir("./databases")))
	})
	if err := eg.Wait(); err != nil {
		log.Panicln("服务启动失败: ", err.Error())
	}
}
