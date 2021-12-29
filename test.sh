#!/bin/bash
set -e

ROOT_DIR=$(cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd)

for dir in ./*/
do 
    multirepo="$ROOT_DIR/$dir"
    if [ -f "$multirepo/Makefile" ]; then
        cd $multirepo 1>/dev/null
        make test
    else
        echo "WARNING: could not call tests from subrepo "$dir" with make test"
    fi
done

cd $ROOT_DIR