// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Injectable } from '@angular/core';
import { HttpClient } from './http-client';
import { EnvConfig } from './cla.env.utils';
import 'rxjs/Rx';

@Injectable()
export class CincoService {
  cincoApiUrl: String;

  constructor(public http: HttpClient) {
    this.cincoApiUrl = EnvConfig['cinco-api-url'];
  }

  //////////////////////////////////////////////////////////////////////////////

  /**
   * This service should ONLY contain methods calling CINCO API
   **/

  //////////////////////////////////////////////////////////////////////////////

  /**
   * Projects
   * Resources to expose and manipulate details of projects
   **/

  getProjectStatuses() {
    return this.http.get(this.cincoApiUrl + '/project/status').map((res) => res.json());
  }

  getProjectCategories() {
    return this.http.get(this.cincoApiUrl + '/project/categories').map((res) => res.json());
  }

  getProjectSectors() {
    return this.http.get(this.cincoApiUrl + '/project/sectors').map((res) => res.json());
  }

  getMyProjects() {
    return this.http.get(this.cincoApiUrl + '/project').map((res) => res.json());
  }

  getEvents() {
    return this.http.get(this.cincoApiUrl + '/events').map((res) => res.json());
  }

  getAllProjects() {
    return this.http
      .get(this.cincoApiUrl + '/project') // CINCO changes
      .map((res) => res.json());
  }

  //  Disabled for PMC v1.0
  //  createProject(newProject) {
  //    return this.http.post('/projects', newProject)
  //            .map((res) => res.json());
  //  }

  getProject(projectId, getMembers) {
    if (getMembers) {
      projectId = projectId + '?members=true';
    }
    return this.http.get(this.cincoApiUrl + '/project/' + projectId).map((res) => {
      let project = res.json();
      if (project.config.logoRef) {
        project.config.logoRef += '?' + new Date().getTime();
      }
      return project;
    });
  }

  editProject(projectId, editProject) {
    return this.http.put(this.cincoApiUrl + '/project/' + projectId, editProject).map((res) => res.json());
  }

  getProjectConfig(projectId) {
    return this.http.get(this.cincoApiUrl + '/project/' + projectId + '/config').map((res) => res.json());
  }

  editProjectConfig(projectId, updatedConfig) {
    return this.http
      .patch(this.cincoApiUrl + '/projects/' + projectId + '/config', updatedConfig)
      .map((res) => res.json());
  }

  updateProjectManagers(projectId, updatedManagers) {
    return this.http
      .put(this.cincoApiUrl + '/project/' + projectId + '/managers', updatedManagers)
      .map((res) => res.json());
  }

  getProjectLogos(projectId) {
    return this.http.get(this.cincoApiUrl + '/project/' + projectId + '/logos').map((res) => res.json());
  }

  /**
   * This endpoint will return an array of all documents for a project.
   * Each object contains a url property indicating a temporary URL where the actual document can be retrieved with a GET request.
   * The expiresOn property is an ISO-8601 timestamp indicating the last millisecond the URL will be valid to issue a GET request to.
   **/
  getProjectDocuments(projectId) {
    return this.http.get(this.cincoApiUrl + '/project/' + projectId + '/documents').map((res) => res.json());
  }

  /**
   * PUT /project/{projectId}/logos/{classifier}/
   * For logos, any program manager/admin should be able to submit a PUT to `/projects/{projectId}/logos/{classifier}`.
   * The request must have a Content-Type header of `image/**` (e.g. `image/png` or `image/jpeg`)
   * and can optionally include the image bytes in the body.
   * `{classifier}` can be anything, but it's meant to be something like "main" or "black-and-white" or "thumbnail".
   * On a successful PUT, the endpoint is going to respond with a 307 whose
   * Location header (and body) will contain a URL to S3 that the client can resubmit the PUT to.
   * That will actually upload the image.
   **/
  obtainLogoS3URL(projectId, classifier, image) {
    return this.http
      .put(this.cincoApiUrl + '/project/' + projectId + '/logos/' + classifier, image, image.contentType)
      .map((res) => res);
  }

  /**
   * PUT /project/{projectId}/documents/{classifier}/{filename}/
   * Accepts a request whose body is optionally blank (the body will be ignored).
   * The request headers must specify the Content-Type of the document.
   * The classifier path parameter corresponds as a single word class of this logo.
   * The values of the classifier path are limited to known and expected values, such as "minutes", or "bylaws".
   * A successful request will result in a 307, indicating that the client should retry the request with
   * the returned URI found in the responses location header and the binary contents of the logo in the request body
   * The response body contains details about the request. The putUrl.url property is where the real upload request should occur.
   * The putUrl.expiresOn property is an ISO-8601 timestamp indicating the last millisecond the putURL will be valid.
   **/
  obtainDocumentS3URL(projectId, classifier, file, filename, contentType) {
    return this.http
      .put(this.cincoApiUrl + '/project/' + projectId + '/documents/' + classifier + '/' + filename, file, contentType)
      .map((res) => res);
  }

  uploadToS3(S3URL, file, contentType) {
    return this.http.putS3(S3URL, file, contentType).map((res) => res);
  }

  //////////////////////////////////////////////////////////////////////////////

  /**
   * Projects - Mailing List - groups.io
   * Resources for getting details about project members
   **/

  getAllProjectGroups(projectId) {
    return this.http.get(this.cincoApiUrl + '/project/' + projectId + '/mailinglists').map((res) => res.json());
  }

  createProjectGroup(projectId, group) {
    return this.http.post(this.cincoApiUrl + '/project/' + projectId + '/mailinglists', group).map((res) => res.json());
  }

  getProjectGroup(projectId, groupId) {
    return this.http
      .get(this.cincoApiUrl + '/project/' + projectId + '/mailinglists' + groupId)
      .map((res) => res.json());
  }

  removeProjectGroup(projectId, groupId) {
    return this.http
      .delete(this.cincoApiUrl + '/project/' + projectId + '/mailinglists/' + groupId)
      .map((res) => res.json());
  }

  addGroupParticipant(projectId, groupId, participant) {
    return this.http
      .post(this.cincoApiUrl + '/project/' + projectId + '/mailinglists/' + groupId + '/participants', participant)
      .map((res) => res.json());
  }

  removeGroupParticipant(projectId, groupId, participantEmail) {
    return this.http
      .delete(
        this.cincoApiUrl + '/project/' + projectId + '/mailinglists/' + groupId + '/participants/' + participantEmail
      )
      .map((res) => res.json());
  }
  //////////////////////////////////////////////////////////////////////////////

  /**
   * Projects - Members
   * Resources for getting details about project members
   **/

  getProjectMembers(projectId) {
    return this.http.get(this.cincoApiUrl + '/project/' + projectId + '/members').map((res) => res.json());
  }

  getMember(projectId, memberId) {
    return this.http.get(this.cincoApiUrl + '/project/' + projectId + '/members/' + memberId).map((res) => res.json());
  }

  //////////////////////////////////////////////////////////////////////////////

  /**
   * Projects - Members - Contacts
   * Resources for getting and manipulating contacts of project members
   **/

  getMemberContactRoles() {
    return this.http.get(this.cincoApiUrl + '/project/members/contacts/types').map((res) => res.json());
  }

  getMemberContacts(projectId, memberId) {
    return this.http
      .get(this.cincoApiUrl + '/project/' + projectId + '/members/' + memberId + '/contacts')
      .map((res) => res.json());
  }

  addMemberContact(projectId, memberId, contactId, newContact) {
    return this.http
      .post(this.cincoApiUrl + '/project/' + projectId + '/members/' + memberId + '/contacts/' + contactId, newContact)
      .map((res) => res.json());
  }

  removeMemberContact(projectId, memberId, contactId, roleId) {
    return this.http
      .delete(
        this.cincoApiUrl +
          '/project/' +
          projectId +
          '/members/' +
          memberId +
          '/contacts/' +
          contactId +
          '/roles/' +
          roleId
      )
      .map((res) => res.json());
  }

  updateMemberContact(projectId, memberId, contactId, roleId, updatedContact) {
    return this.http
      .put(
        this.cincoApiUrl +
          '/project/' +
          projectId +
          '/members/' +
          memberId +
          '/contacts/' +
          contactId +
          '/roles/' +
          roleId,
        updatedContact
      )
      .map((res) => res.json());
  }

  //////////////////////////////////////////////////////////////////////////////

  /**
   * Organizations - Contacts
   * Resources for getting and manipulating contacts of organizations
   **/

  getOrganizationContactTypes() {
    return this.http.get(this.cincoApiUrl + '/organizations/contacts/types').map((res) => res.json());
  }

  getOrganizationContacts(organizationId) {
    return this.http.get(this.cincoApiUrl + '/organizations/' + organizationId + '/contacts').map((res) => res.json());
  }

  createOrganizationContact(organizationId, newContact) {
    return this.http
      .post(this.cincoApiUrl + '/organizations/' + organizationId + '/contacts', newContact)
      .map((res) => res.json());
  }

  getOrganizationContact(organizationId, contactId) {
    return this.http
      .get(this.cincoApiUrl + '/organizations/' + organizationId + '/contacts/' + contactId)
      .map((res) => res.json());
  }

  updateOrganizationContact(organizationId, contactId, contact) {
    return this.http
      .put(this.cincoApiUrl + '/organizations/' + organizationId + '/contacts/' + contactId, contact)
      .map((res) => res.json());
  }

  //////////////////////////////////////////////////////////////////////////////

  /**
   * Organizations - Projects:
   * Resources for getting details about an organizations project membership
   **/

  getOrganizationProjectMemberships(organizationId) {
    return this.http
      .get(this.cincoApiUrl + '/organizations/' + organizationId + '/project_members')
      .map((res) => res.json());
  }

  //////////////////////////////////////////////////////////////////////////////

  /**
   * User
   * Resources to manage internal LF users and roles
   **/

  searchUser(username) {
    return this.http.get(this.cincoApiUrl + '/user/search/username/' + username).map((res) => res.json());
  }

  searchUserTerm(term) {
    return this.http.get(this.cincoApiUrl + '/user/search/' + term).map((res) => res.json());
  }

  getCurrentUser() {
    return this.http.get(this.cincoApiUrl + '/user').map((res) => res.json());
  }

  createUser(user) {
    return this.http.post(this.cincoApiUrl + '/user', user).map((res) => res.json());
  }

  removeUser(userId) {
    return this.http.delete(this.cincoApiUrl + '/user/' + userId).map((res) => res.json());
  }

  getUser(userId) {
    return this.http.get(this.cincoApiUrl + '/user/' + userId).map((res) => res.json());
  }

  getUserRoles() {
    return this.http.get(this.cincoApiUrl + '/user/roles').map((res) => res.json());
  }

  updateUser(userId, user) {
    return this.http.put(this.cincoApiUrl + '/user/' + userId, user).map((res) => res.json());
  }

  addUserRole(userId, role) {
    return this.http.put(this.cincoApiUrl + '/user/' + userId + '/role/' + role, null).map((res) => res.json());
  }

  removeUserRole(userId, role) {
    return this.http.delete(this.cincoApiUrl + '/user/' + userId + '/role/' + role).map((res) => res.json());
  }

  //////////////////////////////////////////////////////////////////////////////
}
