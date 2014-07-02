#!/usr/bin/env bash

set -e

CONFIG_DIR=$1

if [[ -z $CONFIG_DIR ]]
then
  echo "Usage: setup_syslog_forwarder [config dir]"
  exit 1
fi

# Place to spool logs if the upstream server is down
mkdir -p /var/vcap/sys/rsyslog/buffered
chown -R syslog:adm /var/vcap/sys/rsyslog/buffered

cp $CONFIG_DIR/syslog_forwarder.conf /etc/rsyslog.d/00-syslog_forwarder.conf

/usr/sbin/service rsyslog restart
