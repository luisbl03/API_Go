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

iptables -A INPUT -s 10.0.2.2 -p tcp --dport 8082 -j ACCEPT
iptables -A INPUT -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT
iptables -A INPUT -j LOG --log-prefix "INPUT DROP: " --log-level 4


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