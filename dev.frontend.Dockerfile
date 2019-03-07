FROM node:10-alpine 
ARG               PROJECT_DIR=/go/src/git.metrosystems.net/reliability-engineering/ustress/
RUN               mkdir -p ${PROJECT_DIR}/web/ui 
WORKDIR           ${PROJECT_DIR}/web/ui
COPY              web/ui .
RUN               npm install
RUN               ls -al
CMD               ["npm","start"]
