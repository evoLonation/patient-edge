# /bin/bash
VERSION=v1.0.4
IMAGENAME=patient-edge:${VERSION}
docker build . -t ${IMAGENAME}
