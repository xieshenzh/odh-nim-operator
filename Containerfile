# Copyright (c) 2023 Red Hat, Inc.

FROM registry.access.redhat.com/ubi9/ubi-minimal:9.3

USER root
# build requires vpn and is performed outside the scope this image build
COPY build/odhnimoperator /usr/bin/odhnimoperator
COPY LICENSE /licenses/odhnimoperator-license
ENTRYPOINT ["/usr/bin/odhnimoperator"]
USER 1001
