package log

import (
	"io/ioutil"
	stlog "log"
	"net/http"
	"os"
)

var log *stlog.Logger

type fileLog string

// 写日志到文件
func (fl fileLog) Write(data []byte) (int, error) {
	f, err := os.OpenFile(string(fl), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return f.Write(data)
}

func Run(destination string) {
	log = stlog.New(fileLog(destination), "[log] - ", stlog.LstdFlags)
}

func RegisterHandlers() {
	http.HandleFunc("/log", func(rw http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			msg, err := ioutil.ReadAll(r.Body)
			if err != nil || len(msg) == 0 {
				rw.WriteHeader(http.StatusBadRequest)
				return
			}
			write(string(msg))
		}
	})
}

func write(message string) {
	log.Printf("%v\n", message)
}
