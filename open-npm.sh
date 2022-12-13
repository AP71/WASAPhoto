#!/usr/bin/env sh

docker run -it --rm -v "$(pwd):/src" --network host --workdir /src/webui node:lts /bin/bash
