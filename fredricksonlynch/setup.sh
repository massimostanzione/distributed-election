#!/bin/bash
# generate node and serviceregistry code using the protocol buffer compiler
modules=( node serviceregistry )
for item in "${modules[@]}"
do
    echo "Generating $item stub..."
    protoc --go_out=.                           \
           --go_opt=paths=source_relative       \
           --go-grpc_out=.                      \
           --go-grpc_opt=paths=source_relative  \
           $item/pb/proto$item.proto
    echo "Done."
    echo ""
done

# arrange dependencies
echo "Arranging module with relative dependencies..."
go mod init fredricksonlynch
go mod tidy
echo "Done."
echo ""    
    
# compile entities
for item in "${modules[@]}"
do
    echo "Building $item..."
    cd $item/cmd
    go build -v -o ../../bin/$item
    cd ../..
    echo "Done."
    echo ""
done

