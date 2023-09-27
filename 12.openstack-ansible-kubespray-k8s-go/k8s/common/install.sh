#!/bin/bash

# for log
export INSTALL_LOG=/var/log/k8s_install_0.log
echo "START SCRIPT >>>>>>> INSTALL_LOG=${INSTALL_LOG}, HOME=${HOME}"
sudo touch $INSTALL_LOG
sudo chown $USER $INSTALL_LOG

echo "wait a moment, sleep 10s"
sleep 10
echo ">> upgrade" &>> $INSTALL_LOG
sudo apt-get -y upgrade &>> $INSTALL_LOG
echo ">>> update packages - HOME=${HOME}" &>> $INSTALL_LOG
sudo apt-get -y update &>> $INSTALL_LOG

# install net-tools
echo ">>> install net-tools" &>> $INSTALL_LOG
sudo apt-get -y install net-tools jq &>> $INSTALL_LOG

# How to setup SCTP in ubuntu
echo ">>> install SCTP" &>> $INSTALL_LOG
sudo apt-get install -y libsctp-dev lksctp-tools &>> $INSTALL_LOG
sudo modprobe sctp &>> $INSTALL_LOG
sudo lsmod | grep sctp &>> $INSTALL_LOG

# for conntrack
echo ">>> enable conntrack" &>> $INSTALL_LOG
sudo modprobe nf_conntrack &>> $INSTALL_LOG
sudo lsmod | grep conntrack &>> $INSTALL_LOG

# end of script.
