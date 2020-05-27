FROM concourse/golang-builder AS builder
COPY . /go/src/github.com/concourse/datadog-event-resource
ENV CGO_ENABLED 0
RUN go build -o /assets/out github.com/concourse/datadog-event-resource/cmd/out
RUN go build -o /assets/in github.com/concourse/datadog-event-resource/cmd/in
RUN go build -o /assets/check github.com/concourse/datadog-event-resource/cmd/check
RUN set -e; for pkg in $(go list ./...); do \
		go test -o "/tests/$(basename $pkg).test" -c $pkg; \
	done

FROM ubuntu:bionic AS resource
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
  && rm -rf /var/lib/apt/lists/*
COPY --from=builder /assets /opt/resource

FROM resource AS tests
COPY --from=builder /tests /tests
RUN set -e; for test in /tests/*.test; do \
		$test; \
	done

FROM resource