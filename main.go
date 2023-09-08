package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gme/gconfig"
	"gme/gdiscovery"
	"gme/gmysql"
	"gme/router"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func parseArgs() string {
	var configFile string
	flag.StringVar(&configFile, "f", os.Args[0]+".yml", "yml配置文件名")
	flag.Parse()
	path, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	if !strings.Contains(configFile, "/") {
		configFile = path + "/" + configFile
	}
	return configFile
}

func localIPv4s(lan bool, lanNetwork string) ([]string, error) {
	var ips, ipLans, ipWans []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ips, err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && ipnet.IP.IsGlobalUnicast() && ipnet.IP.To4() != nil {
			if ipnet.IP.IsPrivate() {
				ipLans = append(ipLans, ipnet.IP.String())
				if lan && strings.HasPrefix(ipnet.IP.String(), lanNetwork) {
					ips = append(ips, ipnet.IP.String())
				}
			}
			if !ipnet.IP.IsPrivate() {
				ipWans = append(ipWans, ipnet.IP.String())
				if !lan {
					ips = append(ips, ipnet.IP.String())
				}
			}
		}
	}
	if len(ips) == 0 {
		if lan {
			ips = append(ips, ipWans...)
		} else {
			ips = append(ips, ipLans...)
		}
	}
	return ips, nil
}

func main() {
	configFile := parseArgs()

	gconfig.Mgr.LoadConfig(configFile)

	gdiscovery.RegisterInstance(&gconfig.Mgr)

	gmysql.Init(&gconfig.Mgr)

	gin.DefaultWriter = io.Discard
	engine := gin.Default()
	engine.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.SetupRouter(engine)

	//instance, err := namingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
	//	ServiceName: appName,
	//	GroupName:   appGroup,            // 默认值DEFAULT_GROUP
	//	Clusters:    []string{"DEFAULT"}, // 默认值DEFAULT
	//})

	//instanceIp, instancePort := gdiscovery.GetInstanceInfo()
	//
	//fmt.Println("ip: ", instance.Ip, " port: ", instance.Port)
	//
	//router := gin.Default()
	//
	//listenAddr := fmt.Sprintf("%s:%d", instanceIP, appPort)
	//server := &http.Server{
	//	Addr:    listenAddr,
	//	Handler: router,
	//}
	//
	//// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	//signalChan := make(chan os.Signal)
	//signal.Notify(signalChan, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)
	//sig := <-signalChan
	//log.Println("Get Signal:" + sig.String())
	//log.Println("Shutdown Server")
	//
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	//if err := server.Shutdown(ctx); err != nil {
	//	log.Println("Server Shutdown err:" + err.Error())
	//}
	//log.Println("Server exit")

}
