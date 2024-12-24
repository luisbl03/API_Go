FROM miubuntu 

COPY auth /app/auth
COPY config/config.toml /app/config.toml
COPY entrypoints.sh /app/entrypoint.sh

WORKDIR /app

ENTRYPOINT [ "entrypoint.sh" ]

CMD ["/app/auth"]
