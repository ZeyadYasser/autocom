FROM golang:1.11

ENV GOBIN="$GOPATH/bin"
WORKDIR /src

# Fetch dependencies first; they are less susceptible to change on every build
# and will therefore be cached for speeding up the next build
COPY ./go.mod ./go.sum ./
RUN go mod download

# Import the code.
COPY . .

RUN go install -v .

CMD [ "autocom" ]