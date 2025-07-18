package main

import (
	"fmt"
	"os"
	"runtime/pprof"

	"github.com/fatih/color"

	"github.com/golangmc/minecraft-server/impl"
	"github.com/golangmc/minecraft-server/impl/conf"

	"log"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// выполнить функцию через 10 секунд
	// go func() {
	// 	time.Sleep(10 * time.Second)
	// 	f, _ := os.Create("mem1.prof")
	// 	pprof.WriteHeapProfile(f)
	// 	fmt.Println("mem1.prof created")
	// }()

	// в случае паники сохранить профайл
	defer func() {
		if r := recover(); r != nil {
			f, _ := os.Create("panic.prof")
			pprof.WriteHeapProfile(f)
			fmt.Println("panic.prof created")
		}
	}()

	// go func() {
	// 	time.Sleep(46 * time.Second)
	// 	f, _ := os.Create("mem2.prof")
	// 	pprof.WriteHeapProfile(f)
	// 	fmt.Println("mem2.prof created")
	// }()
	color.NoColor = false
	config := readConfig("serverConfig.yml")

	server := impl.NewServer(config)
	server.Load()
}

func readConfig(filename string) (config conf.ServerInfo) {
	// read from yaml file

	return *conf.NewServerInfo(filename)
}
