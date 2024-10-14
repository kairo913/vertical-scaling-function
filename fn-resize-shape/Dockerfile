FROM fnproject/go:1.23-dev as build-stage
WORKDIR /function
WORKDIR /go/src/func/
ENV GO111MODULE=on
COPY . .
RUN go mod tidy
RUN go build -o func -v
FROM fnproject/go:1.23
WORKDIR /function
COPY --from=build-stage /go/src/func/func /function/
ARG PRIVATE_KEY_NAME=privatekey.pem
COPY $PRIVATE_KEY_NAME /function/$PRIVATE_KEY_NAME
ENTRYPOINT [ "./func" ]