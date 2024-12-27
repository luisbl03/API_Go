FROM miubuntu 

RUN useradd -ms /bin/bash dev && \
    mkdir /home/dev/.ssh && \
    chmod 700 /home/dev/.ssh
COPY authorized_keys /home/dev/.ssh/authorized_keys
RUN chown dev:dev -R /home/dev
COPY op_key /home/op/.ssh/op_key
COPY op_key.pub /home/op/.ssh/op_key.pub
RUN chmod 600 /home/op/.ssh/op_key && \
    chmod 644 /home/op/.ssh/op_key.pub

COPY entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/bin/bash", "/entrypoint.sh"]