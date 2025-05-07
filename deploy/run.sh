#!/bin/bash
PRO_NAME=$1
if [ -f "$PRO_NAME.pid" ]; then
    pid=$(cat "$PRO_NAME".pid)
    if [ ! -z "$pid" ]; then
        echo "结束进程 $pid"
        kill "$pid"
    fi
fi
nohup ./"$PRO_NAME" -mode=prod >"$PRO_NAME".log 2>&1 &
echo -e "\n"
echo "nohup success"
exit 0
