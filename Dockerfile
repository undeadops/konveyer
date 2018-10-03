# build stage
FROM golang:1.11-stretch AS build-env
ADD . /go/konveyer
RUN cd /go/konveyer/cmd/konveyer-server && go install ./...

# final stage
FROM debian:stretch-slim
COPY --from=build-env /go/bin/* /usr/bin/
CMD ["konveyer-server"]
