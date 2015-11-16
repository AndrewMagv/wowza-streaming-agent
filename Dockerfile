FROM scratch
MAINTAINER YI-HUNG JEN <yihungjen@gmail.com>

COPY ca-certificates.crt /etc/ssl/certs/
COPY wowza-streaming-agent /
ENTRYPOINT ["/wowza-streaming-agent"]
CMD ["--help"]

EXPOSE 20080
