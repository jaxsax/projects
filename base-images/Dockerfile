FROM ubuntu:20.04

ENV DEBIAN_FRONTEND="noninteractive"

RUN apt-get -q update \
    && apt-get install -qy curl htop ca-certificates dnsutils locales \
    && locale-gen "en_US.UTF-8" \
    && apt-get autoremove -qy \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

ENV LANG="en_US.UTF-8" \
    LANGUAGE="en_US:en" \
    LC_ALL="en_US.UTF-8" \
    TZ="Etc/UTC"

# Add Tini
ENV TINI_VERSION v0.19.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini
RUN chmod +x /tini
ENTRYPOINT ["/tini", "--"]
