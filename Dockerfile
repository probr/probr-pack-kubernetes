FROM golang:1.14.4-alpine AS probr-build
WORKDIR /probr
COPY . .
RUN go build -o /out/probr .

FROM node:alpine
RUN mkdir -p /probr/testoutput
COPY test /probr/test
COPY view /probr/view

WORKDIR /probr/view
RUN npm ci

WORKDIR /probr
COPY --from=probr-build /out/probr .
COPY run.sh .
ENTRYPOINT ["./run.sh"]
