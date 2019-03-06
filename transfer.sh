#!/bin/bash

./block send  张三 李四 10 班长 "张三转李四10"
./block send  张三 王五 20 班长 "张三转王五20"
./checkBalance.sh

echo "======================"
./block send  王五 李四 2 班长 "王五转李四2"
./block send  王五 李四 3 班长 "王五转李四3"
./block send  王五 张三 5 班长 "王五转张三5"
./checkBalance.sh

echo "======================"
./block send  李四 赵六 14 班长 "李四转赵六14"
./checkBalance.sh

