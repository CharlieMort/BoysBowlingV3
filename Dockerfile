FROM golang:1.19

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY backend/*.go ./
COPY frontend/ ../frontend

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /boysbowling

EXPOSE 8888

# Run
CMD ["/boysbowling"]