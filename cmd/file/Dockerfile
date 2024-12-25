FROM miubuntu 

COPY file /app/file
COPY config/config.toml /app/config/config.toml
COPY entrypoints.sh /app/entrypoint.sh

WORKDIR /app

RUN chmod +x /app/entrypoint.sh
ENTRYPOINT [ "/app/entrypoint.sh" ]

CMD ["/app/file"]