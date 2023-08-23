package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/v2"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
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
	k := koanf.New(".")
	if err := k.Load(file.Provider(configFile), yaml.Parser()); err != nil {
		panic("fail to load config file: " + configFile)
	}
	configServerHost := k.String("go.config.host")
	configServerPort := k.Int64("go.config.port")
	userName := k.String("go.config.username")
	password := k.String("go.config.password")

	appGroup := k.String("go.application.group")
	appName := k.String("go.application.name")
	appPort := k.Int64("go.application.port")
	env := k.String("go.config.env")
	dataId := appName + "-" + env

	sc := []constant.ServerConfig{
		*constant.NewServerConfig(configServerHost, uint64(configServerPort), constant.WithContextPath("/nacos")),
	}

	//sc := []constant.ServerConfig{
	//	*constant.NewServerConfig("59.56.77.17", 58848, constant.WithContextPath("/nacos")),
	//}

	//create ClientConfig
	cc := constant.ClientConfig{
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		Username:            userName,
		Password:            password,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
	}

	// create config client
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

	if err != nil {
		log.Fatalln("连接配置中心失败 err: {}", err.Error())
		return
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  appGroup,
	})

	if err != nil {
		log.Fatalln("读取配置中心配置失败, err: ", err.Error())
		return
	}

	if err := k.Load(rawbytes.Provider([]byte(content)), yaml.Parser()); err != nil {
		panic(err.Error())
	}

	err = configClient.ListenConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  appGroup,
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("config changed group:" + group + ", dataId:" + dataId + ", content:" + data)
		},
	})

	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic("初, err: " + err.Error())
	}

	instanceIP := k.String("go.instance.ip")
	if instanceIP == "" {
		ipList, err := localIPv4s(k.Bool("go.application.lan"), k.String("go.application.lannet"))
		if err != nil {
			panic("获取本地IP失败")
		}
		instanceIP = ipList[0]
	}

	if !k.Bool("go.instance.debug") {
		_, err := namingClient.RegisterInstance(vo.RegisterInstanceParam{
			Ip:          instanceIP,
			Port:        uint64(appPort),
			ServiceName: appName,
			Weight:      10,
			Enable:      true,
			Healthy:     true,
			Ephemeral:   true,
			Metadata:    map[string]string{"debug": k.String("go.instance.debug")},
			ClusterName: "",       // default value is DEFAULT
			GroupName:   appGroup, // default value is DEFAULT_GROUP
		})

		if err != nil {
			panic("注册服务实例失败 err: " + err.Error())
		}

		log.Println("服务实例注册成功")
	}

	instance, err := namingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: appName,
		GroupName:   appGroup,            // 默认值DEFAULT_GROUP
		Clusters:    []string{"DEFAULT"}, // 默认值DEFAULT
	})

	fmt.Println("ip: ", instance.Ip, " port: ", instance.Port)

	router := gin.Default()

	listenAddr := fmt.Sprintf("%s:%d", instanceIP, appPort)
	server := &http.Server{
		Addr:    listenAddr,
		Handler: router,
	}

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)
	sig := <-signalChan
	log.Println("Get Signal:" + sig.String())
	log.Println("Shutdown Server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown err:" + err.Error())
	}
	log.Println("Server exit")

}
