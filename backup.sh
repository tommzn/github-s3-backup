#!/bin/bash

GS3B_VERSION="v0.0.1"
GS3B_REPOSITORIES="repositories.list"
GS3B_WORKDIR="./"

printf "\nS3 Backup for Github ($GS3B_VERSION)\n"

while [[ $# -gt 0 ]]; do
    case $1 in
        --repos)
            GS3B_REPOSITORIES="$2"
            shift # past argument
            shift # past value
            ;;
        --s3-bucket)
            GS3B_BUCKET="s3://$2/"
            shift # past argument
            shift # past value
            ;;
        --s3-target-path)
            GS3B_TARGETPATH="$2"
            shift # past argument
            shift # past value
            ;;
        --workdir)
            GS3B_WORKDIR="$2"
            shift # past argument
            shift # past value
            ;;
        *)
            shift # past argument
            ;;
    esac
done

if [ -z "$GS3B_BUCKET" ]; then
    printf "\nMissing S3 bucket. Use --s3-bucket to specify a backup bucket"
    exit 1
fi

while read -r repo; do 

    printf "\n Backup: $repo"

    printf "\n Cloning: $repo"
    GS3B_TARGETDIR=$(echo "$repo" | sed -e 's#https://github.com/##' -e 's/.git//' -e 's#/#.#')
    git clone --mirror "$repo" "$GS3B_WORKDIR$GS3B_TARGETDIR"

    printf "\n Archive: $repo"
    GS3B_ARCHIVE="$GS3B_WORKDIR$GS3B_TARGETDIR.tar.gz"
    tar -czv -C "$GS3B_WORKDIR" -f "$GS3B_ARCHIVE" "$GS3B_TARGETDIR"

    printf "\n Upload: $GS3B_ARCHIVE to $GS3B_BUCKET$GS3B_TARGETPATH"
    aws s3 cp "$GS3B_ARCHIVE" "$GS3B_BUCKET$GS3B_TARGETPATH"

    printf "\n Cleanup: $GS3B_WORKDIR$GS3B_TARGETDIR"
    rm -rf "$GS3B_WORKDIR$GS3B_TARGETDIR"
    rm -rf "$GS3B_ARCHIVE"

     printf "\n Finished: $repo"
     
done < "$GS3B_REPOSITORIES"