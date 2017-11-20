#!/bin/bash
# http://stackoverflow.com/questions/59895/getting-the-source-directory-of-a-bash-script-from-within
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PYTHONPATH=$PYTHONPATH:$DIR
#cd tests/
nose2 --coverage=cla --coverage-report=html
echo "Open htmlcov/index.html in your favorite browser for more details."
