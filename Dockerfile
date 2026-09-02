FROM ubuntu:24.04

RUN apt-get update && apt-get install -y golang make ca-certificates

WORKDIR /app

COPY . .

CMD ["make","run"]

