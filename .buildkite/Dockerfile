FROM --platform=linux/amd64 alpine

ARG buildkite_version=3.41.0

LABEL name="catalog-build" \
    maintainers="prod-ops@greymatter.io" \
    base="alpine"

# Install Required Build Packages
RUN apk update
RUN apk add curl wget git bash go npm openssh

# Install Buildah
RUN apk add buildah cni-plugins iptables ip6tables

# Install Buildkite
RUN mkdir /buildkite && cd /buildkite
RUN wget https://github.com/buildkite/agent/releases/download/v${buildkite_version}/buildkite-agent-linux-amd64-${buildkite_version}.tar.gz
RUN tar -xvf buildkite-agent-linux-amd64-${buildkite_version}.tar.gz -C /buildkite && rm buildkite-agent-linux-amd64-${buildkite_version}.tar.gz 
ENV PATH=$PATH:/buildkite

# Add Build Scripts
COPY rex-init.sh /scripts/rex-init.sh
COPY rex-build-backend.sh /scripts/rex-build-backend.sh
COPY rex-push-backend.sh /scripts/rex-push-backend.sh
RUN chmod +x /scripts/rex-init.sh /scripts/rex-build-backend.sh \
    /scripts/rex-push-backend.sh