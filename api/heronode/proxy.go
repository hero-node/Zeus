package heronode

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"zeus/utils/config"

	"github.com/gin-gonic/gin"
)

func ReverseProxy() gin.HandlerFunc {

	return func(c *gin.Context) {
		target := config.GetIpfsHost()
		targetUrl, err := url.Parse(target)
		if err != nil {
			log.Fatal(err)
		}
		targetHost := targetUrl.Host
		path := c.Request.URL.Path
		version := "v0"

		if strings.Contains(path, "ipfs") {
			path = strings.Replace(path, "/ipfs", "", -1)
			path = "/api/" + version + path
		}

		director := func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = targetHost
			req.URL.Path = path
			fmt.Println(req.URL)
		}

		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
