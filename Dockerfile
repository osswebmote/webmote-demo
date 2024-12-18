FROM node:lts AS build
WORKDIR /build

COPY . .

RUN yarn && yarn astro build

FROM golang:alpine AS build-svr
WORKDIR /build

COPY . .

RUN go build  -o app ./

FROM gcr.io/distroless/static-debian12 AS runtime
WORKDIR /app

COPY --from=build /build/dist /app/dist
COPY --from=build-svr /build/app /app/app
COPY --from=build-svr /build/tls.* /app/

CMD ["/app/app"]
EXPOSE 8000