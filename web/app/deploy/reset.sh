#!/usr/bin/env bash

echo "Stopping peer and membersrvc..."
ps -ef | grep membersrvc | grep -v grep | awk '{print $2}' | xargs kill &> /tmp/kill
ps -ef | grep peer | grep -v grep | awk '{print $2}' | xargs kill &> /tmp/kill
echo "   - Stopped"

echo "Removing previous configuration and state"
#reset configuration
rm -rf /var/hyperledger/production
echo "   - Removed"