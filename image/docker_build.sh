#!/bin/bash
docker rm -f kubespray-$1
# docker rmi kubespray:$1
docker build -t kubespray:$1 .
# docker push kubespray:$1
# docker run -it --name kubespray-$1 kubespray:$1 /bin/sh
