#!/bin/bash
set -xeuo pipefail

# Add key to ssh-agent
mkdir -p ~/.ssh
eval "$(ssh-agent -s)"
ssh-add -k /secrets/ssh-key

# Add github to known hosts
ssh-keyscan github.com >> ~/.ssh/known_hosts
# Clone repository
git clone $BUILDKITE_REPO
pushd $BUILDKITE_PIPELINE_SLUG
git checkout $BUILDKITE_COMMIT
popd
tar -cvf $BUILDKITE_PIPELINE_SLUG-repository.tar $BUILDKITE_PIPELINE_SLUG
buildkite-agent artifact upload $BUILDKITE_PIPELINE_SLUG-repository.tar