package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jander/golog/logger"
)

var host string

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", time.Now())
}

func main() {
	flag.StringVar(&host, "host", "localhost:8080", "主机名:端口")
	flag.Parse()
	logger.SetHandlers(logger.NewFileHandler("hello.txt"), logger.NewConsoleHandler())
	logger.Println(os.Args)
	http.HandleFunc("/", greet)
	http.ListenAndServe(host, nil)
	logger.Println("exit!")
}
