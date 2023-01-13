# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM golang:1.18
WORKDIR /build
COPY . .
ENV CGO_ENABLED=0 
ENV GOOS=linux
ARG TARGETOS TARGETARCH
RUN cd cmd/agent; GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /runner -mod mod -a .

FROM zcube/bitnami-compat-kubectl:1.26.0
USER root
RUN apt-get update && apt-get install -y wget
RUN mkdir -p $HOME/.kubepug  && \
    cd $HOME/.kubepug && \
    wget https://github.com/rikatz/kubepug/releases/download/v1.1.3/kubepug_linux_amd64.tar.gz && \
    tar -zxvf kubepug_linux_amd64.tar.gz && \
    rm -f kubepug_linux_amd64.tar.gz
ENV PATH="$HOME/.kubepug/:${PATH}"
WORKDIR /root/
COPY --from=0 /runner /bin/runner
USER 1001
ENTRYPOINT ["/bin/runner"]
