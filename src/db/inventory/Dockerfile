FROM docker.io/golang:latest as build
COPY . /
WORKDIR /

RUN CGO_ENABLED=0 GOOS=linux go build -a -o inventory

## Stage 2
FROM scratch
COPY --from=build /inventory .
CMD [ "/inventory" ]