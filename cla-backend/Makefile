.PHONY: setup
setup:
	npm install;
	pip install virtualenv;
	virtualenv -p python3 ~/.env/lf-cla;
	source ~/.env/lf-cla/bin/activate; \
	pip install -r requirements.txt; \

.PHONY: deploy
deploy:
	source ~/.env/lf-cla/bin/activate; \
	node_modules/.bin/serverless deploy --stage ${STAGE} --region us-east-1; \
