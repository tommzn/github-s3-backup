from alpine:latest 

RUN apk add --no-cache aws-cli

WORKDIR /opt/github-s3-backup

COPY backup.sh repositories.list ./

ENTRYPOINT ["/opt/github-s3-backup/backup.sh"]
