FROM alpine

LABEL maintainer="Dmitry Mozzherin"

ENV LAST_FULL_REBUILD 2021-03-12

WORKDIR /bin

COPY ./gnverify/gnverify /bin

ENTRYPOINT [ "gnverify" ]

CMD ["-w", "8181"]
