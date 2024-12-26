FROM miubuntu 

COPY broker /app/broker
COPY config/config.toml /app/config/config.toml
COPY entrypoints.sh /app/entrypoint.sh

WORKDIR /app

ENTRYPOINT [ "/app/entrypoint.sh" ]