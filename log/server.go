package log

import (
	"io/ioutil"
	stlog "log"
	"net/http"
	"os"
)

var log *stlog.Logger
var logfile string

type filelog string
type mylog int32

func (ml mylog) Write(data []byte) (int, error) {
	f, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return f.Write(data)

}

func (fl filelog) Write(data []byte) (int, error) {
	f, err := os.OpenFile(string(fl), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return f.Write(data)

}

// func (fl filelog) Write(data []byte) (int, error) {
// 	f, err := os.OpenFile(string(fl), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
// 	if err != nil {
// 		return 0, err
// 	}
// 	defer f.Close()
// 	return f.Write(data)

// }
func Run(destination string) {
	logfile = destination
	log = stlog.New(mylog(1), "go: ", stlog.LstdFlags)

}

// func Run(destination string) {
// log = stlog.New(filelog(destination), "go: ", stlog.LstdFlags)

// }

func RegisterHandlers() {
	http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			msg, err := ioutil.ReadAll(r.Body)
			if err != nil || len(msg) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			write(string(msg))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			// return
		}
	})
}

func write(messages string) {
	log.Printf("%v\n", messages)
}
