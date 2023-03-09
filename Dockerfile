FROM alpine:3.17

LABEL maintainer="Dmitry Mozzherin"

WORKDIR /bin

COPY ./gnverifier /bin

ENTRYPOINT [ "gnverifier" ]

CMD ["-p", "8181"]
