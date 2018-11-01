package main

import (
	"flag"
	"fmt"
	"zeus/api/bootstrap"
	"zeus/api/core"
	c "zeus/utils/config"
	"zeus/utils/global"

	"github.com/Mercy-Li/Goconfig/config"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
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

	err := config.InitConfigEnv(env)
	err = config.LoadConfigFile()
	if err != nil {
		fmt.Println("No config file. Use default parameters.")
	}

	go c.GetValidEthHost()

	port := c.GetHttpPort()
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))
	core.InitRoute(router)

	router.Run(port)
}
