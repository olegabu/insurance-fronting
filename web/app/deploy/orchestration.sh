#!/usr/bin/env bash

ssh vp0 'bash -s' < reset.sh
ssh vp1 'bash -s' < reset.sh
ssh vp2 'bash -s' < reset.sh
ssh vp3 'bash -s' < reset.sh


cat start_vp0.sh | ssh vp0
cat start_vp1.sh | ssh vp1
cat start_vp2.sh | ssh vp2
cat start_vp3.sh | ssh vp3


#vp0
echo "peer network login citi -p 4nXSrfoYGFCP" | ssh vp0
echo "peer network login auditor -p yg5DVhm0er1z" | ssh vp0
echo "peer network login bermuda -p b7pmSxzKNFiw" | ssh vp0

#vp1
echo "peer network login allianz -p W8G0usrU7jRk" | ssh vp1

#vp2
echo "peer network login nigeria -p H80SiB5ODKKQ" | ssh vp2

#vp3
echo "peer network login art -p YsWZD4qQmYxo" | ssh vp3



