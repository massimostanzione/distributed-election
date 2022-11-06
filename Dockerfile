# syntax=docker/dockerfile:1
FROM golang:1.19-alpine
ADD . /app
WORKDIR /app
COPY . .

RUN apk add protobuf curl
#RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
#PB_REL="https://github.com/protocolbuffers/protobuf/releases"
RUN curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v3.15.8/protoc-3.15.8-linux-x86_64.zip
RUN unzip protoc-3.15.8-linux-x86_64.zip -d $HOME/.local
RUN export PATH="$PATH:$HOME/.local/bin"
RUN protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative serviceregistry/pb/protoserviceregistry.proto

RUN go mod init distributedelection
RUN go mod tidy


#RUN go mod download
RUN cd serviceregistry/cmd
RUN go build -v -o ../../bin/serviceregistry
ENTRYPOINT ["/home/app/bin/serviceregistry"]










#------------------------
# syntax=docker/dockerfile:1
#FROM golang:1.19-alpine

# update dependencies
#RUN apk update

#WORKDIR /home/foo/app
#ADD . /bin
#COPY . .
#WORKDIR /APP
# build
#RUN chmod -R 777 .
#RUN ["chmod", "+x", "/home/foo/app/bin/node"]
#RUN ["chmod", "+x", "/home/foo/app/bin/serviceregistry"]
#CMD ["chmod", "+x", "/home/foo/app/bin/node"]
#CMD ["chmod", "+x", "/home/foo/app/bin/serviceregistry"]

#RUN chmod +x -R .
#RUN cd bin
#RUN chmod +x -R .
#RUN chmod +x -R serviceregistry
#RUN ./serviceregistry
#CMD ["./serviceregistry"]#, "-a", "b"]
