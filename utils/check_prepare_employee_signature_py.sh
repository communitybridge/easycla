#!/bin/bash
# API_URL=https://3f13-147-75-85-27.ngrok-free.app (defaults to localhost:5000)
# user_id='9dcf5bbc-2492-11ed-97c7-3e2a23ea20b5'
# company_id='862ff296-6508-4f10-9147-2bc2dd7bfe80'
# project_id='88ee12de-122b-4c46-9046-19422054ed8d'
# DEBUG=1 ./utils/check_prepare_employee_signature_py 9dcf5bbc-2492-11ed-97c7-3e2a23ea20b5 862ff296-6508-4f10-9147-2bc2dd7bfe80 88ee12de-122b-4c46-9046-19422054ed8d
# DEBUG=1 ./utils/check_prepare_employee_signature_py.sh 65d22813-1ac0-4292-bb68-fdcb278473a5 4930fe6e-e023-4f56-9767-6f1996a7b730 43c546ff-bc79-4a32-9454-77dabd6afaee

if [ -z "$1" ]
then
  echo "$0: you need to specify user_id as a 1st parameter"
  exit 1
fi
export user_id="$1"

if [ -z "$2" ]
then
  echo "$0: you need to specify company_id as a 2nd parameter"
  exit 2
fi
export company_id="$2"

if [ -z "$3" ]
then
  echo "$0: you need to specify project_id as a 3rd parameter"
  exit 3
fi
export project_id="$3"

if [ -z "$API_URL" ]
then
  export API_URL="http://localhost:5000"
fi

if [ ! -z "$DEBUG" ]
then
  echo "curl -s -XPOST -H 'Content-Type: application/json' '${API_URL}/v2/check-prepare-employee-signature' -d '{\"project_id\":\"${project_id}\",\"user_id\":\"${user_id}\",\"company_id\":\"${company_id}\"}\"}' | jq -r '.'"
fi
curl -s -XPOST -H "Content-Type: application/json" "${API_URL}/v2/check-prepare-employee-signature" -d "{\"project_id\":\"${project_id}\",\"user_id\":\"${user_id}\",\"company_id\":\"${company_id}\"}" | jq -r '.'
