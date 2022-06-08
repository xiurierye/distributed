package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
)

func RegisterService(r Registration) error {
	serveicUpdateURL, err := url.Parse(r.ServiceUpdateURl)
	if err != nil {
		return err
	}

	http.Handle(serveicUpdateURL.Path, &serviceUpdateHandler{})

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err = enc.Encode(r)

	if err != nil {
		return err
	}

	resp, err := http.Post(ServicesURL, "application/json", buf)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to register service. Registry serice responded with code %d", resp.StatusCode)
	}

	return nil
}

type serviceUpdateHandler struct{}

func (suh serviceUpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	var pat patch
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&pat)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Printf("Updated receive %v\n", pat)

	prov.Update(pat)

}

func ShutdowService(url string) error {

	req, err := http.NewRequest(http.MethodDelete, ServicesURL, bytes.NewBuffer([]byte(url)))
	if err != nil {
		log.Panicln(err)
	}
	req.Header.Add("Content-Type", "text/plain")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Panicln(err)
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to deregister service. Registry"+
			"server responed with code %v", res.StatusCode)
	}

	return nil

}

type providers struct {
	services map[ServiceName][]string
	mutex    *sync.RWMutex
}

var prov = providers{
	services: make(map[ServiceName][]string),
	mutex:    new(sync.RWMutex),
}

func (p *providers) Update(pat patch) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, patchEntry := range pat.Added {
		if _, ok := p.services[patchEntry.Name]; !ok {
			p.services[patchEntry.Name] = make([]string, 0)
		}
		p.services[patchEntry.Name] = append(p.services[patchEntry.Name], patchEntry.URL)
	}

	for _, patchEntry := range pat.Removed {
		if provideURLs, ok := p.services[patchEntry.Name]; ok {
			for i := range provideURLs {
				if provideURLs[i] == patchEntry.URL {
					p.services[patchEntry.Name] = append(provideURLs[:i],
						provideURLs[i+1:]...)
				}
			}
		}
	}
}

// 本来是返回 服务对应的多个url 单当前服务很简单  ,    只有一个url提供服务就 偷懒直接返回string
func (p providers) get(name ServiceName) (string, error) {
	if providers, ok := p.services[name]; ok {
		return providers[int(rand.Float32()*float32(len(providers)))], nil
	}

	return "", fmt.Errorf("No providers available for service %v", name)

}

func GetProvoider(name ServiceName) (string, error) {
	return prov.get(name)
}
