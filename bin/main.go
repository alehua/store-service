package main

import (
	"context"
	"github.com/alehua/store-service/internal"
	"github.com/alehua/store-service/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"log"
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

	server := internal.Server{
		Handler:    engine,
		ServerAddr: addr,
		AdminAddr:  filePort,
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := server.Start(ctx); err != nil {
		log.Panicln("服务启动失败: ", err.Error())
	}
}
