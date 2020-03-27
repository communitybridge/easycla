// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Injectable } from '@angular/core';
import { Http, RequestOptions } from '@angular/http';
import { AuthService } from './auth.service';

import 'rxjs/Rx';
import { Observable } from 'rxjs/Rx';
import {Response} from "@angular/http/src/static_response";

@Injectable()
export class ClaService {
  http: any;
  authService: AuthService;
  claApiUrl: string = '';
  s3LogoUrl: string = '';
  localTesting = false;
  v1ClaAPIURLLocal = 'http://localhost:5000';
  v2ClaAPIURLLocal = 'http://localhost:5000';
  v3ClaAPIURLLocal = 'http://localhost:8080';

  constructor(http: Http, authService: AuthService) {
    this.http = http;
    this.authService = authService;
  }

  private getV1APIEndpoint(path: string) {
    return new URL(this.claApiUrl + path);
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

  public setS3LogoUrl(s3LogoUrl: string) {
    this.s3LogoUrl = s3LogoUrl;
  }

  public setHttp(http: any) {
    this.http = http; // allow configuration for alternate http library
  }

  //////////////////////////////////////////////////////////////////////////////

  /**
   * This service should ONLY contain methods calling CLA API
   */

  //////////////////////////////////////////////////////////////////////////////

  /**
   * GET /user
   */
  getUsers() {
    const url: URL = this.getV1Endpoint('/v1/user');
    return this.http.get(url).map((res) => res.json());
  }

  /**
   * POST /v1/user
   * @param user the user payload
   */
  postUser(user) {
    /*
      {
        'user_email': 'user@email.com',
        'user_name': 'User Name',
        'user_company_id': '<org-id>',
        'user_github_id': 12345
      }
     */
    const url: URL = this.getV1Endpoint('/v1/user');
    return this.http.post(url, user).map((res) => res.json());
  }

  /**
   * PUT /v1/user
   * @param user the user payload
   */
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
    return this.http.put(url, user).map((res) => res.json());
  }

  /**
   * GET /user/{user_id}
   */
  getUser(userId) {
    const url: URL = this.getV2Endpoint('/v2/user/' + userId);
    return this.http.get(url).map((res) => res.json());
  }

  /**
   * DELETE /v1/user/userId
   * @param userId the user ID
   */
  deleteUser(userId) {
    const url: URL = this.getV1Endpoint('/v1/user/' + userId);
    return this.http.delete(url).map((res) => res.json());
  }

  /**
   * GET /v1/user/email/{user_email}
   */
  getUserByEmail(userEmail) {
    const url: URL = this.getV1Endpoint('/v1/user/email/' + userEmail);
    return this.http.get(url).map((res) => res.json());
  }

  /**
   * GET /user/github/{user_github_id}
   */
  getUserByGithubId(userGithubId) {
    const url: URL = this.getV1Endpoint('/v1/user/github/' + userGithubId);
    return this.http.get(url).map((res) => res.json());
  }

  /**
   * GET /user/{user_id}/signatures
   */
  getUserSignatures(userId) {
    const url: URL = this.getV1Endpoint('/v1/user/' + userId + '/signature');
    return this.http.get(url).map((res) => res.json());
  }

  /**
   * GET /users/company/{user_company_id}
   */
  getUsersByCompanyId(userCompanyId) {
    const url: URL = this.getV1Endpoint('/v1/users/company/' + userCompanyId);
    return this.http.get(url).map((res) => res.json());
  }

  /**
   * POST /user/{user_id}/request-company-whitelist/{company_id}
   */
  postUserMessageToCompanyManager(userId, companyId, message) {
    /*
      message: {
        'message': 'custom message to manager'
      }
      */
    const url: URL = this.getV2Endpoint('/v2/user/' + userId + '/request-company-whitelist/' + companyId);
    return this.http.post(url, message).map((res) => res.json());
  }

  /**
   * GET /user/{user_id}/active-signature
   */
  getUserSignatureIntent(userId) {
    const url: URL = this.getV1Endpoint('/v2/user/' + userId + '/active-signature');
    return this.http.get(url).map((res) => res.json());
  }

  /**
   * GET /user/{user_id}/project/{project_id}/last-signature
   */
  getLastIndividualSignature(userId, projectId) {
    const url: URL = this.getV2Endpoint('/v2/user/' + userId + '/project/' + projectId + '/last-signature');
    return this.http.get(url).map((res) => res.json());
  }

  /**
   * GET /signature
   */
  getSignatures() {
    const url: URL = this.getV1Endpoint('/v1/signature');
    return this.http.get(url).map((res) => res.json());
  }

  /**
   * POST /signature
   */
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
    return this.http.post(url, signature).map((res) => res.json());
  }

  /**
   * PUT /v1/signature
   * @param signature the signature payload
   */
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
    return this.http.put(url, signature).map((res) => res.json());
  }

  /**
   * GET /v3/signature/{signature_id}
   *
   * @param signatureId the signature ID
   */
  getSignature(signatureId) {
    //const url: URL = this.getV1Endpoint('/v1/signature/' + signatureId);
    // Leverage the new go backend v3 endpoint
    const url: URL = this.getV3Endpoint('/v3/signature/' + signatureId);
    return this.http
      .get(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
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
   */
  getSignaturesUser(userId) {
    //const url: URL = this.getV1Endpoint('/v1/signatures/user/' + userId);
    // Leverage the new go backend v3 endpoint
    const url: URL = this.getV3Endpoint('/v3/signatures/user/' + userId);
    return this.http
      .get(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * GET /v3/signatures/company/{company_id}
   *
   * @param companyId the company ID
   */
  getCompanySignatures(companyId) {
    //const url: URL = this.getV1Endpoint('/v1/signatures/company/' + companyId);
    // Leverage the new go backend v3 endpoint
    const url: URL = this.getV3Endpoint('/v3/signatures/company/' + companyId);
    return this.http
      .get(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * GET /v3/signatures/project/{project_id}/company/{company_id}
   *
   * @param companyId the company ID
   * @param projectId the project ID
   */
  getCompanyProjectSignatures(companyId, projectId) {
    //const url: URL = this.getV1Endpoint('/v1/signatures/company/' + companyId + '/project/' + projectId);
    // Leverage the new go backend v3 endpoint - note the slightly different path layout
    const url: URL = this.getV3Endpoint('/v3/signatures/project/' + projectId + '/company/' + companyId);
    return this.http
      .get(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * GET /v3/signatures/project/{project_id}/company/{company_id}/employee
   *
   * @param companyId the company ID
   * @param projectId the project ID
   */
  getEmployeeProjectSignatures(companyId, projectId) {
    //const url: URL = this.getV1Endpoint('/v1/signatures/company/' + companyId + '/project/' + projectId + '/employee');
    // Leverage the new go backend v3 endpoint - note the different order of the parameters in the path
    const url: URL = this.getV3Endpoint('/v3/signatures/project/' + projectId + '/company/' + companyId + '/employee');
    return this.http
      .get(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * GET /v3/signatures/project/{project_id}
   *
   * @param projectId the project ID
   */
  getProjectSignatures(projectId, lastKeyScanned) {
    // Leverage the new go backend v3 endpoint - note the slightly different path
    if (lastKeyScanned) {
      const url: URL = this.getV3Endpoint(`/v3/signatures/project/${projectId}?pageSize=50&nextKey=${lastKeyScanned}`);
      return this.http
        .get(url)
        .map((res) => res.json())
        .catch((error) => this.handleServiceError(error));
    } else {
      const url: URL = this.getV3Endpoint(`/v3/signatures/project/${projectId}?pageSize=50`);
      return this.http
        .get(url)
        .map((res) => res.json())
        .catch((error) => this.handleServiceError(error));
    }
  }

  /**
   * GET /v3/signatures/{project_id} - v3 backend query which supports pagination.
   *
   * @param projectId the project id
   * @param pageSize the optional page size - default is 50
   * @param nextKey the next key used when asking for the next page of results
   */
  getProjectSignaturesV3(
    projectId,
    pageSize = 50,
    nextKey = '',
    searchTerm = '',
    searchField = '',
    signatureType = '',
    fullMatch = false,
  ) {
    let path: string = '/v3/signatures/project/' + projectId + '?pageSize=' + pageSize;
    if (nextKey !== null && nextKey !== '' && nextKey.trim().length > 0) {
      path += `&nextKey=${nextKey}`;
    }

    if (signatureType) {
      path += `&signatureType=${signatureType}`;
    }

    if (searchTerm !== null && searchTerm !== '' && searchTerm.trim().length > 0) {
      path += `&searchTerm=${searchTerm}&searchField=${searchField}`;

      if (fullMatch) {
        path += '&fullMatch=true';
      }
    }

    const url: URL = this.getV3Endpoint(path);
    return this.http.get(url).map((res) => res.json());
  }

  /**
   * GET /repository
   */
  getRepositories() {
    const url: URL = this.getV1Endpoint('/v1/repository');
    return this.http.get(url).map((res) => res.json());
  }

  /**
   * POST /repository
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
    return this.http.post(url, repository).map((res) => res.json());
  }

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
    return this.http.put(url, repository).map((res) => res.json());
  }

  /**
   * GET /repository/{repository_id}
   **/
  getRepository(repositoryId) {
    const url: URL = this.getV1Endpoint('/v1/repository/' + repositoryId);
    return this.http.get(url).map((res) => res.json());
  }

  /**
   * DELETE /repository/{repository_id}
   */
  deleteRepository(repositoryId) {
    const url: URL = this.getV1Endpoint('/v1/repository/' + repositoryId);
    return this.http.delete(url).map((res) => res.json());
  }

  /**
   * /company
   **/

  /**
   * Returns list of companies for current user
   * GET /v1/company
   */
  getCompanies() {
    const url: URL = this.getV1Endpoint('/v1/company');
    return this.http.get(url).map((res) => res.json());
  }

  /**
   * GET /v2/company/{company_id}
   */
  getAllCompanies() {
    const url: URL = this.getV2Endpoint('/v2/company');
    return this.http.get(url).map((res) => res.json());
  }

  /**
   * GET /v3/company/search?companyName={company_name}
   */
  searchCompaniesByName(companyName) {
    const url: URL = this.getV3Endpoint('/v3/company/search?companyName=' + companyName);
    return this.http
      .get(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  postCompany(company) {
    /*
      {
        'company_name': 'Org Name',
        'company_whitelist': ['safe@email.org'],
        'company_whitelist': ['*@email.org']
      }
     */
    const url: URL = this.getV1Endpoint('/v1/company');
    return this.http.post(url, company).map((res) => res.json());
  }

  putCompany(company) {
    /*
      {
        'company_id': '<company-id>',
        'company_name': 'New Company Name'
      }
     */
    const url: URL = this.getV1Endpoint('/v1/company');
    return this.http.put(url, company).map((res) => res.json());
  }

  /**
   * GET /company/{company_id}
   * @param companyId
   */
  getCompany(companyId) {
    const url: URL = this.getV1Endpoint('/v2/company/' + companyId);
    return this.http.get(url).map((res) => res.json());
  }

  /**
   * DELETE /company/{company_id}
   * @param companyId
   */
  deleteCompany(companyId) {
    const url: URL = this.getV1Endpoint('/v1/company/' + companyId);
    return this.http.delete(url).map((res) => res.json());
  }

  /**
   * GET /project
   */
  getProjects() {
    const url: URL = this.getV1Endpoint('/v1/project');
    return this.http.get(url).map((res) => res.json());
  }

  getProjectsCcla() {
    const url: URL = this.getV1Endpoint('/v1/project/ccla');
    return this.http.get(url).map((res) => res.json());
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
    const url: URL = this.getV3Endpoint('/v3/project');
    return this.http.post(url, project).map((res) => res.json());
  }

  /**
   * PUT /project/{project_id}
   * @param project the project payload
   */
  putProject(project) {
    /*
      {
        'project_id': '<project-id>',
        'project_name': 'New Project Name'
      }
     */
    const url: URL = this.getV3Endpoint('/v3/project');
    return this.http.put(url, project).map((res) => res.json());
  }

  /**
   * GET /project/{project_id}
   * @param projectId the project ID
   */
  getProject(projectId) {
    const url: URL = this.getV2Endpoint('/v2/project/' + projectId);
    return this.http.get(url).map((res) => res.json());
  }

  /**
   * GET /v1/project/external/{externalId}
   * @param externalId the external ID
   */
  getProjectsByExternalId(externalId) {
    const url: URL = this.getV3Endpoint(`/v3/project/external/${externalId}`);
    return this.http
      .get(url)
      .map((res) => res.json())
      .catch((error) => this.handleServiceError(error));
  }

  /**
   * DELETE /v1/project/{projectId}
   * @param projectId the project ID
   */
  deleteProject(projectId) {
    const url: URL = this.getV1Endpoint('/v1/project/' + projectId);
    return this.http.delete(url).map((res) => res.json());
  }

  /**
   * GET /project/{project_id}/repositories
   **/
  getProjectRepositories(projectId) {
    const url: URL = this.getV1Endpoint('/v1/project/' + projectId + '/repositories');
    return this.http.get(url).map((res) => res.json());
  }

  /**
   * GET /project/{project_id}/repositories_by_org
   **/
  getProjectRepositoriesByrOrg(projectId) {
    const url: URL = this.getV1Endpoint('/v1/project/' + projectId + '/repositories_group_by_organization');
    return this.http.get(url).map((res) => res.json());
  }

  /**
   * POST /repository
   **/
  postProjectRepository(repository) {
    const url: URL = this.getV1Endpoint('/v1/repository');
    return this.http.post(url, repository).map((res) => res.json());
  }

  /**
   * DELETE /repository
   */
  removeProjectRepository(repositoryId) {
    const url: URL = this.getV1Endpoint('/v1/repository/' + repositoryId);
    return this.http.delete(url).map((res) => res.json());
  }

  /**
   * GET /project/{project_id}/configuration_orgs_and_repos
   */
  getProjectConfigurationAndRepos(projectId) {
    const url: URL = this.getV1Endpoint('/v1/project/' + projectId + '/configuration_orgs_and_repos');
    return this.http.get(url).map((res) => res.json());
  }

  /**
   * GET /sfdc/${sfid}/github/organizations
   */
  getOrganizations(sfid) {
    const url: URL = this.getV1Endpoint('/v1/sfdc/' + sfid + '/github/organizations');
    return this.http.get(url).map((res) => res.json());
  }

  /**
   * GET /project/{project_id}/companies
   */
  getProjectCompanies(projectId) {
    const url: URL = this.getV2Endpoint('/v2/project/' + projectId + '/companies');
    return this.http.get(url).map((res) => res.json());
  }

  /**
   * GET /project/{project_id}/document/{document_type}
   */
  getProjectDocument(projectId: String, documentType: String) {
    const url: URL = this.getV2Endpoint('/v2/project/' + projectId + '/document/' + documentType);
    return this.http.get(url).map((res) => res.json());
  }

  postProjectDocument(projectId: String, documentType: String, document) {
    /*
      {
        'document_name': 'doc_name.pdf',
        'document_content_type': 'url+pdf',
        'document_content': 'http://url.com/doc.pdf'
      }
     */
    const url: URL = this.getV1Endpoint('/v1/project/' + projectId + '/document/' + documentType);
    return this.http.post(url, document).map((res) => res.json());
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
    return this.http.post(url, document).map((res) => res.json());
  }

  /**
   * DELETE /project/{project_id}/document/{document_type}/{major_version}/{minor_version}
   */
  deleteProjectDocumentRevision(projectId, documentType, majorVersion, minorVersion) {
    const url: URL = this.getV1Endpoint(
      '/v1/project/' + projectId + '/document/' + documentType + '/' + majorVersion + '/' + minorVersion,
    );
    return this.http.delete(url).map((res) => res.json());
  }

  /*
   * GET /project/{project_id}/document/{document_type}/pdf/{document_major_version}/{document_minor_version}
   */
  getProjectDocumentRevisionPdf(projectId, documentType, majorVersion, minorVersion) {
    const url: URL = this.getV1Endpoint(
      '/v1/project/' + projectId + '/document/' + documentType + '/pdf/' + majorVersion + '/' + minorVersion,
    );
    return this.http.get(url).map((res) => {
      return res._body;
    });
  }

  /**
   * POST /request-individual-signature
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
    return this.http.post(url, signatureRequest).map((res) => res.json());
  }

  /**
   * POST /request-employee-signature
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
    return this.http.post(url, signatureRequest).map((res) => res.json());
  }

  /**
   * POST /request-corporate-signature
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
    return this.http.post(url, signatureRequest).map((res) => res.json());
  }

  /**
   * POST /signed/{installation_id}/{github_repository_id}/{change_request_id}
   */
  postSigned(installationId, githubRepositoryId, changeRequestId) {
    const url: URL = this.getV1Endpoint(
      '/v1/signed/' + installationId + '/' + githubRepositoryId + '/' + changeRequestId,
    );
    return this.http.post(url).map((res) => res.json());
  }

  /**
   * GET /return-url/{signature_id}
   */
  getReturnUrl(signatureId) {
    const url: URL = this.getV2Endpoint('/v2/return-url/' + signatureId);
    return this.http.get(url).map((res) => res.json());
  }

  /**
   * GET /repository-provider/{provider}/sign/{installation_id}/{github_repository_id}/{change_request_id}
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
      changeRequestId,
    );
    return this.http.get(url).map((res) => res.json());
  }

  /**
   * GET /repository-provider/{provider}/icon.svg
   */
  getChangeIcon(provider) {
    // This probably won't map to json, but instead to svg/xml
    const url: URL = this.getV2Endpoint('/v2/repository-provider/' + provider + '/icon.svg');
    return this.http.get(url).map((res) => res.json());
  }

  /**
   * POST /repository-provider/{provider}/activity
   */
  postReceivedActivity(provider) {
    const url: URL = this.getV2Endpoint('/v2/repository-provider/' + provider + '/activity');
    return this.http.post(url).map((res) => res.json());
  }

  /**
   * GET /github/organizations
   */
  getGithubOrganizations() {
    const url: URL = this.getV1Endpoint('/v1/github/organizations');
    return this.http.get(url).map((res) => res.json());
  }

  /**
   * POST /github/organizations
   */
  postGithubOrganization(organization) {
    /*
      organization: {
        'organization_project_id': '<project-id>',
        'organization_name': 'org-name'
      }
     */
    const url: URL = this.getV1Endpoint('/v1/github/organizations');
    return this.http.post(url, organization).map((res) => res.json());
  }

  /**
   * GET /github/get/namespace/{namespace}
   */
  getGithubGetNamespace(namespace) {
    const url: URL = this.getV1Endpoint(`/v1/github/get/namespace/${namespace}` );
    return this.http.get(url).map((res) => res.json())
      .catch((error) => Observable.throw(error));
  }

  /**
   * GET /github/check/namespace/{namespace}
   */
  getGithubCheckNamespace(namespace) {
    const url: URL = this.getV1Endpoint('/v1/github/check/namespace/' + namespace);
    return this.http.get(url).map((res) => res.json());
  }

  /**
   * GET /github/organizations/{organization_name}
   */
  getGithubOrganization(organizationName) {
    const url: URL = this.getV1Endpoint('/v1/github/organizations/' + organizationName);
    return this.http.get(url).map((res) => res.json());
  }

  /**
   * DELETE /github/organizations/{organization_name}
   */
  deleteGithubOrganization(organizationName) {
    const url: URL = this.getV1Endpoint('/v1/github/organizations/' + organizationName);
    return this.http.delete(url).map((res) => res.json());
  }

  /**
   * GET /github/organizations/{organization_name}/repositories
   */
  getGithubOrganizationRepositories(organizationName) {
    const url: URL = this.getV1Endpoint('/v1/github/organizations/' + organizationName + '/repositories');
    return this.http.get(url).map((res) => res.json());
  }

  /**
   * GET /github/installation
   */
  getGithubInstallation() {
    const url: URL = this.getV2Endpoint('/v2/github/installation');
    return this.http.get(url).map((res) => res.json());
  }

  /**
   * POST /github/installation
   */
  postGithubInstallation() {
    const url: URL = this.getV2Endpoint('/v2/github/installation');
    return this.http.post(url).map((res) => res.json());
  }

  /**
   * POST /github/activity
   */
  postGithubActivity() {
    const url: URL = this.getV2Endpoint('/v2/github/activity');
    return this.http.post(url).map((res) => res.json());
  }

  /**
   * POST /github/validate
   */
  postGithubValidate() {
    const url: URL = this.getV1Endpoint('/v1/github/validate');
    return this.http.post(url).map((res) => res.json());
  }

  /**
   * GET /salesforce/projects
   */
  getAllProjectsFromSFDC() {
    // Use the deployed API endpoint regardless as the salesforce endpoints are on a different lambda and
    // we don't/can't run this locally
    const url: URL = this.getV1APIEndpoint('/v1/salesforce/projects');
    return this.http.get(url).map((res) => res.json().map((p) => this.addProjectLogoFromS3(p)));
  }

  /**
   * GET /salesforce/projects?id={projectId}
   */
  getProjectFromSFDC(projectId) {
    // Use the deployed API endpoint regardless as the salesforce endpoints are on a different lambda and
    // we don't/can't run this locally
    const url: URL = this.getV1APIEndpoint('/v1/salesforce/project?id=' + projectId);
    return this.http
      .get(url)
      .map((res) => res.json())
      .map((p) => this.addProjectLogoFromS3(p));
  }

  addProjectLogoFromS3(project) {
    let objLogoUrl = {
      logoRef: `${this.s3LogoUrl}/${project.id}.png`,
    };
    return { ...project, ...objLogoUrl };
  }

  getGerritInstance(projectId) {
    const url: URL = this.getV1Endpoint('/v1/project/' + projectId + '/gerrits');
    return this.http.get(url).map((res) => res.json());
  }

  deleteGerritInstance(gerritId) {
    const url: URL = this.getV1Endpoint('/v1/gerrit/' + gerritId);
    return this.http.delete(url).map((res) => res.json());
  }

  postGerritInstance(gerrit) {
    const url: URL = this.getV1Endpoint('/v1/gerrit');
    return this.http.post(url, gerrit).map((res) => res.json());
  }

  getTemplates() {
    const url: URL = this.getV3Endpoint('/v3/template');
    return this.http.get(url).map((res) => res.json());
  }

  postClaGroupTemplate(projectId, data) {
    const url: URL = this.getV3Endpoint('/v3/clagroup/' + projectId + '/template');
    return this.http.post(url, data).map((res) => res.json());
  }

  private handleServiceError(error: any) {
    const errString = String(error);
    if (errString.includes('401')) {
      console.log('authentication error invoking service: ' + error + '. Forcing user to log out...');
      this.authService.logout();
    } else {
      console.log('problem invoking service: ' + error);
    }
  }

  getReleaseVersion() {
    const url: URL = this.getV3Endpoint('/v3/ops/version');
    return this.http.get(url).map((res) => res.json());
  }

  /**
   * Check if the specified GitHub Organization name is valid
   * @param githubOrgName the GitHub Organization Name
   */
  testGitHubOrganization(githubOrgName: string): Observable<Response> {
    const url: URL = this.getV3Endpoint(`/v3/github/org/${githubOrgName}/exists`);
    return this.http.get(url);
  }
}
