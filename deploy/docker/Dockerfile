FROM golang:1.24 AS build
WORKDIR /src/
COPY src/ .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o checklist

FROM scratch
COPY --from=build /src/checklist .
CMD [ "./checklist" ]