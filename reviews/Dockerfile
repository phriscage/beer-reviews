###############################################################################
##  Name:   Dockerfile
##  Date:   2018-01-15
##  Developer:  Chris Page
##  Email:  phriscage@gmail.com
##  Purpose:   This Dockerfile contains the Beers Review xample
################################################################################
FROM golang:1.9.2-alpine3.7

# Set the file maintainer (your name - the file's author)
MAINTAINER Chris Page <phriscage@gmail.com>

#WORKDIR /go/src/github.com/phriscage/beer-reviews

# Install Git, Go dependencies, and build the app
RUN apk --no-cache add --virtual git

# Pull the Go dependencies
RUN go get github.com/pilu/fresh && \
        go get -d -v \
        github.com/olivere/elastic \
        github.com/gorilla/handlers \
        github.com/sirupsen/logrus \
        github.com/gorilla/mux

# Add the sample data
COPY data ./

# Add the application
COPY *.go ./

# Build the package
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

EXPOSE 8080

ENTRYPOINT []

CMD ["./app"]
# Build the package
