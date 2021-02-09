package main

import (
	"bili_comment/api"
	"bili_comment/config"
	"fmt"
	"runtime"
)

func main() {
	cfg := config.Cfg()

	api.Run(cfg.Uid, cfg.Page)

	fmt.Println("DONE!!!")

	if runtime.GOOS == "windows" {
		// EXIT
		var tmp string
		fmt.Println("press enter to exit...")
		fmt.Scanln(&tmp)
	}
}
