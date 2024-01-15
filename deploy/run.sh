#!/bin/bash
PRO_NAME=$1
nohup ./"$PRO_NAME" -mode=prod >"$PRO_NAME".log 2>&1 &
echo -e "\n"
echo "nohup success"
exit 0
