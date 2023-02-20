# /bin/bash
VERSION=v1.0.7

# 卸载之前的部署
. ./offload.sh

# 创建target文件夹
mkdir target
cp deployment.yaml target

# 修改deployment.yaml
sed -i "s/CUP_VERSION/${VERSION}/g" target/deployment.yaml

# 构建docker镜像
docker run --privileged --rm tonistiigi/binfmt --install all
docker build . -t evolonation/patient-edge:${VERSION} --platform arm64 --push

# 创建configmap
kubectl create configmap edge-mysql-scripts --from-file script_edge.sql 
kubectl create configmap cloud-mysql-scripts --from-file script_cloud.sql 

# 部署
kubectl apply -f target/deployment.yaml

# 清除残余
rm -r target



