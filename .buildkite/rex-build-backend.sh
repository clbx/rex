#!/bin/bash

#Download and unarchive repository artifact
buildkite-agent artifact download "rex-repository.tar" --step "init" --build $BUILDKITE_BUILD_ID .
tar -xvf rex-repository.tar
cd rex

# Build Backend Container
buildah bud --layers -f backend.Dockerfile -t clbx/rex .

# Export the image to a tarball
buildah push clbx/rex oci-archive:/build/rex-backend-container.tar 
# Upload artifact 
buildkite-agent artifact upload rex-backend-container.tar
