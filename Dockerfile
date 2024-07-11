FROM golang AS builder

COPY . /usr/project

WORKDIR /usr/project

RUN go build

FROM busybox AS export

COPY --from=builder /usr/project/vortex /usr/bin/vortex

WORKDIR /usr/bin

CMD [ "/usr/bin/vortex" ]