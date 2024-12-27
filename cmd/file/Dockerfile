FROM miubuntu 

COPY file /app/file
COPY config/config.toml /app/config/config.toml
COPY entrypoints.sh /app/entrypoint.sh
COPY certs/fullchain.pem /app/certs/fullchain.pem
COPY certs/privkey.pem /app/certs/privkey.pem
WORKDIR /app

RUN chmod +x /app/entrypoint.sh
ENTRYPOINT [ "/app/entrypoint.sh" ]
