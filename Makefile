bin_dir = /Users/j/pika/yfncc-cdist/cdist/conf/manifest/bin/openbsd
html_dir = /Users/j/pika/yfncc-cdist/cdist/conf/manifest/html
cdist = /Users/j/pika/yfncc-cdist/bin/cdist

warn:
	@echo "read the README.md first."

build :
	go get -u golang.org/x/crypto/acme/autocert
	env GOOS=openbsd GOARCH=amd64 go build

deploy : build
	cp simple-https-server $(bin_dir)/simple-https-server
	$(cdist) config -v pika-web.mit.edu
