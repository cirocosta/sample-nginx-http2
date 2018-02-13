package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/cirocosta/sample-nginx-http2/assets"
)

var (
	port  = flag.Int("port", 8443, "port to listen to")
	http2 = flag.Bool("http2", true, "whether to use http2 or not")
	cert  = flag.String("cert", "", "certificate for TLS")
	key   = flag.String("key", "", "key for TLS")
)

func must(err error) {
	if err == nil {
		return
	}

	panic(err)
}

// handleIndex is the handler for serving `/`.
//
// It first checks if it's possible to push contents via
// the connection. If so, then it pushes `/image.svg` such
// that at the same moment  that the browser is fetching
// `index.html` it can also start retrieving `image.svg`
// (even before it knows about the existence in the html).
func handleIndex(w http.ResponseWriter, r *http.Request) {
	var err error

	if *http2 {
		pusher, ok := w.(http.Pusher)
		if ok {
			must(pusher.Push("/image.svg", nil))
		}
	} else {
		w.Header().Add("Link", "</proxy/image.svg>; rel=preload; as=image")
	}

	w.Header().Add("Content-Type", "text/html")
	_, err = w.Write(assets.Index)
	must(err)
}

// handleImage is the handler for serving `/image.svg`.
//
// It does nothing more than taking the byte array that
// defines our SVG image and sending it downstream.
func handleImage(w http.ResponseWriter, r *http.Request) {
	var err error

	w.Header().Set("Content-Type", "image/svg+xml")
	_, err = w.Write(assets.Image)
	must(err)
}

// main provides the main execution of our server.
//
// It makes sure that we're providing the required
// flags: cert and key.
//
// These two flags are extremely important because
// browsers will only communicate via HTTP2 if we
// serve the content via HTTPS, meaning that we must
// be able to properly terminate TLS connections, thus,
// need a private key and a certificate.
func main() {
	flag.Parse()

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/image.svg", handleImage)

	if *http2 {
		if *key == "" {
			fmt.Println("flag: key must be specified")
			os.Exit(1)
		}

		if *cert == "" {
			fmt.Println("flag: cert must be specified")
			os.Exit(1)
		}

		fmt.Printf("Server HTTP2 on :%d\n", *port)
		must(http.ListenAndServeTLS(":"+strconv.Itoa(*port), *cert, *key, nil))
		return
	}

	fmt.Printf("Server HTTP1.1 on :%d\n", *port)
	must(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}
