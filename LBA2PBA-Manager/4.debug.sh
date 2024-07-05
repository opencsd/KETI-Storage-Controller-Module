#!/bin/bash 

while [ -z $PODNAME ]
do
    PODNAME=`kubectl get po -o=name -A --field-selector=status.phase=Running | grep lba2pba-manager`
    PODNAME="${PODNAME:4}"
done

kubectl logs -f $PODNAME -n storage-platform



