# /bin/bash
VERSION=v1.0.3
IMAGENAME=patient-edge:${VERSION}
docker build . -t ${IMAGENAME}
