FROM scratch
ADD cb .
ENTRYPOINT ["/cb"]
