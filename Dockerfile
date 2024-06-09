# Dockerfile

# Etapa de construcción
FROM golang:1.16 as builder

WORKDIR /workspace

# Copiar los archivos del proyecto
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY main.go main.go
COPY api/ api/
COPY controllers/ controllers/

# Construir el operador
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o email-operator main.go

# Etapa final
FROM alpine:3.14

WORKDIR /
COPY --from=builder /workspace/email-operator .
ENTRYPOINT ["/email-operator"]

