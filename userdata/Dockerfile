FROM docker:17.07.0-ce-dind
RUN apk update && \
    apk add --no-cache \
        curl

COPY userdata.sh .
CMD ["./userdata.sh"]
