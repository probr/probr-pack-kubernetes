FROM golang:1.14.4-alpine
WORKDIR /probr
COPY . .
RUN go build -o /out/probr .

CMD /out/probr