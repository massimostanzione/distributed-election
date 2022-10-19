#!/bin/bash
# run this script if something goes wrong
echo -E "Cleaning up..."
sudo rm -f go.mod go.sum

for item in node serviceRegistry
do
    sudo rm -f pb/$item/*.go \
	           bin/$item
done
sudo rmdir bin
echo -E "Done."
