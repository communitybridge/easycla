// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import {Injectable} from '@angular/core';
import {Http} from '@angular/http';

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
    const url: URL = this.getV2Endpoint('/v2/user/' + userId);
    return this.http.get(url).map((res) => res.json());
  }

  getUserWithAuthToken(userId) {
    const url: URL = this.getV2Endpoint('/v2/user/' + userId);
    return this.http.securedGet(url).map((res) => res.json());
  }

  // creates a new account for Gerrit users, with email.
  postOrGetUserForGerrit() {
    const url: URL = this.getV1Endpoint('/v1/user/gerrit');
    return this.http.securedPost(url).map((res) => res.json());
  }

  /**
   * Request to be added to the company Approved List (formerly WhiteList)
   *
   * /user/{user_id}/request-company-whitelist/{company_id}
   */
  requestToBeOnCompanyApprovedList(userId, companyId, projectId, data) {
    //const url: URL = this.getV2Endpoint('/v2/user/' + userId + '/request-company-whitelist/' + companyId);
    const url: URL = this.getV3Endpoint(`/v3/company/${companyId}/ccla-whitelist-requests/${projectId}`);
    return this.http.post(url, data);// no response .map((res) => res.json());
  }

  /**
   * /user/{user_id}/invite-company-admin
   */
  postEmailToCompanyAdmin(userId, data) {
    const url: URL = this.getV2Endpoint('/v2/user/' + userId + '/invite-company-admin');
    return this.http.post(url, data).map((res) => res.json());
  }

  /**
   * /user/{user_id}/active-signature
   */
  getUserSignatureIntent(userId) {
    const url: URL = this.getV2Endpoint('/v2/user/' + userId + '/active-signature');
    return this.http.get(url).map((res) => res.json());
  }

  /**
   * /user/{user_id}/project/{project_id}/last-signature
   **/

  getLastIndividualSignature(userId, projectId) {
    const url: URL = this.getV2Endpoint('/v2/user/' + userId + '/project/' + projectId + '/last-signature');
    return this.http.get(url).map((res) => res.json());
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
    const url: URL = this.getV2Endpoint('/v2/company');
    return this.http.get(url).map((res) => res.json());
  }

  /**
   * /company/{company_id}
   **/
  getCompany(companyId) {
    const url: URL = this.getV2Endpoint('/v2/company/' + companyId);
    return this.http.get(url).map((res) => res.json());
  }

  /**
   * /project/{project_id}
   **/
  getProject(projectId) {
    const url: URL = this.getV2Endpoint('/v2/project/' + projectId);
    return this.http.get(url).map((res) => res.json());
  }

  getProjectWithAuthToken(projectId) {
    const url: URL = this.getV2Endpoint('/v2/project/' + projectId);
    return this.http.securedGet(url).map((res) => res.json());
  }

  /**
   * /request-individual-signature
   **/
  postIndividualSignatureRequest(signatureRequest) {
    const url: URL = this.getV2Endpoint('/v2/request-individual-signature');
    return this.http.post(url, signatureRequest).map((res) => res.json());
  }

  /**
   * /check-prepare-employee-signature
   **/
  postCheckedAndPreparedEmployeeSignature(data) {
    const url: URL = this.getV2Endpoint('/v2/check-prepare-employee-signature');
    return this.http.post(url, data).map((res) => res.json());
  }

  /**
   * /request-employee-signature
   **/
  postEmployeeSignatureRequest(signatureRequest) {
    const url: URL = this.getV2Endpoint('/v2/request-employee-signature');
    return this.http.post(url, signatureRequest).map((res) => res.json());
  }

  /**
   *  /company/{companyID}/ccla-whitelist-requests/{projectID}
   */
  postCCLAWhitelistRequest(companyID, projectID, user) {
    const url: URL = this.getV3Endpoint('/v3/company/' + companyID + '/ccla-whitelist-requests/' + projectID);
    return this.http.post(url, user);
  }

  getGerrit(gerritId) {
    const url: URL = this.getV2Endpoint('/v2/gerrit/' + gerritId);
    return this.http.securedGet(url).map((res) => res.json());
  }

  getProjectGerrits(projectId) {
    const url: URL = this.getV1Endpoint('/v1/project/' + projectId + '/gerrits');
    return this.http.securedGet(url).map((res) => res.json());
  }

  getReleaseVersion() {
    const url: URL = this.getV3Endpoint('/v3/ops/version');
    return this.http.get(url).map((res) => res.json());
  }
}
