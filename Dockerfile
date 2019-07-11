FROM golang:latest

ARG BUILD_DATE
ARG VCS_REF

LABEL org.label-schema.build-date=$BUILD_DATE \
      org.label-schema.vcs-url="https://github.com/Zibby/tplink-exporter.git" \
      org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.schema-version="1.0.0-rc1"


RUN mkdir /app
ADD . /app/
WORKDIR /app
EXPOSE 8089/tcp
RUN go get "github.com/prometheus/client_golang/prometheus"
RUN go get "github.com/prometheus/client_golang/prometheus/promhttp"
RUN go get "github.com/sausheong/hs1xxplug"
RUN go get "github.com/sirupsen/logrus"

RUN go build -o main .
ENTRYPOINT ["/app/main"]

ARG GIT_COMMIT=unspecified
LABEL git_commit=$GIT_COMMIT
