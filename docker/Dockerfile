# builder
FROM public.ecr.aws/docker/library/golang:1.25.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api ./cmd/api

# runner
FROM public.ecr.aws/docker/library/golang:1.25.0-alpine AS runner
RUN addgroup -S appgroup && adduser -S nonroot -G appgroup

WORKDIR /home/nonroot
COPY --from=builder /app/api .
RUN chown -R nonroot:appgroup /home/nonroot

USER nonroot
CMD ["./api"]