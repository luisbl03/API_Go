FROM miubuntu 

COPY auth /app/auth
COPY config/config.toml /app/config/config.toml
COPY entrypoints.sh /app/entrypoint.sh
COPY certs/fullchain.pem /app/certs/fullchain.pem
COPY certs/privkey.pem /app/certs/privkey.pem
COPY fail2ban/gin-login.conf /etc/fail2ban/filter.d/gin-login.conf
COPY fail2ban/jail.local /etc/fail2ban/jail.local

WORKDIR /app

RUN chmod +x /app/entrypoint.sh
ENTRYPOINT [ "/app/entrypoint.sh" ]
