package main

import (
	"flag"
	l "log"
	"zeus/api/bootstrap"
	"zeus/api/heronode"
	c "zeus/utils/config"
	"zeus/utils/global"

	"github.com/Mercy-Li/Goconfig/config"
	"github.com/gin-gonic/gin"
)

func main() {
	testnet := flag.Bool("testnet", false, "is testnet")
	pconfig := flag.String("config", "./heronode.conf", "config file")
	boot := flag.String("bootlist", "./bootstrap.list", "bootstrap list")
	flag.Parse()

	bootstrap.InitBootStrap(*boot)

	var psection string
	if *testnet {
		psection = "dev"
		global.TEST_NET = true
	} else {
		psection = "pro"
		global.TEST_NET = false
	}

	env := map[string]string{
		"config":  *pconfig,
		"section": psection,
	}

	config.InitConfigEnv(env)
	err := config.LoadConfigFile()
	if err != nil {
		l.Fatalln("Loadconfig File from %s failed. err=%v", *pconfig, err)
	}

	port := c.GetHttpPort()
	gin.SetMode(gin.DebugMode)
	router := gin.Default()

	heronode.InitRoute(router)

	router.Run(port)
}
