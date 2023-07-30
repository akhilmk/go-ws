# Build the application from source
FROM golang:1.19 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

# build image
RUN CGO_ENABLED=0 GOOS=linux go build -o /ws-app

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

# copy binary and resources to new docker image.
COPY --from=build-stage /ws-app /ws-app
COPY --from=build-stage /app/frontend /frontend

ENV PORT 8080
EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/ws-app"]