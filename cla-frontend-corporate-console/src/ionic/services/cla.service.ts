// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import {Injectable} from '@angular/core';
import {Http} from '@angular/http';
import {AuthService} from './auth.service';

import 'rxjs/Rx';
import 'rxjs/add/operator/catch';
import {Observable} from 'rxjs';

@Injectable()
export class ClaService {
  http: any;
  authService: AuthService;
  claApiUrl: string = '';
  localTesting = false;
  v1ClaAPIURLLocal = 'http://localhost:5000';
  v2ClaAPIURLLocal = 'http://localhost:5000';
  v3ClaAPIURLLocal = 'http://localhost:8080';

  constructor(http: Http, authService: AuthService) {
    this.http = http;
    this.authService = authService;
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

  //////////////////////////////////////////////////////////////////////////////

  /**
   * This service should ONLY contain methods calling CLA API
   **/

  //////////////////////////////////////////////////////////////////////////////

  /**
   * GET /v3/users/{userId}
   **/
  getUserByUserId(userId) {
    const url: URL = this.getV3Endpoint('/v3/users/' + userId);
    return this.http
      .getWithCreds(url)
      .map((res) => {
        return res.json();
      })
      .catch(this.handleServiceError);
  }

  /**
   * GET /v3/users/username/{userName}
   **/
  getUserByUserName(userName) {
    const url: URL = this.getV3Endpoint('/v3/users/username/' + userName);
    return this.http
      .getWithCreds(url)
      .map((res) => res.json())
      .catch((err) => this.handleServiceError(err));
  }

  /**
   * POST /v3/user
   **/
  createUserV3(user) {
    /*
      {
        'user_email': 'user@email.com',
        'user_name': 'User Name',
        'user_company_id': '<org-id>',
        'user_github_id': 12345
      }
     */
    const url: URL = this.getV3Endpoint('/v3/users');
    return this.http
      .post(url, user)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * PUT /v3/user
   **/
  updateUserV3(user) {
    /*
      {
        "lfUsername": "<some-username>",
        "lfEmail": "<some-updated-value>",
        "companyID": "<some-updated-value>"
      }
     */
    const url: URL = this.getV3Endpoint('/v3/users');
    return this.http
      .putWithCreds(url, user)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * GET /v2/user/{user_id}
   **/
  getUser(userId) {
    const url: URL = this.getV2Endpoint('/v2/user/' + userId);
    return this.http
      .get(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  putSignature(signature) {
    /*
      signature: {
        'signature_id': '<signature-id>',
        'signature_type': ('cla' | 'dco'),
        'signature_signed': true,
        'signature_approved': true,
        'signature_sign_url': 'http://sign.com/here',
        'signature_return_url': 'http://cla-system.com/signed',
        'signature_project_id': '<project-id>',
        'signature_reference_id': '<ref-id>',
        'signature_reference_type': ('individual' | 'corporate'),
      }
      */
    const url: URL = this.getV1Endpoint('/v1/signature');
    return this.http
      .put(url, signature)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * GET /v3/signatures/user/{user_id}
   *
   * @param userId the user ID
   * @param pageSize the optional page size - default is 50
   * @param nextKey the next key used when asking for the next page of results
   */
  getSignaturesUser(userId, pageSize = 50, nextKey = '') {
    //const url: URL = this.getV1Endpoint('/v1/signatures/user/' + userId);
    // Leverage the new go backend v3 endpoint
    let path: string = '/v3/signatures/user/' + userId + '?pageSize=' + pageSize;
    if (nextKey != null && nextKey !== '' && nextKey.trim().length > 0) {
      path += '&nextKey=' + nextKey;
    }
    const url: URL = this.getV3Endpoint(path);
    return this.http.getWithCreds(url).map((res) => res.json());
  }

  /**
   * GET /v3/signatures/company/{company_id}
   *
   * @param companyId the company ID
   * @param pageSize the optional page size - default is 50
   * @param nextKey the next key used when asking for the next page of results
   */
  getCompanySignatures(companyId, pageSize = 50, nextKey = '') {
    //const url: URL = this.getV1Endpoint('/v1/signatures/company/' + companyId);
    // Leverage the new go backend v3 endpoint
    let path: string = '/v3/signatures/company/' + companyId + '?pageSize=' + pageSize;
    if (nextKey != null && nextKey !== '' && nextKey.trim().length > 0) {
      path += '&nextKey=' + nextKey;
    }
    const url: URL = this.getV3Endpoint(path);
    return this.http.getWithCreds(url).map((res) => res.json());
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
    //const url: URL = this.getV1Endpoint('/v1/signatures/company/' + companyId + '/project/' + projectId);
    // Leverage the new go backend v3 endpoint - note the slightly different path layout
    let path: string = '/v3/signatures/project/' + projectId + '/company/' + companyId + '?pageSize=' + pageSize;
    if (nextKey != null && nextKey !== '' && nextKey.trim().length > 0) {
      path += '&nextKey=' + nextKey;
    }
    const url: URL = this.getV3Endpoint(path);
    return this.http.getWithCreds(url).map((res) => res.json());
  }

  /**
   * GET /v3/signatures/project/{project_id}/company/{company_id}/employee
   *
   * @param companyId the company ID
   * @param projectId the project ID
   * @param pageSize the optional page size - default is 50
   * @param nextKey the next key used when asking for the next page of results
   */
  getEmployeeProjectSignatures(companyId, projectId, pageSize = 50, nextKey = '') {
    //const url: URL = this.getV1Endpoint('/v1/signatures/company/' + companyId + '/project/' + projectId + '/employee');
    // Leverage the new go backend v3 endpoint - note the different order of the parameters in the path
    let path: string =
      '/v3/signatures/project/' + projectId + '/company/' + companyId + '/employee' + '?pageSize=' + pageSize;
    if (nextKey != null && nextKey !== '' && nextKey.trim().length > 0) {
      path += '&nextKey=' + nextKey;
    }
    const url: URL = this.getV3Endpoint(path);
    return this.http.getWithCreds(url).map((res) => res.json());
  }

  /**
   * GET /v3/signatures/project/{project_id}
   *
   * @param projectId the project ID
   * @param pageSize the optional page size - default is 50
   * @param nextKey the next key used when asking for the next page of results
   */
  getProjectSignatures(projectId, pageSize = 50, nextKey = '') {
    // Leverage the new go backend v3 endpoint - note the slightly different path
    let path: string = '/v3/signatures/project/' + projectId + '?pageSize=' + pageSize;
    if (nextKey != null && nextKey !== '' && nextKey.trim().length > 0) {
      path += '&nextKey=' + nextKey;
    }
    const url: URL = this.getV3Endpoint(path);
    return this.http.getWithCreds(url).map((res) => res.json());
  }

  /**
   * GET /v3/company/user/{userID}/invites Returns list of companies for current
   * user where they are in the access control list or have a pending or rejected invitation
   *
   * Response looks the same as the above getCompaniesByUserManager(userID) except it
   * also includes invitations.
   *
   * of note, when reviewing the output JSON:
   * 1) when the user's name is in the ACL - you should see that the status is “approved”
   * 2) when the user's is NOT in the ACL - you should see the status as “pending” or “rejected”
   */
  getCompaniesByUserManagerWithInvites(userID) {
    const url: URL = this.getV3Endpoint('/v3/company/user/' + userID + '/invites');
    return this.http
      .get(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * GET /v2/company Returns all the companies
   */
  getAllCompanies() {
    const url: URL = this.getV2Endpoint('/v2/company');
    return this.http
      .get(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * GET /v3/company Returns all the companies
   */
  getAllV3Companies() {
    const url: URL = this.getV3Endpoint('/v3/company');
    return this.http
      .get(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * POST /v1/company
   */
  postCompany(company) {
    /*
      {
        'company_name': 'Org Name',
        'company_whitelist': ['safe@email.org'],
        'company_whitelist': ['*@email.org']
      }
     */
    const url: URL = this.getV1Endpoint('/v1/company');
    return this.http
      .post(url, company)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * PUT /v1/company
   */
  putCompany(company) {
    /*
      {
        'company_id': '<company-id>',
        'company_name': 'New Company Name'
      }
     */
    const url: URL = this.getV1Endpoint('/v1/company');
    return this.http
      .put(url, company)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * GET /v2/company/{company_id}
   */
  getCompany(companyId) {
    const url: URL = this.getV2Endpoint('/v2/company/' + companyId);
    return this.http
      .get(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * GET /v1/signature/{signature_id}/manager
   */
  getCLAManagers(signatureId) {
    const url: URL = this.getV1Endpoint('/v1/signature/' + signatureId + '/manager');
    return this.http
      .get(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * POST /v1/signature/{signature_id}/manager
   */
  postCLAManager(signatureId, payload) {
    const url: URL = this.getV1Endpoint('/v1/signature/' + signatureId + '/manager');
    return this.http
      .post(url, { lfid: payload.managerLFID })
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * DELETE /signature/{signature_id}/manager/{lfid}
   */
  deleteCLAManager(projectId, lfid) {
    const url: URL = this.getV1Endpoint('/v1/signature/' + projectId + '/manager/' + lfid);
    return this.http
      .delete(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * GET /v1/project
   */
  getProjects() {
    const url: URL = this.getV1Endpoint('/v1/project');
    return this.http
      .get(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  getCompanyUnsignedProjects(companyId) {
    const url: URL = this.getV1Endpoint('/v1/company/' + companyId + '/project/unsigned');
    return this.http
      .get(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * GET /v2/project/{project_id}
   */
  getProject(projectId) {
    const url: URL = this.getV2Endpoint('/v2/project/' + projectId);
    return this.http
      .get(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * POST /v1/request-corporate-signature
   */
  postCorporateSignatureRequest(signatureRequest) {
    /*
      {
        'project_id': <project-id>,
        'company_id': <company-id>,
        'return_url': <optional-return-url>,
      }
     */
    const url: URL = this.getV1Endpoint('/v1/request-corporate-signature');
    return this.http
      .post(url, signatureRequest)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * /github/login
   */
  githubLogin(companyID, corporateClaID) {
    if (this.localTesting) {
      window.location.assign(
        this.v3ClaAPIURLLocal +
        `/v3/github/login?callback=https://${window.location.host}/#/company/${companyID}/project/${corporateClaID}/orgwhitelist`
      );
    } else {
      window.location.assign(
        this.claApiUrl +
        `/v3/github/login?callback=https://${window.location.host}/#/company/${companyID}/project/${corporateClaID}/orgwhitelist`
      );
    }
  }

  /**
   * GET /signatures/{signatureID}
   **/
  getGithubOrganizationWhitelistEntries(signatureID) {
    const url: URL = this.getV3Endpoint(`/v3/signatures/${signatureID}/gh-org-whitelist`);
    return this.http.getWithCreds(url).map((res) => res.json());
  }

  /**
   * POST /signatures/{signatureID}
   **/
  addGithubOrganizationWhitelistEntry(signatureID, organizationId) {
    const path = `/v3/signatures/${signatureID}/gh-org-whitelist`;
    const data = { organization_id: organizationId };
    const url: URL = this.getV3Endpoint(path);
    return this.http.postWithCreds(url, data).map((res) => res.json());
  }

  /**
   * DELETE /signatures/{signatureID}
   **/
  removeGithubOrganizationWhitelistEntry(signatureID, organizationId) {
    const path = `/v3/signatures/${signatureID}/gh-org-whitelist`;
    const data = { organization_id: organizationId };
    const url: URL = this.getV3Endpoint(path);
    return this.http.deleteWithCredsAndBody(url, data).map((res) => res.json());
  }

  /**
   * POST /v3/company/{companyID}/cla/{corporateClaID}/whitelist/githuborg
   **/
  addGithubOrganizationWhitelist(signatureID, companyID, corporateClaID, organizationId) {
    const path = `/v3/company/${companyID}/cla/${corporateClaID}/whitelist/githuborg`;
    const data = { id: organizationId };
    const url: URL = this.getV3Endpoint(path);
    return this.http.postWithCreds(url, data).map((res) => res.json());
  }

  /**
   * /company/{companyID}/cla/{corporateClaID}/whitelist/githuborg
   **/
  removeGithubOrganizationWhitelist(companyID, corporateClaID, organizationId) {
    const path = `/v3/company/${companyID}/cla/${corporateClaID}/whitelist/githuborg`;
    const data = { id: organizationId };
    const url: URL = this.getV3Endpoint(path);
    return this.http.deleteWithCredsAndBody(url, data).map((res) => res.json());
  }

  /**
   * POST /v3/company/{companyId}/invite-request
   **/
  sendInviteRequestEmail(companyId: string, userId: string, userEmail: string, userName: string): Observable<any> {
    const url: URL = this.getV3Endpoint(`/v3/company/${companyId}/cla/accesslist/request`);
    const data = {
      userID: userId,
      lfUsername: userId,
      lfEmail: userEmail,
      username: userName
    };
    return this.http.putWithCreds(url, data);
  }

  /**
   * GET /v3/company/{companyID}/cla/invitelist:
   * Example:
   * [
   * {
   *      "inviteId": "1e5debac-57fd-4b1f-9669-cfa4c30d7b22",
   *      "status": "pending",
   *      "userEmail": "ahampras@proximabiz.com",
   *      "userLFID": "ahampras",
   *      "userName": "Abhijeet Hampras"
   *  },
   * {
   *      "inviteId": "1e5debac-57fd-4b1f-9669-cfa4c30d7b33",
   *      "status": "rejected",
   *      "userEmail": "davifowl@microsoft.com",
   *      "userLFID": "davidfowl"
   *  }
   * ]
   *
   * @param companyId the company ID
   * @param status the status value - one of: {pending, approved, rejected}
   **/
  getCompanyInvites(companyId: string, status: string) {
    const url: URL = this.getV3Endpoint(`/v3/company/${companyId}/cla/invitelist?status=${status}`);
    return this.http.getWithCreds(url).map((res) => res.json());
  }

  /**
   * GET /v3/company/{companyID}/{userID}/invitelist:
   * Example (pending):
   * {
   *      "inviteId": "1e5debac-57fd-4b1f-9669-cfa4c30d7b22",
   *      "status": "pending",
   *      "userEmail": "ahampras@proximabiz.com",
   *      "userLFID": "ahampras",
   *      "userName": "Abhijeet Hampras",
   *      "companyName: "my company name"
   * }
   *
   * Example (rejected):
   * {
   *      "inviteId": "1e5debac-57fd-4b1f-9669-cfa4c30d7b22",
   *      "status": "rejected",
   *      "userEmail": "ahampras@proximabiz.com",
   *      "userLFID": "ahampras",
   *      "userName": "Abhijeet Hampras",
   *      "companyName: "my company name"
   * }
   *
   * Returns HTTP status 404 if company and user invite is not found
   *
   * @param companyId the company ID
   * @param userID the user ID
   */
  getPendingUserInvite(companyId, userID) {
    const url: URL = this.getV3Endpoint(`/v3/company/${companyId}/${userID}/invitelist`);
    return this.http.getWithCreds(url).map((res) => res.json());
  }

  approveCompanyInvite(companyId: string, requestId: string) {
    const url: URL = this.getV3Endpoint(`/v3/company/${companyId}/cla/accesslist/${requestId}/approve`);
    return this.http.putWithCreds(url);
  }

  rejectCompanyInvite(companyId: string, requestId: string ) {
    const url: URL = this.getV3Endpoint(`/v3/company/${companyId}/cla/accesslist/${requestId}/reject`);
    return this.http.putWithCreds(url);
  }

  /**
   * Creates a new CLA Manager Request with the specified parameters.
   *
   */
  createCLAManagerRequest(companyID: string, projectID: string, userName: string, userEmail: string, userLFID: string) {
    const url: URL = this.getV3Endpoint(`/v3/company/${companyID}/project/${projectID}/cla-manager/requests`);
    const requestBody = {
      userName: userName,
      userEmail: userEmail,
      userLFID: userLFID,
    };
    return this.http.postWithCreds(url, requestBody).map((res) => res.json());
  }

  /**
   * Returns zero or more CLA manager requests based on the user LFID.
   *
   * @param companyID the company ID
   * @param projectID the project ID
   */
  getCLAManagerRequests(companyID: string, projectID: string){
    const url: URL = this.getV3Endpoint(`/v3/company/${companyID}/project/${projectID}/cla-manager/requests`);
    return this.http.getWithCreds(url).map((res) => res.json());
  }

  /**
   * Approves the CLA manager request using the specified parameters.
   *
   * @param companyID the company ID
   * @param projectID the project ID
   * @param requestID the unique request id
   */
  approveCLAManagerRequest(companyID: string, projectID: string, requestID: string) {
    const url: URL = this.getV3Endpoint(`/v3/company/${companyID}/project/${projectID}/cla-manager/request/${requestID}/approve`);
    return this.http.putWithCreds(url).map((res) => res.json());
  }

  /**
   * Deny the CLA manager request using the specified parameters.
   *
   * @param companyID the company ID
   * @param projectID the project ID
   * @param requestID the unique request id
   */
  denyCLAManagerRequest(companyID: string, projectID: string, requestID: string) {
    const url: URL = this.getV3Endpoint(`/v3/company/${companyID}/project/${projectID}/cla-manager/request/${requestID}/deny`);
    return this.http.putWithCreds(url).map((res) => res.json());
  }

  /**
   * Handle service error is a common routine to handle HTTP response errors
   * @param error the error
   */
  private handleServiceError(error: any) {
    if (error.status && error._body && error.status === 500 && error._body.includes('Token is expired')) {
      console.log(error._body + ' - redirecting to login page');
      this.authService.logout();
      window.location.hash = '#/login';
      window.location.reload(true);
      return Observable.throw(error);
    } else if (error.status && error.status === 401) {
      console.log(error._body + ' - redirecting to login page');
      this.authService.logout();
      window.location.hash = '#/login';
      window.location.reload(true);
      return Observable.throw(error);
    }

    console.log('problem invoking service: ');
    console.log(error);
    return Observable.throw(error);
  }

  getReleaseVersion() {
    const url: URL = this.getV3Endpoint('/v3/ops/version');
    return this.http.get(url).map((res) => res.json());
  }

  getProjectWhitelistRequest(companyId: string, projectId: string, status: string) {
    let statusFilter = '';
    if (status != null && status.length > 0) {
      statusFilter = `?status=${status}`;
    }

    const url: URL = this.getV3Endpoint(`/v3/company/${companyId}/ccla-whitelist-requests/${projectId}${statusFilter}`);
    return this.http.get(url).map((res) => res.json())
  }

  approveCclaWhitelistRequest(companyId: string, projectId: string, requestID: string) {
    const url: URL = this.getV3Endpoint(`/v3/company/${companyId}/ccla-whitelist-requests/${projectId}/${requestID}/approve`);
    return this.http.put(url);
  }

  rejectCclaWhitelistRequest(companyId: string, projectId: string, requestID: string) {
    const url: URL = this.getV3Endpoint(`/v3/company/${companyId}/ccla-whitelist-requests/${projectId}/${requestID}/reject`);
    return this.http.put(url);
  }

  //////////////////////////////////////////////////////////////////////////////
}
