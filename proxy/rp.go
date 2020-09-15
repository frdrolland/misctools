/*
Copyright © 2020 Frederic ROLLAND <frolland.2020@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package proxy

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

/*
	Structs
*/

type requestPayloadStruct struct {
	ProxyCondition string `json:"proxy_condition"`
}

var (
	// target URL for reverse proxying
	targetURL string
)

/*
	Utilities
*/

// Get env var or default
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

/*
	Getters
*/

// Get the port to listen on
func getListenAddress() string {
	viper.SetDefault("proxy.listen", "8080")
	port := viper.GetString("proxy.listen")
	return ":" + port
}

/*
	Logging
*/
// Log the typeform payload and redirect url
func logRequestPayload(requestionPayload requestPayloadStruct, proxyURL string) {
	log.Printf("proxy_condition: %s, proxy_url: %s\n", requestionPayload.ProxyCondition, proxyURL)
}

func readResponseBody(resp *http.Response) (data []byte, err error) {
	data, err = ioutil.ReadAll(resp.Body) //Read html
	if err != nil {
		return nil, err
	}
	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}
	//b = bytes.Replace(b, []byte("server"), []byte("schmerver"), -1) // replace html
	body := ioutil.NopCloser(bytes.NewReader(data))
	resp.Body = body
	resp.ContentLength = int64(len(data))
	resp.Header.Set("Content-Length", strconv.Itoa(len(data)))
	return data, nil
}

// Serve a reverse proxy for a given url
func serveReverseProxy(target string, res http.ResponseWriter, req *http.Request) {
	// parse the url
	url, _ := url.Parse(target)

	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)

	// Update the headers to allow for SSL redirection
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = url.Host

	// Note that ServeHttp is non blocking and uses a go routine under the hood
	proxy.ServeHTTP(res, req)
}

// Get a json decoder for a given requests body
func readRequestBody(request *http.Request) []byte {
	// Read body to buffer
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		panic(err)
	}

	// Because go lang is a pain in the ass if you read the body then any susequent calls
	// are unable to read the body again....
	request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return body
}

// Given a request send it to the appropriate url
func handleRequestAndRedirect(res http.ResponseWriter, req *http.Request) {
	body := readRequestBody(req)
	log.WithFields(log.Fields{
		"url": "req.url",
	}).Info(body)

	serveReverseProxy(targetURL, res, req)
}

//
// StartReverseProxy starts a http reverse proxy
//
func StartReverseProxy() {
	// Log setup values
	targetURL = viper.GetString("proxy.target.url")

	log.Printf("Server will run on: %s", getListenAddress())
	log.Printf("Redirecting to Default url: %s", targetURL)

	// start the server
	http.HandleFunc("/", handleRequestAndRedirect)
	if err := http.ListenAndServe(getListenAddress(), nil); err != nil {
		panic(err)
	}
}
