FROM golang:1.20.3-alpine

# set the working directory to /app
WORKDIR /app

# copy the go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# download and install any required GO dependencies
RUN go mod download

# Copy the entire source code to the working directory
COPY . .

# Build the GO app
RUN go build -o main .

# Expose port specified by the PORT environment variable
EXPOSE 8080

# Set the entry point of the container to the executable
CMD ["./main"]