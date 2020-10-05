FROM golang:1.14.4-alpine AS probr-build
WORKDIR /probr
COPY . .
RUN go build -o /out/probr cmd/probr-cli/*.go

FROM node:alpine
COPY integration/clouddriver/kubernetes/assets /probr/integration/clouddriver/kubernetes/assets
COPY probes probr/probes
COPY internal/view /probr/internal/view

WORKDIR /probr/internal/view
RUN npm ci

WORKDIR /probr
COPY --from=probr-build /out/probr .
COPY entrypoint.sh .
ENTRYPOINT ["./entrypoint.sh"]
