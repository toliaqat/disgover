#!/bin/bash

eval $(minikube docker-env)

echo
echo
echo ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~
echo Delete PODs
echo ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~
kubectl delete deployment --all
kubectl delete pods --all
sleep 1m

echo
echo
echo ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~
echo Build NODE-1
echo ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~
imageTag=node1
cd $imageTag
docker build -t disgover:$imageTag . | tee docker-output
imageId=$(cat docker-output | grep 'Successfully built ' | sed "s/Successfully built //g")
rm docker-output
docker tag $imageId localhost:5000/disgover:$imageTag
docker push localhost:5000/disgover:$imageTag
kubectl run disgover-$imageTag --image=localhost:5000/disgover:$imageTag --port=1975 --image-pull-policy=Never
sleep 30s
seedNodeIP=$(kubectl describe pod disgover-node1 | grep -e IP | sed "s/IP://g" | sed 's/ //g')
echo SeedNodeIP is $seedNodeIP
read -p "CHECK Line above and [Enter] if looks ok..."

echo
echo
echo ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~
echo Build NODE-2
echo ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~
imageTag=node2
cd ../$imageTag
docker build -t disgover:$imageTag . | tee docker-output
imageId=$(cat docker-output | grep 'Successfully built ' | sed "s/Successfully built //g")
rm docker-output
docker tag $imageId localhost:5000/disgover:$imageTag
docker push localhost:5000/disgover:$imageTag
kubectl run disgover-$imageTag --image=localhost:5000/disgover:$imageTag --port=1975 --image-pull-policy=Never --env="SEED_NODE_IP=$seedNodeIP"

echo
echo
echo ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~
echo Build NODE-3
echo ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~
imageTag=node3
cd ../$imageTag
docker build -t disgover:$imageTag . | tee docker-output
imageId=$(cat docker-output | grep 'Successfully built ' | sed "s/Successfully built //g")
rm docker-output
docker tag $imageId localhost:5000/disgover:$imageTag
docker push localhost:5000/disgover:$imageTag
kubectl run disgover-$imageTag --image=localhost:5000/disgover:$imageTag --port=1975 --image-pull-policy=Never --env="SEED_NODE_IP=$seedNodeIP"

echo
echo
echo ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~
echo List IPs and Ports
echo ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~  ~~~~
cd ..
sleep 5s
kubectl describe pod disgover-node1 | grep -e IP -e Port
kubectl describe pod disgover-node2 | grep -e IP -e Port
kubectl describe pod disgover-node3 | grep -e IP -e Port
