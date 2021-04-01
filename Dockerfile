FROM alpine

LABEL maintainer="Dmitry Mozzherin"

ENV LAST_FULL_REBUILD 2021-03-12

WORKDIR /bin

COPY ./gnverifier/gnverifier /bin

ENTRYPOINT [ "gnverifier" ]

CMD ["-w", "8181"]
