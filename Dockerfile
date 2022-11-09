FROM golang:1.16
WORKDIR /opt/
COPY src/. ./
# CMD ["/bin/bash"]
RUN go mod download && CGO_ENABLED=0 go build -o http-server

FROM alpine:latest  
WORKDIR /root/
EXPOSE 9090/tcp
COPY --from=0 /opt/http-server ./
CMD ["./http-server"]