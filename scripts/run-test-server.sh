#! /bin/sh

printf "INFO: starting test server\n\n"

python3 -m http.server 8080 -b 127.0.0.1 -d "$(pwd)"
