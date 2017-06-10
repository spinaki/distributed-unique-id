FROM golang
MAINTAINER Pinaki Sinha <spinaki@gmail.com>
# env GOOS=linux GOARCH=amd64 go build -v
COPY ./uniqueidgenerator /go/bin/uniqueidgenerator
CMD ["/go/bin/uniqueidgenerator"]
EXPOSE 8080
