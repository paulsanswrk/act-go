package utils

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"
)

func ListenForUpdatedApp(app_path string, updates_subfolder string, callback func(updated_app_path string)) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	app_path, err = os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	app_dir := filepath.Dir(app_path)
	app_fn := filepath.Base(app_path)
	updated_app_dir := path.Join(app_dir, updates_subfolder)
	updated_app_path := path.Join(updated_app_dir, app_fn)
	println("updated_app_dir: ", updated_app_dir)

	timeout := 1 * time.Second
	ticker := time.NewTicker(timeout)
	defer ticker.Stop()

	modify_event_received := false
	modify_time := time.Now()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				//log.Println("event:", event)
				if event.Has(fsnotify.Write) {
					//log.Println("modified file: ", event.Name)
					if event.Name == updated_app_path {
						modify_event_received = true
						modify_time = time.Now()
					}
				}
			case err, ok := <-watcher.Errors:
				log.Println("error:", err)
				if !ok {
					return
				}
			case <-ticker.C:
				if modify_event_received && time.Now().After(modify_time.Add(timeout)) {
					modify_event_received = false //wait for next update
					callback(updated_app_path)
				}
			}
		}
	}()

	// Add a path.
	err = watcher.Add(updated_app_dir)
	if err != nil {
		log.Fatal(err)
	}

	<-make(chan struct{})

}
