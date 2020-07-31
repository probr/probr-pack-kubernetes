FROM golang:1.14.4-alpine AS probr-build

WORKDIR /probr

COPY . .

RUN go build -o /out/probr .


FROM node:latest  
WORKDIR /probr
COPY --from=probr-build /out/probr .
COPY view .
COPY run.sh .

CMD ["./run.sh"]
