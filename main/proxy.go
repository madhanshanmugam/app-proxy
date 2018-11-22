package main

import (
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
	ServiceTargetEndpoint string
	ServiceLocalDomain   string
	LogPath         string
}
type Config struct {
	Service map[string]*ServiceConfig
}

// Config section ends here

// Proxy section starts here
func getTargetHost(req *http.Request) (*url.URL, error) {
	host := strings.Split(req.Host,":")[0]
	for _, v := range config.Service {
		if host == v.ServiceLocalDomain {
			return url.Parse(v.ServiceTargetEndpoint)
		}
	}
	return url.Parse(appEndPoint)
}

// for logging the response
type proxyTransport struct {
}

func (t *proxyTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	response, err := http.DefaultTransport.RoundTrip(request)

	if nil == response {
		log.Print(strings.Join([]string{"ERROR Endpoint : ", request.URL.Host, " could be down"},""))
	} else if strings.Contains(response.Header.Get("Content-Type"), "application/json") {
		body, err := httputil.DumpResponse(response, true)
		if err != nil {
			return nil, err
		}
		log.Println("------------------------Start of Response-------------")
		log.Println(string(body))
		log.Println("------------------------End of Response-------------")
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
		log.Println("------------------------Start of Request-------------")
		log.Println(string(requestDump))
		log.Println("------------------------End of Request-------------")
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
	log.Println(error)
}
