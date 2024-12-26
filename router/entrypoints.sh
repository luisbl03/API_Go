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

#permitir trafico entrante por el puerto 8080
iptables -A INPUT -s -p tcp --dport 8080 -j ACCEPT
iptables -A INPUT -s -p tcp --dport 8081 -j ACCEPT
iptables -A INPUT -s -p tcp --dport 8082 -j ACCEPT

# redirigir el trafico generado en la red 172.17.0.0/16 en el puerto 8080 a la 10.0.1.4:8080
iptables -t nat -A PREROUTING -p tcp --dport 8080 -j DNAT --to-destination 10.0.1.4:8080
iptables -t nat -A POSTROUTING -o eth1 -p tcp --dport 8080 -s 172.17.0.0/16 -d 10.0.1.4 -j SNAT --to-source 10.0.1.2

#redirigir trafico generado por 10.0.1.4 en el puerto 8081 a 10.0.2.3:8081
iptables -t nat -A PREROUTING -p tcp --dport 8081 -j DNAT --to-destination 10.0.2.3:8081
iptables -t nat -A POSTROUTING -o eth2 -p tcp --dport 8081 -s 10.0.1.0/24 -d 10.0.2.3 -j SNAT --to-source 10.0.2.2

#redirigir trafico generado por 10.0.1.4 en el puerto 8082 a 10.0.2.4:8082
iptables -t nat -A PREROUTING -p tcp --dport 8082 -j DNAT --to-destination 10.0.2.4:8082
iptables -t nat -A POSTROUTING -o eth2 -p tcp --dport 8082 -s 10.0.1.0/24 -d 10.0.2.4 -j SNAT --to-source 10.0.2.2

#permitir trafico entre la red
iptables -A FORWARD -i eth0 -o eth1 -s 172.17.0.0/16 -d 10.0.1.4 -p tcp --dport 8080 -j ACCEPT
iptables -A FORWARD -i eth1 -o eth0 -s 10.0.1.4 -d 172.16.0.0/16 -p tcp --sport 8080 -j ACCEPT 
iptables -A FORWARD -i eth1 -o eth2 -s 10.0.1.4 -d 10.0.2.0/24 -p tcp --dport 8081 -j ACCEPT
iptables -A FORWARD -i eth2 -o eth1 -s 10.0.2.0/24 -d 10.0.1.4 -p tcp --sport 8081 -j ACCEPT
iptables -A FORWARD -i eth1 -o eth2 -s 10.0.1.4 -d 10.0.2.0/24 -p tcp --dport 8082 -j ACCEPT
iptables -A FORWARD -i eth2 -o eth1 -s 10.0.2.0/24 -d 10.0.1.4 -p tcp --sport 8082 -j ACCEPT

# Permitir tráfico establecido y relacionado
iptables -A FORWARD -m state --state ESTABLISHED,RELATED -j ACCEPT


# Iniciar servicios necesarios
service ssh start
service rsyslog start

# Mantener contenedor activo
exec "/bin/bash"
