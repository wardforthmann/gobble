# gobble
Gobble is a tool intended to help test any application that makes http POST requests. It will help you see exactly what
 is being POSTed to make sure it is correct and can help reduce the amount of time spent debugging failing requests.
 There are several applications that already do this around the web but they leave your request data publicly accessible.
 Gobble is self hosted which means your data stays safe inside your own network in situations where you don't want to
 share it with the world.
 
 
### Usage
The current functionality is pretty simple. Simply POST any message to root of the gobble server. It will take your
 request and store it to disk. You can then look up this request by navigating to its location on the gobble server.
 
#### Submitting requests

 To run a gobble server:
 
 ```
 ./gobble -port 8080 -dir ./public
 ```

 Then point your POST request at the root of the gobble server. The request will have its headers and body written to disk 
 under the specified home directory. To ease navigation, the request will be written to a sub directory corresponding
 to the date that the request was received. Gobble will return a string with the path that your request was stored at.
 You can use this to retrieve your request at a later date. There are also several query parameters you may use to override
 some of the default behavior.
 
 Query Parameter | Description 
 :--------------:|:-----------
 dir             | Override the name of the directory your request is written to. This can be used to group requests in a single location.
 no_header       | Do not return the headers at the top of the response. This will try to set the content-type of the response based on the stored header if possible.
 status_code     | Force the http status code on the response. Useful for testing error conditions.
 
 An example post would be:
 ```
 curl -X POST -d '{"sample": "text"}' localhost:8080?dir=test
 ```
 
#### Retrieving requests
 To view stored POST requests navigate a web browser to the location the gobble server is running. It will display a list
 of available directories. You will find your store requests inside these directories. Click on the link corresponding to
 your request and the request will be returned with its headers and body.
 
 
#### Command line args
There are several command line arguments available when running gobble:

 ```text
  -dir string
    	Specifies the root directory which all directories and requests will be stored under (default "public")
  -port string
    	Specifies the port to listen for incoming connections (default "80")
  -tls
    	Tells gobble to listen for secure connections (ie. https)
  -tlsCert string
    	Specifies the path to the x509 certificate (default "cert.pem")
  -tlsKey string
    	Specifies the path to the private key corresponding to the x509 certificate (default "key.pem")
  -tlsPort string
    	Specifies the port to listen for incoming secure connections (default "443")
```

#### Dependencies
Gobble only depends on a single library named chi. It is currently vendored for your convenience but you can find the 
original repo [here](https://github.com/pressly/chi).
