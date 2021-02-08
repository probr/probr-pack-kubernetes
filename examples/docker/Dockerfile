FROM golang:1.14.4-alpine AS probr-build
WORKDIR /probr
COPY . .
RUN go get github.com/markbates/pkger/cmd/pkger
RUN pkger
RUN go build -o probr cmd/main.go

FROM node:alpine
COPY internal/view /probr/internal/view
WORKDIR /probr/internal/view
RUN npm ci

WORKDIR /probr
COPY --from=probr-build /probr .
COPY /examples/docker/. .
RUN chmod +x *.sh

ENTRYPOINT ["/probr/entrypoint.sh"]
