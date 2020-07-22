FROM golang:1.14.4-alpine AS build
WORKDIR /probr
COPY . .
RUN go build -o /out/probr .

CMD /out/probr