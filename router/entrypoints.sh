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

iptables -A FORWARD -p tcp --dport 5000 -j ACCEPT
iptables -A FORWARD -p tcp --sport 5000 -j ACCEPT

iptables -t nat -A PREROUTING -i eth0 -p tcp --dport 5000 -j DNAT --to-destination 10.0.1.4
iptables -t nat -A POSTROUTING -o eth1 -p tcp --dport 5000 -s 172.17.0.0/16 -d 10.0.1.4 -j SNAT --to-source 10.0.1.2

iptables -A INPUT -p tcp --dport 80 -m state --state NEW,ESTABLISHED,RELATED -j ACCEPT
iptables -A INPUT -p tcp --sport 80 -m state --state ESTABLISHED,RELATED -j ACCEPT


# Iniciar servicios necesarios
service ssh start
service rsyslog start

# Mantener contenedor activo
exec "/bin/bash"
