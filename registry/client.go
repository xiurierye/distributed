package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func RegisterService(r Registration) error {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err := enc.Encode(r)

	if err != nil {
		return err
	}

	resp, err := http.Post(ServicesURL, "application/json", buf)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to register service. Registry serice "+
			"responded with code %v", resp.StatusCode)
	}

	return nil
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
