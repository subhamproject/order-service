FROM golang:1.19.4-buster


#copy binary
COPY bin/order-service order-service

EXPOSE 9092
# Run the binary program produced by `go install`
CMD ["./order-service"]


