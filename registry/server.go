package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

const (
	ServerPort  = ":3000"
	ServicesURL = "http://localhost" + ServerPort + "/services"
)

type registry struct {
	registrations []Registration
	mutex         *sync.RWMutex
}

func (r *registry) add(reg Registration) error {
	r.mutex.Lock()
	r.registrations = append(r.registrations, reg)
	r.mutex.Unlock()

	err := r.sendRequiredServices(reg)
	if err != nil {
		return err
	}

	return nil

}

func (r *registry) sendRequiredServices(reg Registration) error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var p patch
	for _, serviceReg := range r.registrations {
		for _, reqService := range reg.RequiredServices {
			if serviceReg.ServiceName == reqService {
				p.Added = append(p.Added, patchEntry{
					Name: reqService,
					URL:  serviceReg.ServiceURL,
				})
			}
		}
	}
	return r.sendPatch(p, reg.ServiceUpdateURl)
}

func (r *registry) sendPatch(p patch, url string) error {
	d, err := json.Marshal(p)
	if err != nil {
		return err
	}
	_, err = http.Post(url, "application/json", bytes.NewBuffer(d))

	if err != nil {
		return err
	}
	return nil

}

func (r *registry) remove(url string) error {
	for i := range reg.registrations {
		if reg.registrations[i].ServiceURL == url {

			// 这个锁 对吗   ?....
			reg.mutex.Lock()

			reg.registrations = append(reg.registrations[:i], reg.registrations[i+1:]...)
			reg.mutex.Unlock()
			return nil
		}
	}
	return fmt.Errorf("Service at URL " + url + "found")
}

var reg = registry{
	registrations: make([]Registration, 0),
	mutex:         new(sync.RWMutex),
}

type RegistryService struct{}

func (s RegistryService) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		dec := json.NewDecoder(r.Body)
		var r Registration
		err := dec.Decode(&r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("Adding service: %v with URL: %s\n", r.ServiceName, r.ServiceURL)
		err = reg.add(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	case http.MethodDelete:
		payload, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Panicln(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		url := string(payload)
		log.Printf("Removeing service at URL:%s", url)
		err = reg.remove(url)
		if err != nil {
			log.Panicln(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return

	case http.MethodGet:
		reg.mutex.Lock()
		defer reg.mutex.Unlock()

		var b bytes.Buffer
		enc := json.NewEncoder(&b)
		err := enc.Encode(reg.registrations)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-type", "application/json")
		w.Write(b.Bytes())

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

}
