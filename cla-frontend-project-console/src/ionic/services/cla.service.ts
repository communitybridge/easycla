// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Injectable } from '@angular/core';
import { Http } from '@angular/http';
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
   * GET /user/{user_id}
   */
  getUser(userId) {
    const url: URL = this.getV2Endpoint('/v2/user/' + userId);
    return this.http.get(url).map((res) => res.json());
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
   * /company
   **/


  /**
   * GET /company/{company_id}
   * @param companyId
   */
  getCompany(companyId) {
    const url: URL = this.getV1Endpoint('/v2/company/' + companyId);
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

  deleteClaProject(projectId: string): Observable<Response> {
    const url: URL = this.getV3Endpoint(`/v3/project/${projectId}`);
    return this.http.delete(url);
  }

  /**
   * GET /project/{project_id}/repositories_by_org
   **/
  // getProjectRepositoriesByrOrg(projectId) {
  //   const url: URL = this.getV1Endpoint('/v1/project/' + projectId + '/repositories_group_by_organization');
  //   return this.http.get(url).map((res) => res.json());
  // }

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
    // Outdated V1
    // const url: URL = this.getV1Endpoint('/v1/sfdc/' + sfid + '/github/organizations');
    const url: URL = this.getV3Endpoint(`/v3/project/${sfid}/github/organizations`);
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
   * POST /github/organizations
   */
  postGithubOrganization(organization) {
    const url: URL = this.getV1Endpoint('/v1/github/organizations');
    return this.http.post(url, organization).map((res) => res.json());
  }

  /**
   * DELETE /github/organizations/{organization_name}
   */
  deleteGithubOrganization(projectSFID:string, organizationName:string) {
    // outdated: v1
    // const url: URL = this.getV1Endpoint('/v1/github/organizations/' + organizationName);
    const url: URL = this.getV3Endpoint(`/v3/project/${projectSFID}/github/organizations/${organizationName}`);
    return this.http.delete(url);
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

  // getGerritInstance(projectId) {
  //   const url: URL = this.getV1Endpoint('/v1/project/' + projectId + '/gerrits');
  //   return this.http.get(url).map((res) => res.json());
  // }

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
