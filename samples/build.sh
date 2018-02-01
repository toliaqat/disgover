#!/bin/bash

eval $(minikube docker-env)

kubectl delete deployment --all
kubectl delete pods --all

cd node1
imageId=$(docker build -t disgover:v1 . | grep 'Successfully built ' | sed "s/Successfully built //g")
docker tag $imageId localhost:5000/disgover:v1
docker push localhost:5000/disgover:v1
kubectl run disgover-node1 --image=localhost:5000/disgover:v1 --port=9001 --image-pull-policy=Never


cd ../node2
imageId=$(docker build -t disgover:v2 . | grep 'Successfully built ' | sed "s/Successfully built //g")
docker tag $imageId localhost:5000/disgover:v2
docker push localhost:5000/disgover:v2
kubectl run disgover-node2 --image=localhost:5000/disgover:v2 --port=9001 --image-pull-policy=Never

cd ../node3
imageId=$(docker build -t disgover:v3 . | grep 'Successfully built ' | sed "s/Successfully built //g")
docker tag $imageId localhost:5000/disgover:v2
docker push localhost:5000/disgover:v3
kubectl run disgover-node3 --image=localhost:5000/disgover:v3 --port=9001 --image-pull-policy=Never

kubectl describe pod disgover-node1 | grep -e IP -e Port
kubectl describe pod disgover-node2 | grep -e IP -e Port
kubectl describe pod disgover-node3 | grep -e IP -e Port
