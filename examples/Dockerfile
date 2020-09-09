FROM golang:1.14.4-alpine AS probr-build
WORKDIR /probr
COPY . .
RUN go build -o /out/probr cmd/probr-cli/main.go

FROM node:alpine
RUN mkdir -p /probr/testoutput
COPY test /probr/test
COPY internal/view /probr/view

WORKDIR /probr/view
RUN npm ci

WORKDIR /probr
COPY --from=probr-build /out/probr .
COPY entrypoint.sh .
ENTRYPOINT ["./entrypoint.sh"]
