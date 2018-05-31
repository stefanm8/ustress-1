package main

import "fmt"

func main() {
	codes := []int{0, 50, 101, 204, 307, 404, 500, 600, 660}

	for _, code := range codes {
		if code > 100 && code < 600 {
			fmt.Printf("true %d\n", code)

		} else {
			fmt.Printf("false %d\n", code)
		}
	}
}

// import (
// 	"context"
// 	"crypto/tls"
// 	"log"
// 	"net"
// 	"net/http"
// 	"time"
// )

// func main() {

// 	address := "10.29.30.8:443"
// 	tr := &http.Transport{}

// 	tr = &http.Transport{
// 		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
// 	}
// 	log.Println(address)

// 	dialer := &net.Dialer{
// 		Timeout:   30 * time.Second,
// 		KeepAlive: 30 * time.Second,
// 		DualStack: true,
// 	}
// 	tr.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
// 		log.Println(address)
// 		log.Println(addr)
// 		return dialer.DialContext(ctx, network, address)
// 	}

// 	client := http.Client{
// 		Timeout:   time.Duration(2 * time.Second),
// 		Transport: tr,
// 	}

// 	resp, err := client.Get("https://idam-pp.metrosystems.net/.well-known/openid-configuration")
// 	if err != nil {
// 		log.Panicln(err.Error())
// 	}
// 	log.Println(resp.Header)
// 	// data, _ := ioutil.ReadAll(resp.Body)
// 	// log.Println(string(data))

// }