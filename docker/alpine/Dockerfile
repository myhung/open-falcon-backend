FROM golang:1.7-alpine

LABEL maintainer cheminlin@cepave.com

ARG ENTRYFILE=run.sh
ENV CI_PROJECT_NAME=open-falcon-backend
ENV FALCON_DIR=/home DOCKER_DIR=docker CONFIG_DIR=/config CONFIG_LINK=config/cfg.json
ENV WORK_PATH=${GOPATH}/src/github.com/Cepave/${CI_PROJECT_NAME}
ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/go/bin:$GOPATH/bin

# Build binary
COPY . ${WORK_PATH}
WORKDIR ${WORK_PATH}
RUN \
  apk add --no-cache ca-certificates bash curl git iproute2 tzdata jq \
  && mkdir -p ${FALCON_DIR} ${CONFIG_DIR} \
  && apk add --no-cache --virtual .build-deps g++ perl make \
  && make all \
  && make pack \
  && tar zxvf open-falcon-v*.tar.gz -C ${FALCON_DIR} \
  && rm -rf ${GOPATH} \
  && apk del .build-deps

# Set timezone
# Set alias in the case of user want to execute control in their terminal
# Set soft links
RUN \
  apk add --no-cache tzdata \
  && cp /usr/share/zoneinfo/Asia/Taipei /etc/localtime \
  && echo "Asia/Taipei" > /etc/timezone \
  && echo "alias ps='pstree'" > ~/.bashrc \
  && touch ${CONFIG_DIR}/agent.json \
  && ln -sf ${CONFIG_DIR}/agent.json ${FALCON_DIR}/agent/${CONFIG_LINK} \
  && rm -f ${CONFIG_DIR}/agent.json \
  && touch ${CONFIG_DIR}/aggregator.json \
  && ln -sf ${CONFIG_DIR}/aggregator.json ${FALCON_DIR}/aggregator/${CONFIG_LINK} \
  && rm -f ${CONFIG_DIR}/aggregator.json \
  && touch ${CONFIG_DIR}/alarm.json \
  && ln -sf ${CONFIG_DIR}/alarm.json ${FALCON_DIR}/alarm/${CONFIG_LINK} \
  && rm -f ${CONFIG_DIR}/alarm.json \
  && touch ${CONFIG_DIR}/fe.json \
  && ln -sf ${CONFIG_DIR}/fe.json ${FALCON_DIR}/fe/${CONFIG_LINK} \
  && rm -f ${CONFIG_DIR}/fe.json \
  && touch ${CONFIG_DIR}/graph.json \
  && ln -sf ${CONFIG_DIR}/graph.json ${FALCON_DIR}/graph/${CONFIG_LINK} \
  && rm -f ${CONFIG_DIR}/graph.json \
  && touch ${CONFIG_DIR}/hbs.json \
  && ln -sf ${CONFIG_DIR}/hbs.json ${FALCON_DIR}/hbs/${CONFIG_LINK} \
  && rm -f ${CONFIG_DIR}/hbs.json \
  && touch ${CONFIG_DIR}/judge.json \
  && ln -sf ${CONFIG_DIR}/judge.json ${FALCON_DIR}/judge/${CONFIG_LINK} \
  && rm -f ${CONFIG_DIR}/judge.json \
  && touch ${CONFIG_DIR}/nodata.json \
  && ln -sf ${CONFIG_DIR}/nodata.json ${FALCON_DIR}/nodata/${CONFIG_LINK} \
  && rm -f ${CONFIG_DIR}/nodata.json \
  && touch ${CONFIG_DIR}/query.json \
  && ln -sf ${CONFIG_DIR}/query.json ${FALCON_DIR}/query/${CONFIG_LINK} \
  && rm -f ${CONFIG_DIR}/query.json \
  && touch ${CONFIG_DIR}/sender.json \
  && ln -sf ${CONFIG_DIR}/sender.json ${FALCON_DIR}/sender/${CONFIG_LINK} \
  && rm -f ${CONFIG_DIR}/sender.json \
  && touch ${CONFIG_DIR}/task.json \
  && ln -sf ${CONFIG_DIR}/task.json ${FALCON_DIR}/task/${CONFIG_LINK} \
  && rm -f ${CONFIG_DIR}/task.json \
  && touch ${CONFIG_DIR}/transfer.json \
  && ln -sf ${CONFIG_DIR}/transfer.json ${FALCON_DIR}/transfer/${CONFIG_LINK} \
  && rm -f ${CONFIG_DIR}/transfer.json \
  && touch ${CONFIG_DIR}/nqm-mng.json \
  && ln -sf ${CONFIG_DIR}/nqm-mng.json ${FALCON_DIR}/nqm-mng/${CONFIG_LINK} \
  && rm -f ${CONFIG_DIR}/nqm-mng.json \
  && touch ${CONFIG_DIR}/f2e-api.json \
  && ln -sf ${CONFIG_DIR}/f2e-api.json ${FALCON_DIR}/f2e-api/${CONFIG_LINK} \
  && rm -f ${CONFIG_DIR}/f2e-api.json

COPY ${DOCKER_DIR}/alpine/${ENTRYFILE} ${FALCON_DIR}/run.sh
COPY ${DOCKER_DIR}/docker-healthcheck /usr/local/bin/

# Port
# Rpc:  10070
# Http: 10080 10081
EXPOSE 10070 10080 10081
WORKDIR ${FALCON_DIR}

# Start
ENTRYPOINT ["/bin/bash", "run.sh"]
HEALTHCHECK --interval=60s --timeout=2s \
  CMD [ "bash", "docker-healthcheck" ]