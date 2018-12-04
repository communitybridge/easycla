# This is a temporary script to deploy all services from dev env

# change this to qa if you are deploting to qa
export STAGE="mani"

# change the REPO_ROOT to your git repo root path
export REPO_ROOT="/home/manikantanr/mani-desk-copy-backup"

deploy_front_end() {
cd "${REPO_ROOT}"
yarn install
# yarn install-aws-profile # do this if you don't configure locally. You should have your credetials at ~/.aws/credentials
cd "${PROJECT_DIR}"
yarn install-frontend
yarn deploy -s "${STAGE}" -r us-east-1 -c 

}


export PROJECT_DIR="${REPO_ROOT}/frontend-project-management-console"
deploy_front_end

export PROJECT_DIR="${REPO_ROOT}/cla-frontend-corporate-console"
deploy_front_end

export PROJECT_DIR="${REPO_ROOT}/cla-frontend-console"
deploy_front_end



echo "Deploying backend"
cd "${REPO_ROOT}"

# export AWS keys

cd cla-backend && npm install && node_modules/.bin/serverless deploy --stage mani --region us-east-1