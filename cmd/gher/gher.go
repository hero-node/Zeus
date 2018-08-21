package main

import (
	"flag"
	l "log"
	"zeus/api/bootstrap"
	"zeus/api/heronode"

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
	} else {
		psection = "pro"
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

	port, err := config.GetConfigString("api.listen")
	if err != nil {
		l.Fatalln("get port config error")
	}

	gin.SetMode(gin.DebugMode)
	router := gin.Default()

	heronode.InitRoute(router)

	router.Run(port)
}
