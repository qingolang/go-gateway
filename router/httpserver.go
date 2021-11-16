package router

import (
	"context"
	"go_gateway/common/lib"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	HTTPSrvHandler *http.Server
)

// HTTPServerRun
func HTTPServerRun() {
	gin.SetMode(lib.GetStringConf("base.base.debug_mode"))
	r := InitRouter()

	// Init HTTP server
	HTTPSrvHandler = &http.Server{
		Addr:           lib.GetStringConf("base.http.addr"),
		Handler:        r,
		ReadTimeout:    time.Duration(lib.GetIntConf("base.http.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(lib.GetIntConf("base.http.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << uint(lib.GetIntConf("base.http.max_header_bytes")),
	}
	go func() {
		log.Printf(" [INFO] HTTPServerRun:%s\n", lib.GetStringConf("base.http.addr"))
		if err := HTTPSrvHandler.ListenAndServe(); err != nil {
			log.Fatalf(" [ERROR] HTTPServerRun:%s err:%v\n", lib.GetStringConf("base.http.addr"), err)
		}
	}()
}

// HTTPServerStop
func HTTPServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := HTTPSrvHandler.Shutdown(ctx); err != nil {
		log.Fatalf(" [ERROR] HTTPServerStop err:%v\n", err)
	}
	log.Printf(" [INFO] HTTPServerStop stopped\n")
}
