FROM golang:1.16
ENV GOPROXY="https://goproxy.cn,direct"

WORKDIR /src
COPY . .

RUN go build .

FROM ubuntu

COPY --from=0 /src/k8s-metrics-logger .

ENTRYPOINT ["./k8s-metrics-logger"]