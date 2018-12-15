FROM golang:latest
RUN mkdir /app
ADD . /app/
WORKDIR /app
EXPOSE :8089
RUN go get "github.com/prometheus/client_golang/prometheus"
RUN go get "github.com/prometheus/client_golang/prometheus/promhttp"
RUN go get "github.com/sausheong/hs1xxplug"

RUN go build -o main .
ENTRYPOINT ["/app/main"]
