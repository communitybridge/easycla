// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Injectable } from "@angular/core";
import { Http, RequestOptions, URLSearchParams } from '@angular/http';
import { Observable } from "rxjs/Observable";

import "rxjs/Rx";

@Injectable()
export class ClaService {
  http: any;
  claApiUrl: string = "";
  s3LogoUrl: string = "";
  localTesting = false;
  v1ClaAPIURLLocal = 'http://localhost:5000';
  v2ClaAPIURLLocal = 'http://localhost:5000';
  v3ClaAPIURLLocal = 'http://localhost:8080';

  constructor(http: Http) {
    this.http = http;
  }

  public isLocalTesting(flag: boolean) {
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
   **/

  //////////////////////////////////////////////////////////////////////////////

  /**
   * /user
   */
  getUsers() {
    if (this.localTesting) {
      return this.http.get(this.v1ClaAPIURLLocal + "/v1/user").map(res => res.json());
    } else {
      return this.http.get(this.claApiUrl + "/v1/user").map(res => res.json());
    }
  }

  postUser(user) {
    /*
      {
        'user_email': 'user@email.com',
        'user_name': 'User Name',
        'user_company_id': '<org-id>',
        'user_github_id': 12345
      }
     */
    if (this.localTesting) {
      return this.http
        .post(this.v1ClaAPIURLLocal + "/v1/user", user)
        .map(res => res.json());
    } else {
      return this.http
        .post(this.claApiUrl + "/v1/user", user)
        .map(res => res.json());
    }
  }

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
    if (this.localTesting) {
      return this.http.put(this.v1ClaAPIURLLocal + "/v1/user", user).map(res => res.json());
    } else {
      return this.http.put(this.claApiUrl + "/v1/user", user).map(res => res.json());
    }
  }

  /**
   * /user/{user_id}
   */
  getUser(userId) {
    if (this.localTesting) {
      return this.http
        .get(this.v2ClaAPIURLLocal + "/v2/user/" + userId)
        .map(res => res.json());
    } else {
      return this.http
        .get(this.claApiUrl + "/v2/user/" + userId)
        .map(res => res.json());
    }
  }

  deleteUser(userId) {
    if (this.localTesting) {
      return this.http
        .delete(this.v1ClaAPIURLLocal + "/v1/user/" + userId)
        .map(res => res.json());
    } else {
      return this.http
        .delete(this.claApiUrl + "/v1/user/" + userId)
        .map(res => res.json());
    }
  }

  /**
   * /user/email/{user_email}
   */
  getUserByEmail(userEmail) {
    if (this.localTesting) {
      return this.http
        .get(this.v1ClaAPIURLLocal + "/v1/user/email/" + userEmail)
        .map(res => res.json());
    } else {
      return this.http
        .get(this.claApiUrl + "/v1/user/email/" + userEmail)
        .map(res => res.json());
    }
  }

  /**
   * /user/github/{user_github_id}
   */
  getUserByGithubId(userGithubId) {
    if (this.localTesting) {
      return this.http
        .get(this.v1ClaAPIURLLocal + "/v1/user/github/" + userGithubId)
        .map(res => res.json());
    } else {
      return this.http
        .get(this.claApiUrl + "/v1/user/github/" + userGithubId)
        .map(res => res.json());
    }
  }

  /**
   * /user/{user_id}/signatures
   */
  getUserSignatures(userId) {
    if (this.localTesting) {
      return this.http
        .get(this.v1ClaAPIURLLocal + "/v1/user/" + userId + "/signatures")
        .map(res => res.json());
    } else {
      return this.http
        .get(this.claApiUrl + "/v1/user/" + userId + "/signatures")
        .map(res => res.json());
    }
  }

  /**
   * /users/company/{user_company_id}
   */
  getUsersByCompanyId(userCompanyId) {
    if (this.localTesting) {
      return this.http
        .get(this.v1ClaAPIURLLocal + "/v1/users/company/" + userCompanyId)
        .map(res => res.json());
    } else {
      return this.http
        .get(this.claApiUrl + "/v1/users/company/" + userCompanyId)
        .map(res => res.json());
    }
  }

  /**
   * /user/{user_id}/request-company-whitelist/{company_id}
   */
  postUserMessageToCompanyManager(userId, companyId, message) {
    /*
      message: {
        'message': 'custom message to manager'
      }
      */
    if (this.localTesting) {
      return this.http
        .post(this.v2ClaAPIURLLocal + "/v2/user/" + userId + "/request-company-whitelist/" + companyId, message)
        .map(res => res.json());
    } else {
      return this.http
        .post(this.claApiUrl + "/v2/user/" + userId + "/request-company-whitelist/" + companyId, message)
        .map(res => res.json());
    }
  }

  /**
   * /user/{user_id}/active-signature
   */
  getUserSignatureIntent(userId) {
    if (this.localTesting) {
      return this.http
        .get(this.v2ClaAPIURLLocal + "/v2/user/" + userId + "/active-signature")
        .map(res => res.json());
    } else {
      return this.http
        .get(this.claApiUrl + "/v2/user/" + userId + "/active-signature")
        .map(res => res.json());
    }
  }

  /**
   * /user/{user_id}/project/{project_id}/last-signature
   */
  getLastIndividualSignature(userId, projectId) {
    if (this.localTesting) {
      return this.http
        .get(this.v2ClaAPIURLLocal + "/v2/user/" + userId + "/project/" + projectId + "/last-signature")
        .map(res => res.json());
    } else {
      return this.http
        .get(this.claApiUrl + "/v2/user/" + userId + "/project/" + projectId + "/last-signature")
        .map(res => res.json());
    }
  }

  /**
   * /signature
   */
  getSignatures() {
    if (this.localTesting) {
      return this.http.get(this.v1ClaAPIURLLocal + "/v1/signature").map(res => res.json());
    } else {
      return this.http.get(this.claApiUrl + "/v1/signature").map(res => res.json());
    }
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
    return this.http
      .post(this.claApiUrl + "/v1/signature", signature)
      .map(res => res.json());
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
    return this.http
      .put(this.claApiUrl + "/v1/signature", signature)
      .map(res => res.json());
  }

  /**
   * /signature/{signature_id}
   **/

  getSignature(signatureId) {
    return this.http
      .get(this.claApiUrl + "/v1/signature/" + signatureId)
      .map(res => res.json());
  }

  deleteSignature(signatureId) {
    return this.http
      .delete(this.claApiUrl + "/v1/signature/" + signatureId)
      .map(res => res.json());
  }

  /**
   * /signatures/user/{user_id}
   **/

  getSignaturesUser(userId) {
    return this.http
      .get(this.claApiUrl + "/v1/signatures/user/" + userId)
      .map(res => res.json());
  }

  /**
   * /signatures/company/{company_id}
   **/

  getCompanySignatures(companyId) {
    return this.http
      .get(this.claApiUrl + "/v1/signatures/company/" + companyId)
      .map(res => res.json());
  }

  /**
   * /signatures/company/{company_id}/project/{project_id}
   **/

  getCompanyProjectSignatures(companyId, projectId) {
    return this.http
      .get(
        this.claApiUrl +
        "/v1/signatures/company/" +
        companyId +
        "/project/" +
        projectId
      )
      .map(res => res.json());
  }

  /**
   * /signatures/project/{project_id}
   **/

  getProjectSignatures(projectId) {
    return this.http
      .get(this.claApiUrl + "/v1/signatures/project/" + projectId)
      .map(res => res.json());
  }

  /**
   * /repository
   **/

  getRepositories() {
    return this.http.get(this.claApiUrl + "/v1/repository").map(res => res.json());
  }

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
    return this.http
      .post(this.claApiUrl + "/v1/repository", repository)
      .map(res => res.json());
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
    return this.http
      .put(this.claApiUrl + "/v1/repository", repository)
      .map(res => res.json());
  }

  /**
   * /repository/{repository_id}
   **/

  getRepository(repositoryId) {
    return this.http
      .get(this.claApiUrl + "/v1/repository/" + repositoryId)
      .map(res => res.json());
  }

  deleteRepository(repositoryId) {
    return this.http
      .delete(this.claApiUrl + "/v1/repository/" + repositoryId)
      .map(res => res.json());
  }

  /**
   * /company
   **/


  // Returns list of companies for current user
  getCompanies() {
    return this.http.get(this.claApiUrl + "/v1/company").map(res => res.json());
  }

  getAllCompanies() {
    return this.http.get(this.claApiUrl + "/v2/company").map(res => res.json());
  }

  postCompany(company) {
    /*
      {
        'company_name': 'Org Name',
        'company_whitelist': ['safe@email.org'],
        'company_whitelist': ['*@email.org']
      }
     */
    return this.http
      .post(this.claApiUrl + "/v1/company", company)
      .map(res => res.json());
  }

  putCompany(company) {
    /*
      {
        'company_id': '<company-id>',
        'company_name': 'New Company Name'
      }
     */
    return this.http
      .put(this.claApiUrl + "/v1/company", company)
      .map(res => res.json());
  }

  /**
   * /company/{company_id}
   **/

  getCompany(companyId) {
    return this.http
      .get(this.claApiUrl + "/v2/company/" + companyId)
      .map(res => res.json());
  }

  deleteCompany(companyId) {
    return this.http
      .delete(this.claApiUrl + "/v1/company/" + companyId)
      .map(res => res.json());
  }

  /**
   * /project
   **/

  getProjects() {
    return this.http.get(this.claApiUrl + "/v1/project").map(res => res.json());
  }

  getProjectsCcla() {
    return this.http
      .get(this.claApiUrl + "/v1/project/ccla")
      .map(res => res.json());
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
    return this.http
      .post(this.claApiUrl + "/v1/project", project)
      .map(res => res.json());
  }

  putProject(project) {
    /*
      {
        'project_id': '<project-id>',
        'project_name': 'New Project Name'
      }
     */
    return this.http
      .put(this.claApiUrl + "/v1/project", project)
      .map(res => res.json());
  }

  /**
   * /project/{project_id}
   **/

  getProject(projectId) {
    return this.http
      .get(this.claApiUrl + "/v2/project/" + projectId)
      .map(res => res.json());
  }

  getProjectsByExternalId(externalId) {
    return this.http
      .get(this.claApiUrl + "/v1/project/external/" + externalId)
      .map(res => res.json());
  }

  deleteProject(projectId) {
    return this.http
      .delete(this.claApiUrl + "/v1/project/" + projectId)
      .map(res => res.json());
  }

  /**
   * /project/{project_id}/repositories
   **/

  getProjectRepositories(projectId) {
    return this.http
      .get(this.claApiUrl + "/v1/project/" + projectId + "/repositories")
      .map(res => res.json());
  }

  /**
   * /project/{project_id}/repositories_by_org
   **/

  getProjectRepositoriesByrOrg(projectId) {
    return this.http
      .get(this.claApiUrl + "/v1/project/" + projectId + "/repositories_group_by_organization")
      .map(res => res.json());
  }

  /**
   * /repository
   **/

  postProjectRepository(repository) {
    return this.http
      .post(this.claApiUrl + "/v1/repository/", repository)
      .map(res => res.json());
  }

  /**
   * /repository
   **/

  removeProjectRepository(repositoryId) {
    return this.http
      .delete(this.claApiUrl + "/v1/repository/" + repositoryId)
      .map(res => res.json());
  }

  /**
   * /project/{project_id}/configuration_orgs_and_repos
   **/

  getProjectConfigurationAndRepos(projectId) {
    return this.http
      .get(this.claApiUrl + "/v1/project/" + projectId + "/configuration_orgs_and_repos")
      .map(res => res.json());
  }

  /**
   * /sfdc/${sfid}/github/organizations
   **/

  getOrganizations(sfid) {
    return this.http
      .get(this.claApiUrl + `/v1/sfdc/${sfid}/github/organizations`)
      .map(res => res.json());
  }

  /**
   * /project/{project_id}/companies
   **/

  getProjectCompanies(projectId) {
    return this.http
      .get(this.claApiUrl + "/v2/project/" + projectId + "/companies")
      .map(res => res.json());
  }

  /**
   * /project/{project_id}/document/{document_type}
   **/

  getProjectDocument(projectId, documentType) {
    return this.http
      .get(
        this.claApiUrl + "/v2/project/" + projectId + "/document/" + documentType
      )
      .map(res => res.json());
  }

  postProjectDocument(projectId, documentType, document) {
    /*
      {
        'document_name': 'doc_name.pdf',
        'document_content_type': 'url+pdf',
        'document_content': 'http://url.com/doc.pdf'
      }
     */
    return this.http
      .post(
        this.claApiUrl + "/v1/project/" + projectId + "/document/" + documentType,
        document
      )
      .map(res => res.json());
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
    return this.http
      .post(
        this.claApiUrl +
        "/v1/project/" +
        projectId +
        "/document/template/" +
        documentType,
        document
      )
      .map(res => res.json());
  }

  /**
   * /project/{project_id}/document/{document_type}/{major_version}/{minor_version}
   */
  deleteProjectDocumentRevision(
    projectId,
    documentType,
    majorVersion,
    minorVersion
  ) {
    return this.http
      .delete(
        this.claApiUrl +
        "/v1/project/" +
        projectId +
        "/document/" +
        documentType +
        "/" +
        majorVersion +
        "/" +
        minorVersion
      )
      .map(res => res.json());
  }

  /*
   * /project/{project_id}/document/{document_type}/pdf/{document_major_version}/{document_minor_version}
   */
  getProjectDocumentRevisionPdf(
    projectId,
    documentType,
    majorVersion,
    minorVersion
  ) {
    return this.http
    .get(this.claApiUrl +
      "/v1/project/" +
      projectId +
      "/document/" +
      documentType +
      "/pdf/" +
      majorVersion +
      "/" +
      minorVersion
    )
    .map(res => res.text());
  }

  /**
   * /request-individual-signature
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
    return this.http
      .post(this.claApiUrl + "/v2/request-individual-signature", signatureRequest)
      .map(res => res.json());
  }

  /**
   * /request-employee-signature
   */
  postEmployeeSignatureRequest(signatureRequest) {
    /*
      {
        'project_id': <project-id>,
        'company_id': <company-id>,
        'user_id': <user-id>
      }
     */
    return this.http
      .post(this.claApiUrl + "/v2/request-employee-signature", signatureRequest)
      .map(res => res.json());
  }

  /**
   * /request-corporate-signature
   */
  postCorporateSignatureRequest(signatureRequest) {
    /*
      {
        'project_id': <project-id>,
        'company_id': <company-id>,
        'return_url': <optional-return-url>,
      }
     */
    return this.http
      .post(this.claApiUrl + "/v1/request-corporate-signature", signatureRequest)
      .map(res => res.json());
  }

  /**
   * /signed/{installation_id}/{github_repository_id}/{change_request_id}
   */
  postSigned(installationId, githubRepositoryId, changeRequestId) {
    return this.http
      .post(
        this.claApiUrl +
        "/v1/signed/" +
        installationId +
        "/" +
        githubRepositoryId +
        "/" +
        changeRequestId
      )
      .map(res => res.json());
  }

  /**
   * /return-url/{signature_id}
   */
  getReturnUrl(signatureId) {
    return this.http
      .get(this.claApiUrl + "/v2/return-url/" + signatureId)
      .map(res => res.json());
  }

  /**
   * /repository-provider/{provider}/sign/{installation_id}/{github_repository_id}/{change_request_id}
   */
  getSignRequest(
    provider,
    installationId,
    githubRepositoryId,
    changeRequestId
  ) {
    return this.http
      .get(
        this.claApiUrl +
        "/v2/repository-provider/" +
        provider +
        "/sign/" +
        installationId +
        "/" +
        githubRepositoryId +
        "/" +
        changeRequestId
      )
      .map(res => res.json());
  }

  /**
   * /repository-provider/{provider}/icon.svg
   */
  getChangeIcon(provider) {
    // This probably won't map to json, but instead to svg/xml
    return this.http
      .get(this.claApiUrl + "/v2/repository-provider/" + provider + "/icon.svg")
      .map(res => res.json());
  }

  /**
   * /repository-provider/{provider}/activity
   */
  postReceivedActivity(provider) {
    return this.http
      .post(this.claApiUrl + "/v2/repository-provider/" + provider + "/activity")
      .map(res => res.json());
  }

  /**
   * /github/organizations
   */
  getGithubOrganizations() {
    return this.http
      .get(this.claApiUrl + "/v1/github/organizations")
      .map(res => res.json());
  }

  postGithubOrganization(organization) {
    /*
      organization: {
        'organization_project_id': '<project-id>',
        'organization_name': 'org-name'
      }
     */
    return this.http
      .post(this.claApiUrl + "/v1/github/organizations", organization)
      .map(res => res.json());
  }

  /**
   * /github/get/namespace/{namespace}
   */
  getGithubGetNamespace(namespace) {
    return this.http
      .get(this.claApiUrl + "/v1/github/get/namespace/" + namespace)
      .map(res => res.json());
  }

  /**
   * /github/check/namespace/{namespace}
   */
  getGithubCheckNamespace(namespace) {
    return this.http
      .get(this.claApiUrl + "/v1/github/check/namespace/" + namespace)
      .map(res => res.json());
  }

  /**
   * /github/organizations/{organization_name}
   */
  getGithubOrganization(organizationName) {
    return this.http
      .get(this.claApiUrl + "/v1/github/organizations/" + organizationName)
      .map(res => res.json());
  }

  deleteGithubOrganization(organizationName) {
    return this.http
      .delete(this.claApiUrl + "/v1/github/organizations/" + organizationName)
      .map(res => res.json());
  }

  /**
   * /github/organizations/{organization_name}/repositories
   */
  getGithubOrganizationRepositories(organizationName) {
    return this.http
      .get(
        this.claApiUrl +
        "/v1/github/organizations/" +
        organizationName +
        "/repositories"
      )
      .map(res => res.json());
  }

  /**
   * /github/installation
   **/

  getGithubInstallation() {
    return this.http
      .get(this.claApiUrl + "/v2/github/installation")
      .map(res => res.json());
  }

  postGithubInstallation() {
    return this.http
      .post(this.claApiUrl + "/v2/github/installation")
      .map(res => res.json());
  }

  /**
   * /github/activity
   **/

  postGithubActivity() {
    return this.http
      .post(this.claApiUrl + "/v2/github/activity")
      .map(res => res.json());
  }

  /**
   * /github/validate
   **/

  postGithubValidate() {
    return this.http
      .post(this.claApiUrl + "/v1/github/validate")
      .map(res => res.json());
  }

  /**
   * /salesforce/projects
   **/

  getAllProjectsFromSFDC() {
    return this.http
      .get(this.claApiUrl + "/v1/salesforce/projects")
      .map(res => res.json()
        .map(p => this.addProjectLogoFromS3(p))
      );
  }

  getProjectFromSFDC(projectId) {
    return this.http
      .get(this.claApiUrl + `/v1/salesforce/project?id=${projectId}`)
      .map(res => res.json())
      .map(p => this.addProjectLogoFromS3(p));
  }

  addProjectLogoFromS3(project) {
    let objLogoUrl = {
      "logoRef": `${this.s3LogoUrl}/${project.id}.png`
    }
    return { ...project, ...objLogoUrl }
  }


  getGerritInstance(projectId) {
    return this.http
      .get(this.claApiUrl + `/v1/project/${projectId}/gerrits`)
      .map(res => res.json())
  }

  deleteGerritInstance(gerritId) {
    return this.http
      .delete(this.claApiUrl + `/v1/gerrit/${gerritId}`)
      .map(res => res.json())
  }

  postGerritInstance(gerrit) {
    return this.http
      .post(this.claApiUrl + "/v1/gerrit", gerrit)
      .map(res => res.json());
  }

  getTemplates() {
    if (this.localTesting) {
      return this.http
        .get(this.v3ClaAPIURLLocal + `/v3/template`)
        .map(res => res.json())
    } else {
      return this.http
        .get(this.claApiUrl + `/v3/template`)
        .map(res => res.json())
    }
  }

  postClaGroupTemplate(projectId, data) {
    if (this.localTesting) {
      return this.http
        .post(this.v3ClaAPIURLLocal + `/v3/clagroup/${projectId}/template`, data)
        .map(res => res.json());
    } else {
      return this.http
        .post(this.claApiUrl + `/v3/clagroup/${projectId}/template`, data)
        .map(res => res.json());
    }
  }
  //////////////////////////////////////////////////////////////////////////////
}
