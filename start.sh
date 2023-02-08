#!/bin/bash
project="tx"
rm -f $project
rm -f out.log
rm -rf ./runtime
go build -o $project
chmod 0777 $project

nohup /home/www/$project/comics-crawler/$project > out.log 2>&1 &