###############################################################################
##  Name:   Dockerfile
##  Date:   2018-01-15
##  Developer:  Chris Page
##  Email:  phriscage@gmail.com
##  Purpose:   This Dockerfile contains the Beers Review xample
################################################################################
## build stage
FROM golang:1.9.2-alpine3.7 AS build-env

# Set the file maintainer (your name - the file's author)
MAINTAINER Chris Page <phriscage@gmail.com>

# app working directory
WORKDIR /app

# Install Git, Go dependencies, and build the app
RUN apk --no-cache add --virtual git && \
        rm -rf /var/cache/apk/*

# Add the sample data
COPY data /app/data

# Add the application
COPY *.go ./

# Pull the Go dependencies
RUN go get -d -v 

# Build the package
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/main .

EXPOSE 8080

ENTRYPOINT ["/app/main"]


## final stage
FROM alpine:3.7
WORKDIR /app
COPY --from=build-env /app /app
ENTRYPOINT ["/app/main"]
