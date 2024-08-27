# start from base image
FROM golang:alpine

# Add Maintainer info
LABEL maintainer = "George Smith"

# install git
RUN apk update && apk add --no-cache git && apk add --no-cach bash && apk add build-base

# setup folders
RUN mkdir /aggregator
WORKDIR /aggregator

# Copy source from current directory
COPY . .
COPY .env .

# Download dependencies
RUN go get -d -v ./...

RUN go install -v ./...

RUN go build -o /build

#expose port 
EXPOSE 8080

CMD ["/build"]


