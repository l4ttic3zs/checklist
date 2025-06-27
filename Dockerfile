FROM docker.io/library/golang:latest AS build
WORKDIR /
COPY src/ .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o checklist

FROM scratch 
COPY --from=build /checklist .
CMD [ "./checklist" ]