FROM miubuntu 

COPY broker /app/broker
COPY config/config.toml /app/config.toml
COPY entrypoints.sh /app/entrypoint.sh

WORKDIR /app

ENTRYPOINT [ "entrypoint.sh" ]

CMD ["/app/broker"]