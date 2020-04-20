// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Injectable } from '@angular/core';
import { Http } from '@angular/http';

import 'rxjs/Rx';

@Injectable()
export class ClaService {
  http: any;
  claApiUrl: string = '';
  localTesting = false;
  v1ClaAPIURLLocal = 'http://localhost:5000';
  v2ClaAPIURLLocal = 'http://localhost:5000';
  v3ClaAPIURLLocal = 'http://localhost:8080';

  constructor(http: Http) {
    this.http = http;
  }

  /**
   * Constructs a URL based on the path and endpoint host:port.
   * @param path the URL path
   * @returns a URL to the V1 endpoint with the specified path. If running in local mode, the endpoint will point to a
   * local host:port - otherwise the endpoint will point the appropriate environment endpoint running in the cloud.
   */
  private getV1Endpoint(path: string) {
    let url: URL;
    if (this.localTesting) {
      url = new URL(this.v1ClaAPIURLLocal + path);
    } else {
      url = new URL(this.claApiUrl + path);
    }
    return url;
  }

  /**
   * Constructs a URL based on the path and endpoint host:port.
   * @param path the URL path
   * @returns a URL to the V2 endpoint with the specified path. If running in local mode, the endpoint will point to a
   * local host:port - otherwise the endpoint will point the appropriate environment endpoint running in the cloud.
   */
  private getV2Endpoint(path: string) {
    let url: URL;
    if (this.localTesting) {
      url = new URL(this.v2ClaAPIURLLocal + path);
    } else {
      url = new URL(this.claApiUrl + path);
    }
    return url;
  }

  /**
   * Constructs a URL based on the path and endpoint host:port.
   * @param path the URL path
   * @returns a URL to the V3 endpoint with the specified path. If running in local mode, the endpoint will point to a
   * local host:port - otherwise the endpoint will point the appropriate environment endpoint running in the cloud.
   */
  private getV3Endpoint(path: string) {
    let url: URL;
    if (this.localTesting) {
      url = new URL(this.v3ClaAPIURLLocal + path);
    } else {
      url = new URL(this.claApiUrl + path);
    }
    return url;
  }

  public isLocalTesting(flag: boolean) {
    if (flag) {
      console.log('Running in local services mode');
    } else {
      console.log('Running in deployed services mode');
    }
    this.localTesting = flag;
  }

  public setApiUrl(claApiUrl: string) {
    this.claApiUrl = claApiUrl;
  }

  public setHttp(http: any) {
    this.http = http; // allow configuration for alternate http library
  }

  /**
   * /user/{user_id}
   */
  getUser(userId) {
    if (this.localTesting) {
      return this.http.get(this.v2ClaAPIURLLocal + '/v2/user/' + userId).map((res) => res.json());
    } else {
      return this.http.get(this.claApiUrl + '/v2/user/' + userId).map((res) => res.json());
    }
  }

  getUserWithAuthToken(userId) {
    if (this.localTesting) {
      return this.http.securedGet(this.v2ClaAPIURLLocal + '/v2/user/' + userId).map((res) => res.json());
    } else {
      return this.http.securedGet(this.claApiUrl + '/v2/user/' + userId).map((res) => res.json());
    }
  }

  // creates a new account for Gerrit users, with email.
  postOrGetUserForGerrit() {
    if (this.localTesting) {
      return this.http.securedPost(this.v1ClaAPIURLLocal + '/v1/user/gerrit').map((res) => res.json());
    } else {
      return this.http.securedPost(this.claApiUrl + '/v1/user/gerrit').map((res) => res.json());
    }
  }

  /**
   * /user/{user_id}/request-company-whitelist/{company_id}
   */
  postUserMessageToCompanyManager(userId, companyId, data) {
    if (this.localTesting) {
      return this.http
        .post(this.v2ClaAPIURLLocal + '/v2/user/' + userId + '/request-company-whitelist/' + companyId, data)
        .map((res) => res.json());
    } else {
      return this.http
        .post(this.claApiUrl + '/v2/user/' + userId + '/request-company-whitelist/' + companyId, data)
        .map((res) => res.json());
    }
  }

  /**
   * /user/{user_id}/invite-company-admin
   */
  postEmailToCompanyAdmin(userId, data) {
    if (this.localTesting) {
      return this.http
        .post(this.v2ClaAPIURLLocal + '/v2/user/' + userId + '/invite-company-admin/', data)
        .map((res) => res.json());
    } else {
      return this.http
        .post(this.claApiUrl + '/v2/user/' + userId + '/invite-company-admin/', data)
        .map((res) => res.json());
    }
  }

  /**
   * /user/{user_id}/active-signature
   */
  getUserSignatureIntent(userId) {
    if (this.localTesting) {
      return this.http.get(this.v2ClaAPIURLLocal + '/v2/user/' + userId + '/active-signature').map((res) => res.json());
    } else {
      return this.http.get(this.claApiUrl + '/v2/user/' + userId + '/active-signature').map((res) => res.json());
    }
  }

  /**
   * /user/{user_id}/project/{project_id}/last-signature
   **/

  getLastIndividualSignature(userId, projectId) {
    if (this.localTesting) {
      return this.http
        .get(this.v2ClaAPIURLLocal + '/v2/user/' + userId + '/project/' + projectId + '/last-signature')
        .map((res) => res.json());
    } else {
      return this.http
        .get(this.claApiUrl + '/v2/user/' + userId + '/project/' + projectId + '/last-signature')
        .map((res) => res.json());
    }
  }

  /**
   * GET /v3/signatures/project/{project_id}/company/{company_id}
   *
   * @param companyId the company ID
   * @param projectId the project ID
   * @param pageSize the optional page size - default is 50
   * @param nextKey the next key used when asking for the next page of results
   */
  getCompanyProjectSignatures(companyId, projectId, pageSize = 50, nextKey = '') {
    let path: string = '/v3/signatures/project/' + projectId + '/company/' + companyId + '?pageSize=' + pageSize;
    if (nextKey != null && nextKey !== '' && nextKey.trim().length > 0) {
      path += '&nextKey=' + nextKey;
    }
    const url: URL = this.getV3Endpoint(path);
    return this.http.getWithCreds(url).map((res) => res.json());
  }

  getAllCompanies() {
    return this.http.get(this.claApiUrl + '/v2/company').map((res) => res.json());
  }

  /**
   * /company/{company_id}
   **/

  getCompany(companyId) {
    return this.http.get(this.claApiUrl + '/v2/company/' + companyId).map((res) => res.json());
  }

  /**
   * /project/{project_id}
   **/
  getProject(projectId) {
    return this.http.get(this.claApiUrl + '/v2/project/' + projectId).map((res) => res.json());
  }

  getProjectWithAuthToken(projectId) {
    return this.http.securedGet(this.claApiUrl + '/v2/project/' + projectId).map((res) => res.json());
  }

  /**
   * /request-individual-signature
   **/

  postIndividualSignatureRequest(signatureRequest) {
    return this.http
      .post(this.getV2Endpoint('/v2/request-individual-signature'), signatureRequest)
      .map((res) => res.json());
  }

  /**
   * /check-prepare-employee-signature
   **/

  postCheckedAndPreparedEmployeeSignature(data) {
    return this.http.post(this.claApiUrl + '/v2/check-prepare-employee-signature', data).map((res) => res.json());
  }

  /**
   * /request-employee-signature
   **/

  postEmployeeSignatureRequest(signatureRequest) {
    return this.http.post(this.claApiUrl + '/v2/request-employee-signature', signatureRequest).map((res) => res.json());
  }

  /**
   *  /company/{companyID}/ccla_whitelist_requests/{projectID}
   */
  postCCLAWhitelistRequest(companyID, projectID, user) {

    if (this.localTesting) {
      return this.http.post(this.v3ClaAPIURLLocal + '/v3/company/' + companyID + '/ccla_whitelist_requests/' + projectID, user);
    } else {
      return this.http.post(this.claApiUrl + '/v3/company/' + companyID + '/ccla_whitelist_requests/' + projectID, user);
    }
  }

  getGerrit(gerritId) {
    return this.http.securedGet(this.claApiUrl + '/v2/gerrit/' + gerritId).map((res) => res.json());
  }

  getProjectGerrits(projectId) {
    return this.http.securedGet(this.claApiUrl + '/v1/project/' + projectId + '/gerrits').map((res) => res.json());
  }

  getReleaseVersion() {
    const url: URL = this.getV3Endpoint('/v3/ops/version');
    return this.http.get(url).map((res) => res.json());
  }

}
