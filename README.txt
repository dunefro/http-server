# Application
1. `server.go` runs a HTTP/HTTPS server on the port 9090.
2. It gives you 4 endpoints 
```
/
/healthz
/unprotected
/protected
```
3. `/protected` requires username and password that is sourced from a file called `creds.txt`
4. `/unprotected` will give you back some data, to test in wireshark
5. You can run this server in TLS and non-TLS mode to test the connection.

# How to run
1. `go run server.go`
2. If the application is running in `server.go` as `srv.ListenAndServeTLS` then you can access it via `https`
3. Open wireshark and try to access the `Authorization: Basic <base64>` request for `/protected` handler
    a. With non-TLS you can decode the header by `echo <base64> | base64 -d`
    b. For TLS you will not be able to do that
