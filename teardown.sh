#!/bin/bash
# run this script if something goes wrong
echo -E "Cleaning up..."
sudo rm -f go.mod go.sum verbose*

modules=( node serviceregistry )

for item in "${modules[@]}"
do
    sudo rm -f bin/$item        \
               $item/cmd/cmd    \
	           $item/pb/*.go
done
sudo rmdir bin
echo -E "Done."
