FROM alpine:3.21

LABEL maintainer="Dmitry Mozzherin"

RUN adduser -D -H gnverifier

WORKDIR /bin

COPY ./bin/gnverifier /bin/gnverifier

USER gnverifier

ENTRYPOINT [ "gnverifier" ]

CMD ["-p", "8181"]
