#!/bin/bash
PRO_NAME=$1
pid=$(cd .. && cat $PRO_NAME.pid)
echo "结束进程 $pid"
kill "$pid"

nohup ./"$PRO_NAME" -mode=prod >"$PRO_NAME".log 2>&1 &
echo -e "\n"
echo "nohup success"
exit 0
