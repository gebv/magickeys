package main

import (
	"api"
	"flag"
	// "github.com/golang/glog"
	_ "models"
	_ "store"
	"utils"

	"os"
	"os/signal"
	"web"
	"syscall"
)

var flagConfigFile string 

func main() {
	flag.StringVar(&flagConfigFile, "config", "config.json", "")

	flag.Parse()

	utils.LoadConfig(flagConfigFile)

	api.NewServer()
	api.InitApi()
	web.InitWeb()
	api.StartServer()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-c

	api.StopServer()
}