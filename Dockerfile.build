FROM golang
RUN mkdir -p /go/src/cb
WORKDIR /go/src/cb
COPY * ./
RUN go get -v
RUN go build -v -o cb
