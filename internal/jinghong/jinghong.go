package jinghong

import (
	"context"
	"errors"
	"fmt"
	"jonghong/internal/pkg/known"
	"jonghong/internal/pkg/log"
	"jonghong/internal/pkg/middleware"
	"jonghong/pkg/token"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// 配置命令行工具

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		// 定义根命令为jinghong时使用
		Use: "jinghong",
		// 应用描述信息
		Short: "An api project for jinghong app",
		Long:  "An api project for jinghong app, based on golong",
		// 暂时不知道有啥用
		SilenceUsage: true,
		// 运行时
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
		// 禁止没带标签的参数
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}
			return nil
		},
	}
	// 增加flag参数
	// PersistentFlags运用于当前命令及其子命令
	// Flags 只运用于当前命令
	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "The path for config file")

	// 读取配置文件
	cobra.OnInitialize(initDirPath, initConfig, func() {
		log.Init(logOptions())
	})
	return cmd
}

func run() error {

	// 初始化数据库连接
	if err := initStore(); err != nil {
		return err
	}

	if err := initService(); err != nil {
		return err
	}

	// 初始化token数据

	if err := token.Init(viper.GetString("jwt-scret"), known.XUsernameKey); err != nil {
		return err
	}

	// 设置gin的运行模式
	gin.SetMode(viper.GetString("runmode"))
	// 新建路由引擎
	g := gin.New()
	// 头文件中间件
	mws := []gin.HandlerFunc{
		gin.Recovery(),
		middleware.NoCache,
		middleware.Cors,
		middleware.Secure,
	}
	g.Use(mws...)

	if err := initRouter(g); err != nil {
		return err
	}

	// 开启http监听，https的拦截放在nginx服务端

	httpsrv := &http.Server{
		Addr:    viper.GetString("addr"),
		Handler: g,
	}

	// 在协程中启动http服务, 主进程做其他事情
	go func() {
		if err := httpsrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalw(err.Error())
		}
	}()

	// 接受系统停止消息，并停止该进程

	// 建立一个大小为1的通道接收系统消息
	quit := make(chan os.Signal, 1)
	// 监听系统关闭和杀停信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// 阻塞进程 等待系统信号
	<-quit
	log.Infow("Shutting down server ... ")

	// 10秒的容错时间
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpsrv.Shutdown(ctx); err != nil {
		log.Errorw("Insecure Server forced to shutdown", "err", err)
		return err
	}

	log.Infow("Server exiting")

	return nil
}
