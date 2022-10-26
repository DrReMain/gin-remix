package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"go-remix/config"
)

func main() {
	log.Println("Starting server...")

	if gin.Mode() != gin.ReleaseMode {
		err := godotenv.Load()
		if err != nil {
			log.Fatalln("加载 .env 文件错误")
		}
	}

	ctx := context.Background()

	cfg, err := config.LoadConfig(ctx)
	if err != nil {
		log.Fatalf("加载配置错误: %v\n", err)
	}

	ds, err := initDS(ctx, cfg)
	if err != nil {
		log.Fatalf("初始化数据源错误: %v\n", err)
	}

	router, err := inject(ds, cfg)
	if err != nil {
		log.Fatalf("路由注入数据源错误: %v\n", err)
	}

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}
	go func() {
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("初始化服务器错误: %v\n", err)
		}
	}()

	log.Printf("Listening on %v\n", srv.Addr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := ds.close(); err != nil {
		log.Fatalf("A problem occurred gracefully shutting down data sources: %v\n", err)
	}

	log.Println("关闭服务器...")
	if err = srv.Shutdown(ctx); err != nil {
		log.Fatalf("服务器强制关闭错误: %v\n", err)
	}
}
