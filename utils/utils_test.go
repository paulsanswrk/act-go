package utils

import (
	"github.com/fsnotify/fsnotify"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestWatcher(t *testing.T) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Has(fsnotify.Write) {
					log.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Add a path.
	err = watcher.Add("d:\\Work\\Crypto\\ACT_GO\\utils")
	if err != nil {
		log.Fatal(err)
	}

	// Block main goroutine forever.
	<-make(chan struct{})
}

func TestGuid(t *testing.T) {
	println(Guid())
}

func TestRange(t *testing.T) {
	assert.Equal(t, []float64{0.5, 1, 1.5, 2, 2.5, 3, 3.5, 4, 4.5}, Range(0.5, 4.5, 0.5))
	assert.Equal(t, []float64{4.5, 4, 3.5, 3, 2.5, 2, 1.5, 1, 0.5}, Range(4.5, 0.5, 0.5))
	assert.Equal(t, []float64{-0.5, -1, -1.5, -2, -2.5, -3, -3.5, -4, -4.5}, Range(-0.5, -4.5, 0.5))
	assert.Equal(t, []float64{-0.5, -1, -1.5, -2, -2.5, -3, -3.5, -4, -4.5}, Range(-0.5, -4.5, -0.5))

	//fmt.Printf("%#v", Range(4.5, 0.5, 0.5))
}
