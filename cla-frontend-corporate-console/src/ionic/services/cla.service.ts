// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Injectable } from '@angular/core';
import { Http } from '@angular/http';
import { AuthService } from './auth.service';

import 'rxjs/Rx';
import 'rxjs/add/operator/catch';
import { Observable } from 'rxjs';

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
   * GET /v1/user
   **/
  getUsers() {
    const url: URL = this.getV1Endpoint('/v1/user');
    return this.http
      .get(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * GET /v3/users/{userId}
   **/
  getUserByUserId(userId) {
    console.log(userId, 'userIduserId')
    const url: URL = this.getV3Endpoint('/v3/users/' + userId);
    return this.http
      .getWithCreds(url)
      .map((res) => {
        console.log(res, 'res')
        return res.json()
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
   * POST /v1/user
   **/
  createUser(user) {
    /*
      {
        'user_email': 'user@email.com',
        'user_name': 'User Name',
        'user_company_id': '<org-id>',
        'user_github_id': 12345
      }
     */
    const url: URL = this.getV1Endpoint('/v1/user');
    return this.http
      .post(url, user)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
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
   * PUT /v1/user
   **/
  putUser(user) {
    /*
      {
        'user_id': '<user-id>',
        'user_email': 'user@email.com',
        'user_name': 'User Name',
        'user_company_id': '<org-id>',
        'user_github_id': 12345
      }
     */
    const url: URL = this.getV1Endpoint('/v1/user');
    return this.http
      .put(url, user)
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

  /**
   * DELETE /v1/user/{user_id}
   **/
  deleteUser(userId) {
    const url: URL = this.getV1Endpoint('/v1/user/' + userId);
    return this.http
      .delete(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * GET /v1/user/email/{user_email}
   **/
  getUserByEmail(userEmail) {
    const url: URL = this.getV1Endpoint('/v1/user/email/' + userEmail);
    return this.http
      .get(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * GET /v1/user/github/{user_github_id}
   **/
  getUserByGithubId(userGithubId) {
    const url: URL = this.getV1Endpoint('/v1/user/github/' + userGithubId);
    return this.http
      .get(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * GET /v1/user/{user_id}/signatures
   **/
  getUserSignatures(userId) {
    const url: URL = this.getV1Endpoint('/v1/user/' + userId + '/signatures');
    return this.http
      .get(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * GET /v1/users/company/{user_company_id}
   **/
  getUsersByCompanyId(userCompanyId) {
    const url: URL = this.getV1Endpoint('/v1/users/company/' + userCompanyId);
    return this.http
      .get(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * POST /v2/user/{user_id}/request-company-whitelist/{company_id}
   **/
  postUserMessageToCompanyManager(userId, companyId, message) {
    /*
      message: {
        'message': 'custom message to manager'
      }
      */
    const url: URL = this.getV2Endpoint('/v2/user/' + userId + '/request-company-whitelist/' + companyId);
    return this.http
      .post(url, message)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * GET /v2/user/{user_id}/active-signature
   **/
  getUserSignatureIntent(userId) {
    const url: URL = this.getV2Endpoint('/v2/user/' + userId + '/active-signature');
    return this.http
      .get(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * GET /v2/user/{user_id}/project/{project_id}/last-signature
   **/
  getLastIndividualSignature(userId, projectId) {
    const url: URL = this.getV2Endpoint('/v2/user/' + userId + '/project/' + projectId + '/last-signature');
    return this.http
      .get(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * GET /v1/signature
   */
  getSignatures() {
    const url: URL = this.getV1Endpoint('/v1/signature');
    return this.http
      .get(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  postSignature(signature) {
    /*
      signature: {
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
      .post(url, signature)
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
   * GET /v3/signature/{signature_id}
   *
   * @param signatureId the signature ID
   * @param pageSize the optional page size - default is 50
   * @param nextKey the next key used when asking for the next page of results
   */
  getSignature(signatureId, pageSize = 50, nextKey = '') {
    //const url: URL = this.getV1Endpoint('/v1/signature/' + signatureId);
    // Leverage the new go backend v3 endpoint
    let path: string = '/v3/signature/' + signatureId + '?pageSize=' + pageSize;
    if (nextKey != null && nextKey !== '' && nextKey.trim().length > 0) {
      path += '&nextKey=' + nextKey;
    }
    const url: URL = this.getV3Endpoint(path);
    return this.http.getWithCreds(url).map((res) => res.json());
  }

  /**
   * DELETE /v1/signature/{signatureId}
   * @param signatureId the signature id
   */
  deleteSignature(signatureId) {
    const url: URL = this.getV1Endpoint('/v1/signatures/' + signatureId);
    return this.http
      .delete(url)
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
    //const url: URL = this.getV1Endpoint('/v1/signatures/project/' + projectId);
    // Leverage the new go backend v3 endpoint - note the slightly different path
    let path: string = '/v3/signatures/project/' + projectId + '?pageSize=' + pageSize;
    if (nextKey != null && nextKey !== '' && nextKey.trim().length > 0) {
      path += '&nextKey=' + nextKey;
    }
    const url: URL = this.getV3Endpoint(path);
    return this.http.getWithCreds(url).map((res) => res.json());
  }

  /**
   * GET /v1/repository
   */
  getRepositories() {
    const url: URL = this.getV1Endpoint('/v1/repository');
    return this.http
      .get(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * POST /v1/repository
   * @param repository the repository payload
   */
  postRepository(repository) {
    /*
      repository: {
        'repository_project_id': '<project-id>',
        'repository_external_id': 'repo1',
        'repository_name': 'Repo Name',
        'repository_type': 'github',
        'repository_url': 'http://url-to-repo.com'
      }
     */
    const url: URL = this.getV1Endpoint('/v1/repository');
    return this.http
      .post(url, repository)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * PUT /v1/repository
   * @param repository the repository payload
   */
  putRepository(repository) {
    /*
      repository: {
        'repository_id': '<repo-id>',
        'repository_project_id': '<project-id>',
        'repository_external_id': 'repo1',
        'repository_name': 'Repo Name',
        'repository_type': 'github',
        'repository_url': 'http://url-to-repo.com'
      }
     */
    const url: URL = this.getV1Endpoint('/v1/repository');
    return this.http
      .put(url, repository)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * GET /repository/{repository_id}
   */
  getRepository(repositoryId) {
    const url: URL = this.getV1Endpoint('/v1/repository/' + repositoryId);
    return this.http
      .get(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * DELETE /repository/{repository_id}
   */
  deleteRepository(repositoryId) {
    const url: URL = this.getV1Endpoint('/v1/repository/' + repositoryId);
    return this.http
      .delete(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * GET /v1/company Returns list of companies for current user
   */
  getCompanies() {
    const url: URL = this.getV1Endpoint('/v1/company');
    return this.http
      .get(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * GET /v3/company/user/{userID} Returns list of companies for current user where they are in the access control list
   *
   * Typical response might look like:
   * {
   * "companies": [
   *     {
   *         "companyACL": [
   *             "ngozi"
   *         ],
   *         "companyID": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeeJ",
   *         "companyName": "Delta Org",
   *         "created": "2019-08-21T15:42:52.657Z",
   *         "updated": "2019-08-21T15:42:52.657Z"
   *     },
   *     {
   *         "companyACL": [
   *             "ngozi"
   *         ],
   *         "companyID": "aaeaa8aa-bbbb-cccc-d3dd-eeeeeeeeeeeZ",
   *         "companyName": "ngozi",
   *         "created": "2019-08-19T20:16:17.656Z",
   *         "updated": "2019-08-19T20:16:17.656Z"
   *     },
   *     {
   *         "companyACL": [
   *             "ngozi"
   *         ],
   *         "companyID": "aa8aaaaa-bbbb-cccc-d4dd-eeeeeeeeeeeF",
   *         "companyName": "delta org",
   *         "created": "2019-09-03T13:46:47.446Z",
   *         "updated": "2019-09-03T13:46:47.446Z"
   *     },
   *     {
   *         "companyACL": [
   *             "ngozi"
   *         ],
   *         "companyID": "aa5aa5aa-bbbb-cccc-d8dd-eeeeeeeeeeet",
   *         "companyName": "ghghghg",
   *         "created": "2019-07-17T15:10:15.499Z",
   *         "updated": "2019-07-17T15:10:15.499Z"
   *     },
   *     {
   *         "companyACL": [
   *             "ngozi"
   *         ],
   *         "companyID": "aacaadaa-bbbb-cccc-d3dd-eeeeeeeeeeew",
   *         "companyName": "ahahaa",
   *         "created": "2019-07-19T15:49:43.424Z",
   *         "updated": "2019-07-19T15:49:43.424Z"
   *     }
   * ],
   * "resultCount": 5,
   * "totalCount": 560
   *}
   */
  getCompaniesByUserManager(userID) {
    const url: URL = this.getV3Endpoint('/v3/company/user/' + userID);
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
   * GET /v3/company/search?companyName={company_name}
   */
  searchCompaniesByName(companyName) {
    const url: URL = this.getV3Endpoint('/v3/company/search?companyName=' + companyName);
    return this.http.getWithCreds(url).map((res) => res.json());
  }

  /**
   * DELETE /v1/company/{company_id}
   */
  deleteCompany(companyId) {
    const url: URL = this.getV1Endpoint('/v1/company/' + companyId);
    return this.http
      .delete(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * GET /v1/project/{project_id}/manager
   */
  getProjectManagers(projectId) {
    const url: URL = this.getV1Endpoint('/v1/project/' + projectId + '/manager');
    return this.http
      .get(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * POST /v1/project/{project_id}/manager
   */
  postProjectManager(projectId, payload) {
    const url: URL = this.getV1Endpoint('/v1/project/' + projectId + '/manager');
    return this.http
      .post(url, { lfid: payload.managerLFID })
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * DELETE /v1/project/{project_id}/manager/{lfid}
   */
  deleteProjectManager(projectId, payload) {
    const url: URL = this.getV1Endpoint('/v1/project/' + projectId + '/manager/' + payload.lfid);
    return this.http
      .delete(url)
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
  deleteCLAManager(projectId, payload) {
    const url: URL = this.getV1Endpoint('/v1/signature/' + projectId + '/manager/' + payload.lfid);
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

  postProject(project) {
    /*
      {
        'project_external_id': '<proj-external-id>',
        'project_name': 'Project Name',
        'project_ccla_enabled': True,
        'project_ccla_requires_icla_signature': True,
        'project_icla_enabled': True
      }
     */
    const url: URL = this.getV1Endpoint('/v1/project');
    return this.http
      .post(url, project)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  putProject(project) {
    /*
      {
        'project_id': '<project-id>',
        'project_name': 'New Project Name'
      }
     */
    const url: URL = this.getV1Endpoint('/v1/project');
    return this.http
      .put(url, project)
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

  getProjectsByExternalId(externalId) {
    const url: URL = this.getV1Endpoint('/v1/project/external/' + externalId);
    return this.http
      .get(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  deleteProject(projectId) {
    const url: URL = this.getV1Endpoint('/v1/project/' + projectId);
    return this.http
      .delete(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * GET /v1/project/{project_id}/repositories
   */
  getProjectRepositories(projectId) {
    const url: URL = this.getV1Endpoint('/v1/project/' + projectId + '/repositories');
    return this.http
      .get(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * GET /v2/project/{project_id}/companies
   */
  getProjectCompanies(projectId) {
    const url: URL = this.getV2Endpoint('/v2/project/' + projectId + '/companies');
    return this.http
      .get(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * GET /v2/project/{project_id}/document/{document_type}
   */
  getProjectDocument(projectId, documentType) {
    const url: URL = this.getV2Endpoint('/v2/project/' + projectId + '/document/' + documentType);
    return this.http
      .get(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  postProjectDocument(projectId, documentType, document) {
    /*
      {
        'document_name': 'doc_name.pdf',
        'document_content_type': 'url+pdf',
        'document_content': 'http://url.com/doc.pdf'
      }
     */
    const url: URL = this.getV1Endpoint('/v1/project/' + projectId + '/document/' + documentType);
    return this.http
      .post(url, document)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  postProjectDocumentTemplate(projectId, documentType, document) {
    /*
      {
        'document_name': 'project-name_ccla_2017-11-16',
        'document_preamble': '<p>Some <strong>html</strong> content</p>',
        'document_legal_entity_name': 'Some Project Inc.',
        'new_major_version': true|false,
      }
     */
    const url: URL = this.getV1Endpoint('/v1/project/' + projectId + '/document/template/' + documentType);
    return this.http
      .post(url, document)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * DELETE /project/{project_id}/document/{document_type}/{major_version}/{minor_version}
   */
  deleteProjectDocumentRevision(projectId, documentType, majorVersion, minorVersion) {
    const url: URL = this.getV1Endpoint(
      '/v1/project/' + projectId + '/document/' + documentType + '/' + majorVersion + '/' + minorVersion
    );
    return this.http
      .delete(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /*
   * GET /project/{project_id}/document/{document_type}/pdf/{document_major_version}/{document_minor_version}
   */
  getProjectDocumentRevisionPdf(projectId, documentType, majorVersion, minorVersion) {
    const url: URL = this.getV1Endpoint(
      '/v1/project/' + projectId + '/document/' + documentType + '/pdf/' + majorVersion + '/' + minorVersion
    );
    return this.http
      .get(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * POST /v2/request-individual-signature
   */
  postIndividualSignatureRequest(signatureRequest) {
    /*
      {
        'project_id': 'some-project-id',
        'user_id': 'some-user-uuid',
        'return_url': 'https://github.com/linuxfoundation/cla',
        'callback_url': 'http://cla.system/signed-callback'
      }
     */
    const url: URL = this.getV2Endpoint('/v2/request-individual-signature');
    return this.http
      .post(url, signatureRequest)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * POST /v2/request-employee-signature
   */
  postEmployeeSignatureRequest(signatureRequest) {
    /*
      {
        'project_id': <project-id>,
        'company_id': <company-id>,
        'user_id': <user-id>
      }
     */
    const url: URL = this.getV2Endpoint('/v2/request-employee-signature');
    return this.http
      .post(url, signatureRequest)
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
   * /signed/{installation_id}/{github_repository_id}/{change_request_id}
   */
  postSigned(installationId, githubRepositoryId, changeRequestId) {
    const url: URL = this.getV1Endpoint(
      '/v1/signed/' + installationId + '/' + githubRepositoryId + '/' + changeRequestId
    );
    return this.http
      .post(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * GET /v2/return-url/{signature_id}
   */
  getReturnUrl(signatureId) {
    const url: URL = this.getV2Endpoint('/v2/return-url/' + signatureId);
    return this.http
      .get(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * GET /v2/repository-provider/{provider}/sign/{installation_id}/{github_repository_id}/{change_request_id}
   */
  getSignRequest(provider, installationId, githubRepositoryId, changeRequestId) {
    const url: URL = this.getV2Endpoint(
      '/v2/repository-provider/' +
        provider +
        '/sign/' +
        installationId +
        '/' +
        githubRepositoryId +
        '/' +
        changeRequestId
    );
    return this.http
      .get(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * /send-authority-email
   */
  postEmailToCompanyAuthority(data) {
    /*
      {
        'project_name': 'Project Name',
        'company_name': 'Company Name',
        'authority_name': 'John Doe'
        'authority_email': 'johndoe@example.com'
      }
     */
    const url: URL = this.getV2Endpoint('/v2/send-authority-email');
    return this.http
      .post(url, data)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * GET /v2/repository-provider/{provider}/icon.svg
   */
  getChangeIcon(provider) {
    // This probably won't map to json, but instead to svg/xml
    const url: URL = this.getV2Endpoint('/v2/repository-provider/' + provider + '/icon.svg');
    return this.http
      .get(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * POST /v2/repository-provider/{provider}/activity
   */
  postReceivedActivity(provider) {
    const url: URL = this.getV2Endpoint('/v2/repository-provider/' + provider + '/activity');
    return this.http
      .post(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * GET /v1/github/organizations
   */
  getGithubOrganizations() {
    const url: URL = this.getV1Endpoint('/v1/github/organizations');
    return this.http
      .get(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * POST /v1/github/organizations
   */
  postGithubOrganization(organization) {
    /*
      organization: {
        'organization_project_id': '<project-id>',
        'organization_name': 'org-name'
      }
     */
    const url: URL = this.getV1Endpoint('/v1/github/organizations');
    return this.http
      .post(url, organization)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * GET /v1/github/get/namespace/{namespace}
   */
  getGithubGetNamespace(namespace) {
    const url: URL = this.getV1Endpoint('/v1/github/get/namespace/' + namespace);
    return this.http
      .get(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * GET /v1/github/check/namespace/{namespace}
   */
  getGithubCheckNamespace(namespace) {
    const url: URL = this.getV1Endpoint('/v1/github/check/namespace/' + namespace);
    return this.http
      .get(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * GET /v1/github/organizations/{organization_name}
   */
  getGithubOrganization(organizationName) {
    const url: URL = this.getV1Endpoint('/v1/github/organizations/' + organizationName);
    return this.http
      .get(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * DELETE /v1/github/organizations/{organization_name}
   */
  deleteGithubOrganization(organizationName) {
    const url: URL = this.getV1Endpoint('/v1/github/organizations/' + organizationName);
    return this.http
      .delete(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * GET /v1/github/organizations/{organization_name}/repositories
   */
  getGithubOrganizationRepositories(organizationName) {
    const url: URL = this.getV1Endpoint('/v1/github/organizations/' + organizationName + '/repositories');
    return this.http
      .get(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * GET /v2/github/installation
   */
  getGithubInstallation() {
    const url: URL = this.getV2Endpoint('/v2/github/installation');
    return this.http
      .get(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * POST /v2/github/installation
   */
  postGithubInstallation() {
    const url: URL = this.getV2Endpoint('/v2/github/installation');
    return this.http
      .post(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * POST /v2/github/activity
   */
  postGithubActivity() {
    const url: URL = this.getV2Endpoint('/v2/github/activity');
    return this.http
      .post(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * POST /v1/github/validate
   */
  postGithubValidate() {
    const url: URL = this.getV1Endpoint('/v1/github/validate');
    return this.http
      .post(url)
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
   * /company/{companyID}/cla/{corporateClaID}/whitelist/githuborg
   **/
  getGithubOrganizationWhitelist(signatureID, companyID, corporateClaID) {
    const url: URL = this.getV3Endpoint(`/v3/company/${companyID}/cla/${corporateClaID}/whitelist/githuborg`);
    return this.http.getWithCreds(url).map((res) => res.json());
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
    const url: URL = this.getV3Endpoint(`/v3/company/${companyId}/invite-request`);
    const data = {
      userID: userId,
      lfUsername: userId,
      lfEmail: userEmail,
      username: userName
    };
    console.log('Invoking sendInviteRequestEmail() with url: ' + url);
    return this.http
      .postWithCreds(url, data)
      .map((res) => res.json())
      .catch(this.handleServiceError);

    // Observable
    //.catchError((error: HttpErrorResponse) => {

    //  console.log('Caught HttpErrorResponse object:');

    //  console.log(error);
    //if (error.status === 500 && error.)

    //});
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
   **/
  getPendingInvites(companyId) {
    const url: URL = this.getV3Endpoint(`/v3/company/${companyId}/cla/invitelist`);
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

  acceptCompanyInvite(companyId, data) {
    const url: URL = this.getV3Endpoint(`/v3/company/${companyId}/cla/accesslist`);
    return this.http.postWithCreds(url, data).map((res) => res.json());
  }

  declineCompanyInvite(companyId, data) {
    const url: URL = this.getV3Endpoint(`/v3/company/${companyId}/cla/accesslist`);
    return this.http.deleteWithCredsAndBody(url, data).map((res) => res.json());
  }

  /**
   * Creates a new CLA Manager Request with the specified parameters.
   * @param lfid the LF ID of the user
   * @param projectName the project name
   * @param companyName the company name
   * @param userFullName the user full name
   * @param userEmail the user email
   */
  createCLAManagerRequest(lfid, projectName: string, companyName: string, userFullName: string, userEmail: string) {
    console.log('I got called');
    const url: URL = this.getV3Endpoint(`/v3/onboard/cla-manager`);
    const requestBody = {
      lf_id: lfid,
      project_name: projectName,
      company_name: companyName,
      user_full_name: userFullName,
      user_email: userEmail
    };
    //return this.http.post(url, requestBody)
    return this.http.postWithCreds(url, requestBody);
  }

  /**
   * Returns zero or more CLA manager requests based on the user LFID.
   *
   * @param lfid the user's LF ID
   */
  getCLAManagerRequests(lfid: string) {
    const url: URL = this.getV3Endpoint(`/v3/onboard/cla-manager/lfid/${lfid}`);
    return this.http.getWithCreds(url).map((res) => res.json());
  }

  /**
   * Deletes the CLA manager request by request id
   * @param requestID the unique request id
   */
  deleteCLAManagerRequests(requestID: string) {
    const url: URL = this.getV3Endpoint(`/v3/onboard/cla-manager/requests/${requestID}`);
    return this.http.delete(url);
  }

  /**
   * Sends an notification to the specified recipients.
   *
   * @param sender the sender of the message
   * @param subject the subject of the message
   * @param recipients the list of recipients
   * @param messageBody the message body
   */
  sendNotification(sender: string, subject: string, recipients: string[], messageBody: string) {
    const url: URL = this.getV3Endpoint(`/v3/onboard/notification`);
    const payload = {
      "sender_email": sender,
      "subject": subject,
      "recipient_emails": recipients,
      "email_body": messageBody,
    };
    return this.http.post(url, payload)
  }

  /**
   * Handle service error is a common routine to handle HTTP response errors
   * @param error the error
   */
  private handleServiceError(error: any) {
    /*
    console.log('Service error:');
    let errMsg: string;
    if (error instanceof Response) {
      console.log('HTTP Response service error:');
      const body = error.json() || '';
      const err = body.error || JSON.stringify(body);
      errMsg = `${error.status} - ${error.statusText || ''} ${err}`;
    } else {
      console.log('any service error:');
      const errString = String(error);
      console.log("Error:");
      console.log(errString);
      if (errString.includes('401')) {
        console.log('authentication error invoking service: ' + error + '. Forcing user to log out...');
        this.authService.logout();
      } else if (errString.includes('Token is expired')) {
        console.log(errString + ' - redirecting to login page');
        this.authService.logout();
      }
      errMsg = error.message ? error.message : error.toString();
    }
    console.log(errMsg);
    return Observable.throw(errMsg);
     */

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

  //////////////////////////////////////////////////////////////////////////////
}
