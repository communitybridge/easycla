.PHONY: setup
setup:
	yarn install-frontend

.PHONY: run
run:
	cd src; \
	npm run preserve:dev; \
	npm run serve:dev; \

.PHONY: deploy
deploy:
	yarn deploy -s ${STAGE} -r us-east-1 -c
