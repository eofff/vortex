FROM golang AS builder

COPY ./project /usr/project

WORKDIR /usr/project

RUN go build

FROM busybox AS export

COPY --from=builder /usr/project/vortex /usr/bin/vortex
COPY project/migrations /usr/bin/migrations

WORKDIR /usr/bin

CMD [ "/usr/bin/vortex" ]