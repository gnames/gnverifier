FROM alpine:3.14

LABEL maintainer="Dmitry Mozzherin"

ENV LAST_FULL_REBUILD 2021-04-09

WORKDIR /bin

COPY ./gnverifier/gnverifier /bin

ENTRYPOINT [ "gnverifier" ]

CMD ["-p", "8181"]
