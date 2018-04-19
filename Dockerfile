# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

RUN go get -u github.com/go-redis/redis
RUN go get -u github.com/gorilla/mux
RUN go get -u gopkg.in/mgo.v2
RUN go get -u gopkg.in/mgo.v2/bson

RUN mkdir -p /go/src/ml/com/mutants/api-rest/

WORKDIR /go/src/ml/com/mutants/api-rest

# Copy the local package files to the container's workspace.
ADD /api-rest/* /go/src/ml/com/mutants/api-rest/

# Build the app command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go install ml/com/mutants/api-rest

# Run the api-rest command by default when the container starts.
ENTRYPOINT /go/bin/api-rest




