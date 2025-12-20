FROM docker.io/golang:latest as build
COPY . /
WORKDIR /

RUN CGO_ENABLED=0 GOOS=linux go build -a -o catalog-init

## Stage 2
FROM scratch
COPY --from=build /catalog-init /catalog-init
CMD [ "/catalog-init" ]