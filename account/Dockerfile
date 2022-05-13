# first stage: builder. download and install all dependencies and then build this application binary.
# We'll also use this builder stage in our docker-compose development environment later for auto-reloading our application.
# We'll do this by means of the application called reflex, which we download in the builder stage.
# load imgae for container
FROM golang:alpine as builder

# directory for container (every in container related to this directory)
WORKDIR /go/src/app

# set env variable to let Go know it use module
ENV GO111MODULE=on

# bash command Run to run inside container
# Get Reflex for live reloading for the app - (re-run server when files change)
RUN go install github.com/cespare/reflex@latest

# copy go.mod and go.sum to container
COPY go.mod .
COPY go.sum .

# download all packages
RUN go mod download

# get all the download bring in source files
COPY . .

# build go executable to output file name 'run' at root directory
RUN go build -o ./run .


# second stage: lean build. extract the built application from builder, and run it in a separate container.
#This eliminates all of the unnecessary code from the builder stage and creates a leaner final application container which is ready for deployment.
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

#Copy the binary files made before to the new root folder
COPY --from=builder /go/src/app/run .

# expose to gate and run
EXPOSE 8080
CMD ["./run"]