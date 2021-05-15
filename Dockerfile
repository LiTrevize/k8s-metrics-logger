FROM golang:1.16
ENV GOPROXY="https://goproxy.cn,direct"


WORKDIR /src
COPY . .

RUN go build .

FROM nvcr.io/nvidia/cuda:11.2.1-base-ubuntu20.04

COPY --from=0 /src/k8s-metrics-logger .

ENTRYPOINT ["./k8s-metrics-logger"]