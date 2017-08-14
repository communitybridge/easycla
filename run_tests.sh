#!/bin/bash
# http://stackoverflow.com/questions/59895/getting-the-source-directory-of-a-bash-script-from-within
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PYTHONPATH=$PYTHONPATH:$DIR
cd tests/
nosetests --with-coverage --cover-package=cla --cover-html --logging-clear-handlers
echo "Open tests/cover/index.html in your favorite browser for more details."
