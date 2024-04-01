package main

import (
	"fmt"
	"goTSVParser/config"
	"goTSVParser/internal/handler"
	"goTSVParser/internal/service"
	"goTSVParser/internal/storage"
	"goTSVParser/internal/watcher"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cnfg := config.New()
	st, err := storage.NewDBStorage(cnfg)
	if err != nil {
		return
	}
	w := watcher.NewWatcher(cnfg)
	s := service.NewService(st, w, cnfg)
	h := handler.NewHandler(s, cnfg)
	go func() {
		err := s.Scanner()
		fmt.Println(err)
		if err != nil {
			log.Fatal(err)
		}
	}()
	go h.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	log.Println("stopping application")
	st.ShutDown()
	log.Println("shutting down application")
}
