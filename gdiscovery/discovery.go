package gdiscovery

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"log"
	"net"
	"strings"
)

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

func RegisterInstance(config IConfig) {
	configServerHost := config.String("go.config.host")
	configServerPort := config.Int64("go.config.port")
	userName := config.String("go.config.username")
	password := config.String("go.config.password")

	appGroup := config.String("go.application.group")
	appName := config.String("go.application.name")
	appPort := config.Int64("go.application.port")

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

	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic("初, err: " + err.Error())
	}

	instanceIP := config.String("go.instance.ip")
	if instanceIP == "" {
		ipList, err := localIPv4s(config.Bool("go.application.lan"), config.String("go.application.lannet"))
		if err != nil {
			panic("获取本地IP失败")
		}
		instanceIP = ipList[0]
	}

	if !config.Bool("go.instance.debug") {
		_, err := namingClient.RegisterInstance(vo.RegisterInstanceParam{
			Ip:          instanceIP,
			Port:        uint64(appPort),
			ServiceName: appName,
			Weight:      10,
			Enable:      true,
			Healthy:     true,
			Ephemeral:   true,
			Metadata:    map[string]string{"debug": config.String("go.instance.debug")},
			ClusterName: "",       // default value is DEFAULT
			GroupName:   appGroup, // default value is DEFAULT_GROUP
		})

		if err != nil {
			panic("注册服务实例失败 err: " + err.Error())
		}

		log.Println("服务实例注册成功")
	}
}
