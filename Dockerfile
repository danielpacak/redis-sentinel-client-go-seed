FROM scratch

COPY seed /seed

ENTRYPOINT ["/seed"]
