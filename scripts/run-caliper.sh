#!/bin/sh

chan_name=$1
org_msp=$2
pvt_key_name=$3
pub_cert_name=$4
tls_cert_name=$5
usr=$6
peer_host=$7
peer_port=$8

pvt_key_path="/hyperledger/caliper/workspace/peer/$pvt_key_name"
pub_cert_path="/hyperledger/caliper/workspace/peer/$pub_cert_name"
tls_cert_path="/hyperledger/caliper/workspace/peer/$tls_cert_name"
peer_endpoint="$peer_host:$peer_port"

sed -i "s+'<chan-name>'+'$chan_name'+g" caliper/network.yaml
sed -i "s+'<org-msp>'+'$org_msp'+g" caliper/network.yaml
sed -i "s+'<msp-private-key-file-path>'+'$pvt_key_path'+g" caliper/network.yaml
sed -i "s+'<msp-public-cert-path>'+'$pub_cert_path'+g" caliper/network.yaml
sed -i "s+'<tls-root-cert-path>'+'$tls_cert_path'+g" caliper/network.yaml
sed -i "s+'<user-name>'+'$usr'+g" caliper/network.yaml
sed -i "s+'<peer-endpoint>'+'$peer_endpoint'+g" caliper/network.yaml