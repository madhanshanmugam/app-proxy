App proxy is a forward proxy, it will forward all the requests to the configured target(s) and logs the request and response as well. This majorly helps in debugging outgoing API calls without the need to enable logging at the source code. we can proxy more than one target at a time.

This can also be configured as revere proxy and does the logging. Since the proxy is based on simple domain based routing, it doesn't require any ssl certificates to be installed on the client machine.


<h4>To proxy a domain </h4>


**Using Docker**

 `docker run --rm -it -p 9000:9000  madhanshamugam/app-proxy domainNameWithProtocol`
 
  example 
  
 `docker run --rm -it -p 9000:9000  madhanshamugam/app-proxy https://stackoverflow.com`
 
 and use the endpoint  [http://localhost:9000](url) 
 
 **Using go build**
 
 clone the repo and cd to the parent directory. 
    `git clone git@github.com:madhanshanmugam/app-proxy.git`
  building the service `go build main/proxy.go`
  To run it `./proxy domainNameWithProtocol` ex `./proxy https://stackoverflow.com`
  
    
  <h4>For proxing multiple domains </h4>
   The config file `config.gcfg` needs to be edited and can be added as many services. After modifying service needs to started without the `domainNameWithProtocol` param. To route the domain requests to the proxy, /etc/hosts can be used.
    
 
 