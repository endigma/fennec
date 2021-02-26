package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

var config Config

// Config contains the config for the application
type Config struct {
	Port     string    `json:"port"`
	Handlers []Handler `json:"handlers"`
}

// Handler contains info about a handler
type Handler struct {
	Path string `json:"path"`
	Run  string `json:"run"`
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func unpackConfig() Config {
	path, err := os.Getwd()
	checkErr(err)

	configFile, err := os.Open(path + "/config.json")
	checkErr(err)

	defer configFile.Close()

	byteValue, err := ioutil.ReadAll(configFile)
	checkErr(err)

	var config Config

	err = json.Unmarshal(byteValue, &config)
	checkErr(err)

	log.Info("Successfully Unpacked Config")

	return config
}

func getHandler(query string) (int, Handler) {
	for _, h := range config.Handlers {
		if h.Path == query {
			return 0, h
		}
	}
	return 1, Handler{}
}

func catch(rw http.ResponseWriter, req *http.Request) {
	stat, h := getHandler(req.URL.Path)
	if stat != 0 {
		log.WithFields(log.Fields{
			"ip": ip(req),
		}).Warnf("Request Dropped: %s", req.URL.Path)
		fmt.Fprintf(rw, "Fuck off %s!\n", ip(req))
		return
	}

	fmt.Fprintf(rw, "Request is being handled.\n")
	log.Infof("Recieved Request: %s", req.URL.Path)

	exec.Command(h.Run).Run()

	log.Infof("Ran Handler: %s", h.Run)
}

func ip(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

func init() {
	logfmt := new(log.TextFormatter)
	logfmt.TimestampFormat = "2006-01-02 15:04:05"
	logfmt.FullTimestamp = true
	log.SetFormatter(logfmt)

	f, err := os.OpenFile("watcher.log", os.O_WRONLY|os.O_CREATE, 0755)
	checkErr(err)

	mw := io.MultiWriter(os.Stdout, f)
	log.SetOutput(mw)

	config = unpackConfig()
}

func main() {
	http.HandleFunc("/", catch)
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}
