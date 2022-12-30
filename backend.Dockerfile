FROM golang 

WORKDIR /build
COPY . . 
RUN go build -v -o /app/rex 
WORKDIR /app
CMD ["./rex"]



