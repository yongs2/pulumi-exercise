#!/bin/bash

# for log
export INSTALL_LOG=/var/log/k8s_install_2.log
echo "START SCRIPT >>>>>>> INSTALL_LOG=${INSTALL_LOG}"
sudo touch $INSTALL_LOG
sudo chown $USER $INSTALL_LOG

echo "-------------USER=[${USER}], MASTER ONLY" &>> $INSTALL_LOG
echo "wait a moment, sleep 10s"
sleep 10

# master only
mkdir -p $HOME/.kube &>> $INSTALL_LOG
sudo cp -f /etc/kubernetes/admin.conf $HOME/.kube/config &>> $INSTALL_LOG
sudo chown $(id -u):$(id -g) $HOME/.kube/config &>> $INSTALL_LOG
echo "KUBECONFIG="; cat $HOME/.kube/config
sed -i 's/127.0.0.1/'$(hostname -i)'/g' $HOME/.kube/config &>> $INSTALL_LOG

echo 'source <(kubectl completion bash)' >>$HOME/.bashrc
echo 'alias k=kubectl' >>$HOME/.bashrc
echo 'complete -o default -F __start_kubectl k' >>$HOME/.bashrc

echo ">>> remove taints in all nodes" &>> $INSTALL_LOG
kubectl taint nodes --all node-role.kubernetes.io/control-plane- &>> $INSTALL_LOG
kubectl taint nodes --all node-role.kubernetes.io/master- &>> $INSTALL_LOG
# role 에 master 추가
kubectl label nodes k8s-master-01 node-role.kubernetes.io/master= &>> $INSTALL_LOG

echo ">>> get taints in k8s-master-01" &>> $INSTALL_LOG
kubectl describe node k8s-master-01 | grep Taints &>> $INSTALL_LOG
kubectl get node -o wide &>> $INSTALL_LOG

# for root
sudo mkdir -p /root/.kube &>> $INSTALL_LOG
sudo cp -f /etc/kubernetes/admin.conf /root/.kube/config &>> $INSTALL_LOG
sudo sh -c 'echo "source <(kubectl completion bash)" >> /root/.bashrc'
sudo sh -c 'echo "alias k=kubectl" >>/root/.bashrc'
sudo sh -c 'echo "complete -o default -F __start_kubectl k" >>/root/.bashrc'

echo "END SCRIPT <<<<<<<<<<<<<<<<"
# end of script.
