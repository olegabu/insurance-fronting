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
cat login_vp0.sh | ssh vp0

#vp1
cat login_vp1.sh | ssh vp1

#vp2
cat login_vp2.sh | ssh vp2

#vp3
cat login_vp3.sh | ssh vp3





