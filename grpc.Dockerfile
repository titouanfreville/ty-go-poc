FROM scratch
ARG BUILD_NAME="go-poc-rcp"
ARG BUILD_DATE="N/A"
ARG BUILD_VCSREF="N/A"

LABEL maintainer="tfreville@tankyou.co" \
    org.label-schema.name="${BUILD_NAME}" \
    org.label-schema.description="ty-poc-go" \
    org.label-schema.build-date="${BUILD_DATE}" \
    org.label-schema.vcs-ref="${BUILD_VCSREF}"

ENV BUILD_VCSREF=${BUILD_VCSREF} \
    GIN_MODE=release \
    CP_ROOT=/go_poc

WORKDIR ${CP_ROOT}/

EXPOSE 3001
COPY go_poc_rcp go_poc_rcp
ENTRYPOINT ["./go_poc_rcp"]
