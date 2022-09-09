#!/bin/bash
img=192.168.9.12:5000/kubespray:test
docker build -t $img .
docker push $img