package main

import (
	"fmt"
	"zeus/api/bootstrap"
	"runtime"
	"zeus/api/core"
	c "zeus/utils/config"
	"zeus/utils/global"

	"github.com/Mercy-Li/Goconfig/config"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)


var testnet bool
var configFile string
var bootFile string

func main() {


	rootCmd := &cobra.Command{
		Use:   "gher",
		Short: "Heronode official node cmd",
		Long:  "Heronode official node cmd, providing api wrapped eth, ipfs, etc.",
	}

	runCmd := &cobra.Command{
		Use: "run",
		Aliases:[]string{"r"},
		Short:"run heronode api server",
		Long:"run heronode api server",
		Run: run,
	}
	runCmd.PersistentFlags().BoolVarP(&testnet, "testnet", "t", false, "if testnet")
	runCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "./heronode.conf", "config file")
	runCmd.PersistentFlags().StringVarP(&bootFile, "boot", "b", "./bootstrap.list", "bootlist file")
	rootCmd.AddCommand(runCmd)

	versionCmd := &cobra.Command{
		Use:"version",
		Aliases:[]string{"v"},
		Short:"version of gher",
		Long:"version of gher",
		Run: version,
	}
	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}

}

func run(cmd *cobra.Command, args []string)  {
	bootstrap.InitBootStrap(bootFile)
	var section string
	if testnet {
		section = "dev"
		global.TEST_NET = true
	} else {
		section = "pro"
		global.TEST_NET = false
	}

	env := map[string]string{
		"config":  configFile,
		"section": section,
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
	core.InitRoute(router)

	router.Run(port)
}

func version(cmd *cobra.Command, args []string)  {
	v := global.VERSION + runtime.GOOS + "-" + runtime.GOARCH + "/" + runtime.Version()
	fmt.Println(v)
}