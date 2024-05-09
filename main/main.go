package main

import (
	"ACT_GO/bnce"
	"ACT_GO/bngx"
	"ACT_GO/db"
	"ACT_GO/db/entities"
	"ACT_GO/utils"
	cp "github.com/otiai10/copy"
	"github.com/tillberg/autorestart"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
	"runtime"
)

func redirect_log() {
	// Set date, time and filename to the log
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if runtime.GOOS == "linux" {
		log.SetOutput(&lumberjack.Logger{
			Filename:   "/var/log/act_go/act_go.log",
			MaxSize:    500, // megabytes
			MaxBackups: 4,
			MaxAge:     28,   //days
			Compress:   true, // disabled by default
		})
	}
}

func main() {
	//redirect_log()
	go autorestart.RestartOnChange()

	go utils.ListenForUpdatedApp(os.Args[0], "update", func(updated_app_path string) {
		println("Updated app detected  ")
		err := os.Remove(os.Args[0])
		if err != nil {
			log.Fatal(err)
		}
		err = cp.Copy(updated_app_path, os.Args[0])
		if err != nil {
			log.Fatal(err)
		}
		println("Updated app copied")
	})

	println("App started, BuildTime = " + utils.BuildTime)

	//db.Log_Cleanup()
	db.Add_Log(&entities.Log{Message: "App started, BuildTime = " + utils.BuildTime})

	go bnce.Listen_Binance_Klines()
	go bngx.Listen_Account_WS()

	select {}
	//time.Sleep(100 * time.Hour)
}
