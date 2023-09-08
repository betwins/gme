package gconfig

import (
	"fmt"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/v2"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"log"
)

type configMgr struct {
	ko *koanf.Koanf
}

var Mgr configMgr

func (c *configMgr) LoadConfig(configFile string) {
	c.ko = koanf.New(".")
	if err := c.ko.Load(file.Provider(configFile), yaml.Parser()); err != nil {
		log.Fatalln("fail to load config file: " + configFile)
	}
	configServerHost := c.ko.String("go.config.host")
	configServerPort := c.ko.Int64("go.config.port")
	userName := c.ko.String("go.config.username")
	password := c.ko.String("go.config.password")

	appGroup := c.ko.String("go.application.group")
	appName := c.ko.String("go.application.name")
	/*	appPort := c.ko.Int64("go.application.port")*/
	env := c.ko.String("go.config.env")
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
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  appGroup,
	})

	if err != nil {
		log.Fatalln("读取配置中心配置失败, err: ", err.Error())
	}

	if err := c.ko.Load(rawbytes.Provider([]byte(content)), yaml.Parser()); err != nil {
		panic(err.Error())
	}

	err = configClient.ListenConfig(vo.ConfigParam{
		DataId:   dataId,
		Group:    appGroup,
		OnChange: onConfigChange,
	})

}

func onConfigChange(namespace, group, dataId, data string) {
	fmt.Println("config changed group:" + group + ", dataId:" + dataId + ", content:" + data)
	newKo := koanf.New(".")
	if err := newKo.Load(rawbytes.Provider([]byte(data)), yaml.Parser()); err != nil {
		log.Fatalln("fail to load config data: " + data)
	}
	//oldKo := Mgr.ko
	Mgr.ko = newKo

}
