#!/bin/bash
# generate server, client and worker code using the protocol buffer compiler
#modules=( client server worker )
modules=( node serviceRegistry )
for item in "${modules[@]}"
do
    echo "Generating $item stub..."
    cd pb
    #cd $item
    protoc --go_out=. \
           --go_opt=paths=source_relative \
           --go-grpc_out=. \
           --go-grpc_opt=paths=source_relative \
           $item/proto$item.proto
    echo "Done."
    echo ""
    cd ..
done

# arrange dependencies
echo "Arranging module with relative dependencies..."
go mod init bully
go mod tidy
echo "Done."
echo ""    
    
# compile entities
for item in "${modules[@]}"
do
    cd cmd/$item
    echo "Building $item..."
    go build -v -o ../../bin/$item
    echo "Done."
    echo ""
    cd ..
    cd ..
done

