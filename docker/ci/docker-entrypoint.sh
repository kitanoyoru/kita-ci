#!/bin/bash
set -ex

GIT_URL=$1
REPO_BRANCH=$2
TAG=$3

mkdir repo && git clone $GIT_URL repo
if [$? -eq 0]; then
  cd repo
  git checkout $REPO_BRANCH
  if [$? -eq 0]; then
    docker build -t kita .
    if [$? -eq 0]; then
      docker tag kita $TAG
      # TODO: Login to DockerHub registry
      docker push $TAG
    else
      exit $?
    fi
  else
    exit $?
  fi
else
  exit $?
fi
