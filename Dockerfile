FROM golang:alpine as build

COPY src/ /src

WORKDIR /src

RUN mkdir -p /app/dist && go build -i -o /app/dist/ado_server

FROM golang:alpine

COPY --from=build /app/dist/ado_server /app/ado_server

EXPOSE 8080

ENTRYPOINT ["/app/ado_server"]