FROM docker.io/golang:latest as build
COPY . /
WORKDIR /

RUN CGO_ENABLED=0 GOOS=linux go build -a -o messanger

## Stage 2
FROM scratch
COPY --from=build /messanger .
CMD [ "/messanger" ]