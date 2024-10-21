package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var logFile *os.File

func init() {
	var err error
	logFile, err = os.OpenFile("channel_log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logFile)
}

func logOperation(operation string) {
	log.Println(operation)
}

func main() {
	ch := make(chan int, 1)

	go func() {
		for i := 0; i < 5; i++ {
			logOperation(fmt.Sprintf("Sending %d to channel", i))
			ch <- i
			time.Sleep(1 * time.Second)
		}
		close(ch)
	}()

	go func() {
		for val := range ch {
			logOperation(fmt.Sprintf("Received %d from channel", val))
			time.Sleep(1 * time.Second)
		}
	}()

	http.HandleFunc("/logs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		http.ServeFile(w, r, "vch/runtime/log/vch.log")
	})

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
