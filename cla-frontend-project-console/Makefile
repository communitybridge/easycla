# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

.PHONY: setup
setup:
	yarn install-frontend

.PHONY: run
run:
	cd src; \
	npm run preserve:${STAGE}; \
	npm run serve:${STAGE}; \

.PHONY: deploy
deploy:
	yarn deploy -s ${STAGE} -r us-east-1 -c
