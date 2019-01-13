.PHONY: setup
setup:
	npm install;
	pip3 install virtualenv;
	virtualenv -p python3 ~/.env/lf-cla;
	source ~/.env/lf-cla/bin/activate; \
	pip3 install -r requirements.txt; \

.PHONY: add_user
add_user:
	aws dynamodb put-item \
		--table-name "cla-local-user-permissions" \
		--item "{ \"user_id\": { \"S\": \"$(user_sub)\" }, \"projects\": { \"SS\": [\"a09J000000KHoZVIA1\",\"a09J000000KHoayIAD\"] } }" \
		--profile lf-cla --region "us-east-1" \
		--endpoint-url http://localhost:8000;

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
