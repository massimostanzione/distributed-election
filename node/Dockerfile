# syntax=docker/dockerfile:1
FROM golang:1.19

# update dependencies
RUN apt-get update

# copia la *intera* repo, per le interdipendenze
ADD .
WORKDIR /
COPY    .. /APP

# build
#RUN ./../setup.sh

#CMD ["./../bin/node", "-a", "b"]
