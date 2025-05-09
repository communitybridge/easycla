#!/bin/bash
# API_URL=https://3f13-147-75-85-27.ngrok-free.app (defaults to localhost:5000)
# user_id='9dcf5bbc-2492-11ed-97c7-3e2a23ea20b5'
# project_id='88ee12de-122b-4c46-9046-19422054ed8d'
# return_url_type='github'
# return_url='http://localhost'
# TOKEN='...' - Auth0 JWT bearer token
# XACL='...' - X-ACL header
# DEBUG=1 TOKEN="$(cat ./auth0.token.secret)" XACL="$(cat ./x-acl.secret)" ./utils/request_individual_signature_go_post.sh 9dcf5bbc-2492-11ed-97c7-3e2a23ea20b5 88ee12de-122b-4c46-9046-19422054ed8d github 'http://localhost'

if [ -z "$TOKEN" ]
then
  # source ./auth0_token.secret
  TOKEN="$(cat ./auth0.token.secret)"
fi

if [ -z "$TOKEN" ]
then
  echo "$0: TOKEN not specified and unable to obtain one"
  exit 1
fi

if [ -z "$XACL" ]
then
  XACL="$(cat ./x-acl.secret)"
fi

if [ -z "$XACL" ]
then
  echo "$0: XACL not specified and unable to obtain one"
  exit 2
fi

if [ -z "$1" ]
then
  echo "$0: you need to specify user_id as a 1st parameter"
  exit 3
fi
export user_id="$1"

if [ -z "$2" ]
then
  echo "$0: you need to specify project_id as a 2nd parameter"
  exit 4
fi
export project_id="$2"

if [ -z "$3" ]
then
  echo "$0: you need to specify return_url_type as a 3rd parameter: github|gitlab|gerrit"
  exit 5
fi
export return_url_type="$3"

if [ -z "$4" ]
then
  echo "$0: you need to specify return_url as a 4th parameter"
  exit 6
fi
export return_url="$4"

if [ -z "$API_URL" ]
then
  export API_URL="http://localhost:5000"
fi

if [ ! -z "$DEBUG" ]
then
  echo "curl -s -XPOST -H 'X-ACL: ${XACL}' -H 'Authorization: Bearer ${TOKEN}' -H 'Content-Type: application/json' '${API_URL}/v4/request-individual-signature' -d '{\"project_id\":\"${project_id}\",\"user_id\":\"${user_id}\",\"return_url_type\":\"${return_url_type}\",\"return_url\":\"${return_url}\"}' | jq -r '.'"
fi
curl -s -XPOST -H "X-ACL: ${XACL}" -H "Authorization: Bearer ${TOKEN}" -H "Content-Type: application/json" "${API_URL}/v4/request-individual-signature" -d "{\"project_id\":\"${project_id}\",\"user_id\":\"${user_id}\",\"return_url_type\":\"${return_url_type}\",\"return_url\":\"${return_url}\"}" | jq -r '.'
