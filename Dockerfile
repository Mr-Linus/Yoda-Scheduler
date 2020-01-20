FROM debian:stretch-slim

WORKDIR /

COPY yoda-scheduler /usr/local/bin

CMD ["yoda-scheduler"]