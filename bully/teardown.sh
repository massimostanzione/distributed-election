#!/bin/bash
# run this script if something goes wrong
echo -E "Cleaning up..."
sudo rm -f bin
sudo rm -f go.mod go.sum

modules=( node serviceregistry )

for item in "${modules[@]}"
do
    sudo rm -f $item/cmd/$item  \
	           $item/pb/*.go
done
#sudo rmdir bin
echo -E "Done."
