# Start from golang base image
FROM golang:alpine as development

# Add Maintainer info
LABEL maintainer="Oleg Komisarenko"

# Add a work directory inside container
WORKDIR /app

# Cache and install dependencies
#COPY go.mod go.sum .env ./
# Copy the source from the current directory to the working Directory inside the container
COPY . .
RUN go mod download

# Expose port 8080 to the outside world
EXPOSE 8080

#Command to run the executable - Start app
CMD go run main.go