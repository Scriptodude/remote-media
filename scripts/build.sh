#!/bin/bash
# usage:
# chmod a+x ./build.sh
# ./build.sh

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

cd $SCRIPT_DIR/../web/_remotemediaweb
flutter build web -o $(pwd)/../release/

cd $SCRIPT_DIR

./service-linux.sh
