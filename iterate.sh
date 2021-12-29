#!/bin/bash

#
# Iterate over file running make $command
#

set -e # stop if error

command=$1

ROOT_DIR=$(cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd)

for dir in ./*/
do 
    multirepo="$ROOT_DIR/$dir"
    if [ -f "$multirepo/Makefile" ]; then
        cd $multirepo
        bash -c "make $command"
    else
        echo "WARNING: Makefile not found from subrepo "$dir""
    fi
done

cd $ROOT_DIR