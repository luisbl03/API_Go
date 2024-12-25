#!/bin/bash
echo 1 > /proc/sys/net/ipv4/ip_forward

iptables -P INPUT DROP
iptables -P FORWARD DROP
iptables -P OUTPUT ACCEPT

iptables -A INPUT -i lo -j ACCEPT
iptables -A INPUT -p icmp -j ACCEPT
iptables -A FORWARD -p icmp -j ACCEPT

iptables -A INPUT -s 10.0.1.2 -p tcp --dport 8080 -j ACCEPT

iptables -A INPUT -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT


service ssh start
service rsyslog start

exec /bin/bash
