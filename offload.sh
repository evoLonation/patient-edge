# /bin/bash
kubectl delete -f deployment.yaml
kubectl delete configmap edge-mysql-scripts
kubectl delete configmap cloud-mysql-scripts