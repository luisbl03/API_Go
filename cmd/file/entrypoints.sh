#!/bin/bash

echo 1 > /proc/sys/net/ipv4/ip_forward

ufw allow 8080/tcp
ufw allow 8082/tcp

iptables -P INPUT DROP
iptables -P FORWARD DROP
iptables -P OUTPUT ACCEPT

iptables -A INPUT -i lo -j ACCEPT
iptables -A INPUT -p icmp -j ACCEPT
iptables -A FORWARD -p icmp -j ACCEPT

iptables -A INPUT -p tcp --dport 5000 -i eth0 -s 10.0.1.4 -j ACCEPT
iptables -A INPUT -p tcp --dport 80 -m state --state NEW,ESTABLISHED,RELATED -j ACCEPT
iptables -A INPUT -p tcp --sport 80 -m state --state ESTABLISHED,RELATED -j ACCEPT

service ssh start
service rsyslog start

ip route del default
ip route add default via 10.0.2.2 dev eth0

./file &

if [ -z "$@" ]; then
    exec /bin/bash
else
    exec $@
fi