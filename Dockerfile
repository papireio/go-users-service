FROM golang:1.21

WORKDIR /app
COPY . .
RUN make vendor && make
WORKDIR ./target/

CMD ["./go-users"]
