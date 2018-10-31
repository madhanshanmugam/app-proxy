package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/gcfg.v1"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// Config section stats here
var config Config

const appEndPoint = "localhost:9000"

type ServiceConfig struct {
	RestartCmd      string
	HealthCheckUrl  string
	ServiceEndpoint string
	ServiceDomain   string
	LogPath         string
}
type Config struct {
	Service map[string]*ServiceConfig
}

// Config section ends here

// Proxy section starts here
func getTargetHost(req *http.Request) (*url.URL, error) {
	for _, v := range config.Service {
		if req.Host == v.ServiceDomain {
			return url.Parse(v.ServiceEndpoint)
		}
	}
	return url.Parse(appEndPoint)
}

// for logging the response
type proxyTransport struct {
}

func (t *proxyTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	response, err := http.DefaultTransport.RoundTrip(request)

	if strings.Contains(response.Header.Get("Content-Type"), "application/json") {
		body, err := httputil.DumpResponse(response, true)
		if err != nil {
			return nil, err
		}
		log.Print(string(body))
	}
	return response, err
}

func appProxy(target *url.URL) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		requestDump, err := httputil.DumpRequest(req, true)
		if nil != err {
			panic("error in the request dumping")
		}
		//todo have to set the host header from the config
		fmt.Println(string(requestDump))
		targetUrl, _ := getTargetHost(req)
		req.URL.Scheme = targetUrl.Scheme
		req.URL.Host = targetUrl.Host //config.Service["app"].ServiceEndpoint
		req.URL.Path = req.URL.Path
	}
	return &httputil.ReverseProxy{Director: director}
}

// Proxy section ends here

func main() {
	cfgStr := "config.gcfg"
	err := gcfg.ReadFileInto(&config, cfgStr)
	if err != nil {
		log.Fatalf("Failed to parse gcfg data: %s", err)
	}
	router := mux.NewRouter()
	origin, _ := url.Parse("")
	reverseProxy := appProxy(origin)
	reverseProxy.Transport = &proxyTransport{} // Custom transport is needed for purpose like logging the response, measuring the latency etc
	router.PathPrefix("/").Handler(reverseProxy);
	error := http.ListenAndServe(appEndPoint, router)
	fmt.Println(error)
}
