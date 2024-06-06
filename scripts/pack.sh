#!/bin/bash

ver=$1
ver_label=$(echo "$ver" | tr '.' '_')
ver_ip=$(echo "$ver" | tr '.' '-')

cd ../pkg || exit

echo -e "{
  \"address\": \"asset-cc-$ver_ip:6051\",
  \"dial_timeout\": \"10s\",
  \"tls_required\": false
}" > connection.json

echo -e "{
  \"type\": \"external\",
  \"label\": \"asset_$ver_label\"
}" > metadata.json

tar cfz code.tar.gz connection.json
tar cfz "asset_v$ver_label.tar.gz" metadata.json code.tar.gz
rm code.tar.gz