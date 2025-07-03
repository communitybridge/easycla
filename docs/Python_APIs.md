# V1/V2 python APIs still used

1. EasyCLA backend [here](https://github.com/linuxfoundation/easycla/blob/main/.github/workflows/deploy-prod.yml#L127).
- `/v2/health`. Health check while `dev`/`prod` deployment.

2. EasyCLA contributor console [here](https://github.com/communitybridge/easycla-contributor-console/blob/main/src/app/core/services/cla-contributor.service.ts), found via `` find.sh . '*' /vN ``.
- `/v1/api/health` [here](https://github.com/communitybridge/easycla-contributor-console/blob/main/test/functional/cypress/integration/api-tests/health-check.spec.ts#L4).
- `/v1/user/gerrit`.
- `/v2/user/<userId>`.
- `/v2/project/<projectId>`.
- `/v2/user/<userId>/active-signature`.
- `/v2/check-prepare-employee-signature`.
- `/v2/request-employee-signature`.
- `/v2/user/<userId>/project/<projectId>/last-signature`.
- `/v2/project/<projectId>`.

3. EasyCLA corporate console [here](https://github.com/LF-Engineering/lfx-corp-cla-console/blob/main/backend/src/data/cla-api.ts):
- No V1 or V2 APIs used.

4. PCC [here](https://github.com/linuxfoundation/lfx-pcc/blob/main/apps/v1-backend/src/modules/cla-services/model/index.ts):
- No V1 or V2 APIs used.

5. GitHub/Gitlab/Gerrit ({provider}) app/bot [here]() (there is no list of which particular APIs are used by GitHub/GitLab/Gerrit):
- `/v2/github/activity`.
- `/v2/repository-provider/{provider}/sign/{installation_id}/{github_repository_id}/{change_request_id}`.
- `/v2/github/installation`.
- `/v1/user/gerrit`.

- `/v2/signed/individual/{installation_id}/{github_repository_id}/{change_request_id}` ?
- `/v2/repository-provider/{provider}/activity` ?
- `/v2/repository-provider/{provider}/oauth2_redirect` ?
- `/v2/signed/gitlab/individual/{user_id}/{organization_id}/{gitlab_repository_id}/{merge_request_id}` ?
- `/v2/signed/gerrit/individual/{user_id}` ?
- `cla-backend/cla/routes.py`: `GitHub Routes`, `Gerrit Routes`, `Gerrit instance routes`.

6. Manually check which APIs were actually called on `dev` and `prod` via:
- `prod` analysis: `` DEBUG=1 NO_ECHO=1 STAGE=prod REGION=us-east-1 DTFROM='10 days ago' DTTO='1 second ago' OUT=api-logs-prod.json ./utils/search_aws_logs.sh 'LG:api-request-path' && jq -r '.[].message' api-logs-prod.json | grep -o 'LG:api-request-path:[^[:space:]]*' | sed 's/^LG:api-request-path://' | sed -E 's/[0-9a-fA-F-]{36}/<uuid>/g' | sed -E 's/\b[0-9]{2,}\b/<id>/g' | sort | uniq -c | sort -nr ``:
```
```
- `dev` analysis (but this can contain API calls made by developer and not actually used): `` DEBUG=1 STAGE=dev REGION=us-east-1 DTFROM='10 days ago' DTTO='1 second ago' OUT=api-logs-dev.json ./utils/search_aws_logs.sh 'LG:api-request-path' && jq -r '.[].message' api-logs-dev.json | grep -o 'LG:api-request-path:[^[:space:]]*' | sed 's/^LG:api-request-path://' | sed -E 's/[0-9a-fA-F-]{36}/<uuid>/g' | sed -E 's/\b[0-9]{2,}\b/<id>/g' | sort | uniq -c | sort -nr ``:
```
    113 /v2/project/<uuid>
     53 /v2/user/<uuid>
     39 /v2/github/activity
     14 /v2/user/<uuid>/active-signature
     12 /v2/user-from-token
     12 /v2/repository-provider/github/sign/<id>/<id>/<id>
      8 /v2/github/installation
      5 /v2/health
      4 /v2/project/null
      1 /v2/users/company/abcd
      1 /v2/user-from-session
      1 /v2/return-url/<uuid>
      1 /v2/repository-provider/github/sign/<id>/<id>/1
      1 /v2/check-prepare-employee-signature
      1 /v2/.env
      1 /v1/users/company/abcd
      1 /v1/.env
```
