# syntax = docker/dockerfile:experimental
FROM ubuntu:22.04

ARG DEBIAN_FRONTEND=noninteractive
RUN apt update && apt upgrade -y && apt install ca-certificates python3 hostname libunwind8 netbase \
        ethtool iproute2 ncat libunbound-dev procps libatomic1 kmod iptables python3-netifaces python3-sortedcontainers \
        tcpdump ipset curl uuid-runtime openssl inetutils-ping arping ndisc6 \
        logrotate dnsutils net-tools strongswan strongswan-pki libcharon-extra-plugins \
        libcharon-extauth-plugins libstrongswan-extra-plugins libstrongswan-standard-plugins -y --no-install-recommends && \
        rm -rf /var/lib/apt/lists/* && \
        rm -rf /etc/localtime

WORKDIR /network-pinger

COPY network-pinger /network-pinger/network-pinger