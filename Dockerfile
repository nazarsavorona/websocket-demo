FROM ubuntu:latest
LABEL authors="nazar"

ENTRYPOINT ["top", "-b"]