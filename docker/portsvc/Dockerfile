ARG IMAGE_TAG

FROM ${IMAGE_TAG}

ENV PORT ${PORT:-8081}
ENV LOGLEVEL ${LOGLEVEL:-info}

# copy config to a container
ADD ./config/config.ini /app

WORKDIR /app
