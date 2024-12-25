FROM miubuntu 

COPY auth /app/auth
COPY config/config.toml /app/config/config.toml
COPY entrypoints.sh /app/entrypoint.sh

WORKDIR /app

RUN chmod +x /app/entrypoint.sh
ENTRYPOINT [ "/app/entrypoint.sh" ]

CMD ["/app/auth"]
