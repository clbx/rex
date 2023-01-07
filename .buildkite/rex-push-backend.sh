#!/bin/bash 

# Login to DockerHub
buildah login -u clbx --password-stdin < /secrets/docker-access-token docker.io

buildkite-agent artifact download "rex-backend-container.tar" --step "build-backend" --build $BUILDKITE_BUILD_ID .

container_hash=`buildah pull oci-archive:rex-backend-container.tar`


# tag and push
if [[ -n "$BUILDKITE_TAG" ]]; then
    echo "Tagged Release, taggoing: clbx/rex:$BUILDKITE_TAG"
    buildah tag $container_hash clbx/rex:${BUILDKITE_TAG:1}
    buildah push clbx/rex:${BUILDKITE_TAG:1}
else
    if [[ "$BUILDKITE_BRANCH" == "main" ]]; then
        echo "Main Branch, tagging: clbx/rex:latest"
        buildah tag $container_hash clbx/rex:latest
        buildah push clbx/rex:latest
    else
        echo "Development Branch, tagging: clbx/rex:${BUILDKITE_BRANCH}"
        buildah tag $container_hash clbx/rex:${BUILDKITE_BRANCH}
        buildah push clbx/rex:${BUILDKITE_BRANCH}
    fi
fi

