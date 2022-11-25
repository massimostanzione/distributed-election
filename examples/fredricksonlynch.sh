#!/bin/bash
# a quick, ready-to-use demo

NODES_NO=12
BASE_ADDRESS=40042

echo "*** distributed-election demo ***"
echo "   fredrickson-lynch algorithm   "
echo -E "" 
if ! [ -f "./../go.mod" ]; then 
    echo "Please run ./setup.sh first"
    exit -1
fi
echo "This script will run a demo of the 'distributed-election' program"
echo "The following entities will be run, in separate windows:"
echo "- a Service Registry"
echo "- $NODES_NO nodes"
echo "Each entity will be identified by the terminal window title."
echo -E ""
read -n 1 -s -r -p "Press any key to start..."
echo -E ""
echo -E "-------------------------------------"

# detect default terminal
terms=($TERMINAL x-terminal-emulator urxvt rxvt termit terminator xterm xfce4-terminal xdg-terminal gnome-terminal iterm)
for t in ${terms[*]}
do
    if [ $(command -v $t) ]
    then
        detected_term=$t
        break
    fi
done

# run service registry
$detected_term -T "ServReg" -e "bash -c './../bin/serviceregistry;bash'"   

#run nodes
for i in $(seq 1 $NODES_NO)
do
    $detected_term -T "Node n. $i" -e "bash -c './../bin/node -c ../configs/config.ini -a fl -ncl ABSENT -p $(($BASE_ADDRESS + $i));bash'"    
done
