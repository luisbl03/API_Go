#!/bin/bash

# Habilitar reenvío de paquetes
echo 1 > /proc/sys/net/ipv4/ip_forward


# Políticas por defecto (más seguras)
iptables -P INPUT DROP
iptables -P FORWARD DROP
iptables -P OUTPUT ACCEPT

# Permitir tráfico local
iptables -A INPUT -i lo -j ACCEPT

# Permitir tráfico ICMP
iptables -A INPUT -p icmp -j ACCEPT
iptables -A FORWARD -p icmp -j ACCEPT
iptables -t nat -A POSTROUTING -o eth0 -p icmp -j MASQUERADE

#ssh
iptables -A FORWARD -i eth0 -o eth1 -p tcp --syn --dport 22 -m state --state NEW -j ACCEPT
iptables -A FORWARD -i eth0 -o eth1 -m state --state ESTABLISHED,RELATED -j ACCEPT
iptables -A FORWARD -i eth1 -o eth0 -m state --state ESTABLISHED,RELATED -j ACCEPT


iptables -t nat -A PREROUTING -i eth0 -p tcp --dport 22 -j DNAT --to-destination 10.0.1.3
iptables -t nat -A POSTROUTING -o eth1 -p tcp --dport 22 -s 172.17.0.0/16 -d 10.0.1.3 -j SNAT --to-source 10.0.1.2
iptables -t nat -A POSTROUTING -o eth3 -p tcp --dport 22 -s 172.17.0.0/16 -d 10.0.3.3 -j SNAT --to-source 10.0.3.2

iptables -A FORWARD -i eth1 -o eth2 -p tcp --dport 22 -j ACCEPT
iptables -A FORWARD -i eth2 -o eth1 -p tcp --sport 22 -j ACCEPT
iptables -A FORWARD -i eth2 -o eth1 -p tcp --dport 22 -j ACCEPT
iptables -A FORWARD -i eth1 -o eth2 -p tcp --sport 22 -j ACCEPT
iptables -A FORWARD -i eth3 -o eth2 -p tcp --dport 22 -j ACCEPT
iptables -A FORWARD -i eth2 -o eth3 -p tcp --sport 22 -j ACCEPT
iptables -A FORWARD -i eth2 -o eth3 -p tcp --dport 22 -j ACCEPT
iptables -A FORWARD -i eth3 -o eth2 -p tcp --sport 22 -j ACCEPT
iptables -A FORWARD -i eth1 -o eth3 -p tcp --dport 22 -j ACCEPT
iptables -A FORWARD -i eth3 -o eth1 -p tcp --sport 22 -j ACCEPT
iptables -A FORWARD -i eth3 -o eth1 -p tcp --dport 22 -j ACCEPT
iptables -A FORWARD -i eth1 -o eth3 -p tcp --sport 22 -j ACCEPT

iptables -A INPUT -p tcp --dport 22 -i eth2 -s 10.0.3.3 -j ACCEPT
iptables -A INPUT -p tcp --dport 22 -s 10.0.3.3 -j ACCEPT
iptables -A INPUT -p tcp --dport 22 -i eth3 -s 10.0.3.3 -j ACCEPT

#http

iptables -A FORWARD -p tcp --dport 5000 -j ACCEPT
iptables -A FORWARD -p tcp --sport 5000 -j ACCEPT

iptables -t nat -A PREROUTING -i eth0 -p tcp --dport 5000 -j DNAT --to-destination 10.0.1.4
iptables -t nat -A POSTROUTING -o eth1 -p tcp --dport 5000 -s 172.17.0.0/16 -d 10.0.1.4 -j SNAT --to-source 10.0.1.2

#http
iptables -A INPUT -p tcp --dport 80 -m state --state NEW,ESTABLISHED,RELATED -j ACCEPT
iptables -A INPUT -p tcp --sport 80 -m state --state ESTABLISHED,RELATED -j ACCEPT

#https
iptables -A INPUT -p tcp --dport 443 -m state --state NEW,ESTABLISHED,RELATED -j ACCEPT
iptables -A INPUT -p tcp --sport 443 -m state --state ESTABLISHED,RELATED -j ACCEPT


# Iniciar servicios necesarios
service ssh start
service rsyslog start

# Mantener contenedor activo
exec "/bin/bash"
