# gobble
gobble is a tool intended to help test any application
 that makes http POST requests.
 
###Usage
The current functionality is pretty simple. Simply POST any message to root of the gobble server. It will take your
 request and store it to disk. 
 
####Submitting requests
 Point your POST request at the root of the gobble server. The request will have its headers and body written to disk 
 under the specified home directory. To ease navigation, the request will be written to a sub directory corresponding
 to the date that the request was received. You may provide a query parameter named "dir" to override this with a name
 of your choosing. Gobble will return a string with the path that your request was stored at. You can use this to retrieve
 your request at a later date.
 
####Retrieving requests
 To view stored POST requests navigate a web browser to the location the gobble server is running. It will display a list
 of available directories. You will find your store requests inside these directories. Click on the link corresponding to
 your request and the request will be returned with its headers and body.
 
 
####Command line args
There are several command line arguments available when running gobble:

 ```text
  -dir string
    	Specifies the root directory which all directories and requests will be stored under (default "public")
    
  -port string
      	Specifies the port to listen for incoming connections (default "80")
```

####Dependencies
Gobble only depends on a single library named chi. You can retrieve this by using the go get tool.
 ```text
go get github.com/pressly/chi
```