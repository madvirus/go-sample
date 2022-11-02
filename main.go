package main

import (
	"flag"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"go-sample/appctx"
	"go-sample/server"
)

func main() {
	var configFlag = flag.String("config", "", "config file path")
	flag.Parse()

	err := loadEnv(*configFlag)
	if err != nil {
		log.Errorf("main: fail to load env file %s: %v", *configFlag, err)
		return
	}
	config := appctx.CreateConfigFromEnv()
	appCtx, err := appctx.Initialize(config)
	if err != nil {
		log.Fatalf("main: fail to Initialize: %v", err)
		return
	}
	log.Infof("main: initialized AppCtx")

	apiServer, err := server.CreateApiServer(config, appCtx)
	if err != nil {
		log.Fatalf("main: fail to CreateApiServer: %v", err)
		return
	}
	apiServer.Run()
}

func loadEnv(envFile string) error {
	if envFile != "" {
		log.Infof("use env file path: %s", envFile)
		return godotenv.Load(envFile)
	} else {
		log.Infof("use default env file")
		return godotenv.Load()
	}
}
