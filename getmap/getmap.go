package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

var apiKey = os.Getenv("API_KEY")

// loadMap generates a new map or returns a cached one
func loadMap(address string, title string) ([]byte, error) {
	mapHost := "http://maps.googleapis.com/maps/api/staticmap"

	v := url.Values{}
	v.Add("center", address)
	v.Add("zoom", "15")
	v.Add("scale", "2")
	v.Add("size", "400x350")
	v.Add("markers", address)
	v.Add("sensor", "false")
	v.Add("key", apiKey)
	mapURL := mapHost + "?" + v.Encode()

	filename := title + ".png"
	for i := 0; i < 30; i++ {
		body, err := ioutil.ReadFile(filename)
		if err == nil {
			return body, nil
		}

		if i == 0 {
			fmt.Println("Getting map for " + address + " and saving it to " + filename)
			fmt.Println("Map URL: " + mapURL)
			httpClientGetMap, _ := http.Get(mapURL)
			body, _ := ioutil.ReadAll(httpClientGetMap.Body)
			ioutil.WriteFile(filename, []byte(body), 0600)
		}

		i++
	}

	return nil, errors.New("Failed to get an image")
}

// viewHandler returns the map image
func viewHandler(w http.ResponseWriter, r *http.Request) {

	address := r.URL.Query().Get("address")

	md5 := md5.Sum([]byte(address))
	title := hex.EncodeToString(md5[:])
	p, err := loadMap(address, title)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusGatewayTimeout)
		return
	}

	w.Header().Set("content-type", "image/png")
	w.Write(p)
}

func main() {
	http.HandleFunc("/getmap", viewHandler)

	http.ListenAndServe(":8080", nil)
}
