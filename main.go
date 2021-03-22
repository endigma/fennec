package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

var config Config
var basePath string
var debug = flag.Bool("d", false, "debug information")

// Config contains the config for the application
type Config struct {
	Port     string    `json:"port"`
	Handlers []Handler `json:"handlers"`
}

// Req contains a request item
type Req struct {
	IP   string
	Path string

	Secret string      `json:"secret"`
	Data   interface{} `json:"data"`
}

// Handler contains info about a handler
type Handler struct {
	Path        string `json:"path"`
	Command     string `json:"command"`
	ForwardData bool   `json:"forwarddata"`
	Secret      string `json:"secret"`
}

func (h Handler) run(data interface{}) {
	var err error

	d, err := json.Marshal(data)

	if h.ForwardData {
		err = exec.Command(h.Command, string(d)).Run()
	} else {
		err = exec.Command(h.Command).Run()
	}

	if err != nil {
		log.WithFields(log.Fields{
			"cmd":  h.Command,
			"args": string(d),
		}).Warn("Error Running Handler")
		log.Warn(err)
	} else {
		log.WithFields(log.Fields{
			"cmd":  h.Command,
			"args": string(d),
		}).Info("Ran Handler")
	}

}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func unpackConfig() Config {
	configFile, err := os.Open(os.Args[1])
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

func unpackReq(req *http.Request) (Req, error) {
	var request Req
	request.Path = req.URL.Path
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		return Req{}, err
	}

	return request, nil
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
	if req.Method == "POST" {
		request, err := unpackReq(req)
		if err != nil {
			logReq(req, "Dropped", "Invalid JSON")
			rw.WriteHeader(http.StatusBadRequest)

			return
		}

		stat, handler := getHandler(request.Path)
		if stat != 0 {
			logReq(req, "Dropped", "Invalid Path")
			rw.WriteHeader(http.StatusNotFound)

			return
		}

		log.Info(request.Data)

		if request.Secret != handler.Secret {
			logReq(req, "Refused", "Incorrect Secret")
			rw.WriteHeader(http.StatusUnauthorized)

			return

		}

		logReq(req, "Accepted", "OK")
		rw.WriteHeader(http.StatusOK)

		go handler.run(request.Data)
	} else {
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func logReq(req *http.Request, status string, msg string) {
	log.WithFields(log.Fields{
		"path": req.URL.Path,
		"ip":   ip(req),
	}).Infof("Request %s: %s", status, msg)
}

func ip(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

func init() {
	flag.Parse()
	if len(os.Args) == 1 {
		log.Fatal("Please provide a valid config file.")
	}

	if _, err := os.Stat(os.Args[1]); os.IsNotExist(err) {
		log.Fatal("Please provide a valid config file.")
	}

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	ex, err := os.Executable()
	checkErr(err)

	basePath = filepath.Dir(ex)

	logfmt := new(log.TextFormatter)
	logfmt.TimestampFormat = "2006-01-02 15:04:05"
	logfmt.FullTimestamp = true
	log.SetFormatter(logfmt)
	log.SetOutput(os.Stdout)

	config = unpackConfig()
}

func main() {
	http.HandleFunc("/", catch)
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}
