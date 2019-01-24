.PHONY: setup
setup:
	npm install;
	pip3 install virtualenv;
	virtualenv -p python3 ~/.env/lf-cla;
	source ~/.env/lf-cla/bin/activate; \
	pip3 install -r requirements.txt; \

.PHONY: setup_circle
setup_circle:
	npm install; \
	sudo pip install -r requirements.txt; \

# username: LFID username
# projects: quoted list of SFDC project IDs. e.g. '\"abcd\",\"1234\"' 
.PHONY: add_user
add_user:
	aws dynamodb put-item \
		--table-name "cla-$(STAGE)-user-permissions" \
		--item "{ \"username\": { \"S\": \"$(username)\" }, \"projects\": { \"SS\": [$(projects)] } }" \
		--region "us-east-1";

# local_file: path to local file
# remote_filename: logo name, in the form of <SFDC_ID>.png
.PHONY: upload_logo
upload_logo:
	aws s3 cp $(local_file) s3://cla-project-logo-$(STAGE)/$(remote_filename)

.PHONY: run_dynamo
run_dynamo:
	source ~/.env/lf-cla/bin/activate; \
	node_modules/.bin/serverless dynamodb install -s 'local'; \
	node_modules/.bin/serverless dynamodb start -s 'local'; \

.PHONY: run_s3
run_s3:
	source ~/.env/lf-cla/bin/activate; \
	node_modules/.bin/serverless s3 start -s 'local'; \

.PHONY: run_lambda
run_lambda:
	source ~/.env/lf-cla/bin/activate; \
	node_modules/.bin/serverless wsgi serve -s 'local'; \

.PHONY: deploy
deploy:
	# Stop deployment if env.json exists
	test ! -s ./env.json || exit 1

	source ~/.env/lf-cla/bin/activate; \
	node_modules/.bin/serverless deploy --stage ${STAGE} --region us-east-1; \
