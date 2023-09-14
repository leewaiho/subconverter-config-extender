FROM debian:buster-slim

COPY subconverter-config-extender /usr/local/bin/subconverter-config-extender

RUN apt-get update -y && \
    apt-get install -y ca-certificates && \
    chmod a+x /usr/local/bin/subconverter-config-extender

CMD ["/usr/local/bin/subconverter-config-extender"]