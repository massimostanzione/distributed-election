# syntax=docker/dockerfile:1
FROM golang:1.19-alpine AS build

ADD . /distributed-election
WORKDIR /distributed-election

# install dependencies
RUN apk add protobuf curl unzip --no-cache

# add support for gRPC and protocol buffer, not natively included in golang-alpine
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

RUN curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v3.15.8/protoc-3.15.8-linux-x86_64.zip && \
    unzip protoc-3.15.8-linux-x86_64.zip -d $HOME/.local && \
    export PATH="$PATH:$HOME/.local/bin"

# copy project folder
COPY . .

# compile stubs
RUN protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative serviceregistry/pb/protoserviceregistry.proto && \
    protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative node/pb/protonode.proto
    
RUN go mod init distributedelection
RUN go mod tidy

# compile node
WORKDIR /distributed-election/node/cmd
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o /distributed-election/bin/node

# new stage: drastically reduce image size by using only strictly necessary executable(s) and file(s)
FROM scratch
COPY --from=build /distributed-election/bin/ /distributed-election
WORKDIR /distributed-election
ENTRYPOINT ["./node", "-a", "b", "-sh", "serviceregistry"]
