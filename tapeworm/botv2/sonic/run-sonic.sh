#!/bin/bash

# -v /path/to/your/sonic/store/:/var/lib/sonic/store/
docker run \
    -p 1491:1491 \
    -v $PWD/sonic.cfg:/etc/sonic.cfg:ro \
    valeriansaliou/sonic:v1.3.0
