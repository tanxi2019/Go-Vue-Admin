package initialize

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"server/app/service/system"
	"server/config"
	"server/global"
	"server/middleware"
	"server/router"
	"syscall"
	"time"
)

// InitializeRun
func InitRun() {
	// 操作日志中间件处理日志时没有将日志发送到rabbitmq或者kafka中, 而是发送到了channel中
	// 这里开启3个goroutine处理channel将日志记录到数据库
	logDao := system.NewLogService()
	for i := 0; i < 3; i++ {
		go logDao.SaveOperationLogChannel(middleware.OperationLogChan)
	}

	// 注册所有路由
	r := router.InitRouter()

	host := config.Conf.System.Host
	port := config.Conf.System.Port

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: r,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			global.Log.Fatalf("listen: %s\n", err)
		}
	}()

	global.Log.Info("服务启动完成")

	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	global.Log.Info("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		global.Log.Fatal("Server forced to shutdown:", err)
	}

	global.Log.Info("Server exiting!")
}
