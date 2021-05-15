FROM golang:1.16
ENV GOPROXY="https://goproxy.cn,direct"


WORKDIR /src
COPY . .

RUN go build .

FROM nvcr.io/nvidia/cuda:10.2-base-ubuntu18.04

COPY --from=0 /src/k8s-metrics-logger .

ENTRYPOINT ["./k8s-metrics-logger"]