build:
	go build -i -v

fmt:
	go fmt

run: build
	sudo ./sample-nginx-http2 \
		-port=443 \
		-key=./certs/key_example.com.pem \
		-cert=./certs/cert_example.com.pem

open-chrome:
	SSLKEYLOGFILE=~/tlskey.log \
        /Applications/Google\ Chrome.app/Contents/MacOS/Google\ Chrome &

.PHONY: build fmt run open-chrome
