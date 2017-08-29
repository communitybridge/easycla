import { Injectable } from '@angular/core';
import { HttpClient } from './http-client';

import 'rxjs/Rx';

import { CINCO_API_URL } from './constants'; // TODO: Make sure CINCO_API_URL maps to CINCO URL

@Injectable()
export class CincoService {

  cincoApiUrl: String;

  constructor(public http: HttpClient) {
    this.cincoApiUrl = CINCO_API_URL;
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
    return this.http.get(this.cincoApiUrl + '/project/status')
      .map(res => res.json());
  }

  getProjectCategories() {
    return this.http.get(this.cincoApiUrl + '/project/categories')
      .map(res => res.json());
  }

  getProjectSectors() {
    return this.http.get(this.cincoApiUrl + '/project/sectors')
      .map(res => res.json());
  }

  getAllProjects() {
    return this.http.get(this.cincoApiUrl + '/projects')
      .map(res => res.json());
  }

  //  Disabled for PMC v1.0
  //  createProject(newProject) {
  //    return this.http.post('/projects', newProject)
  //            .map((res) => res.json());
  //  }

  getProject(projectId, getMembers) {
    if (getMembers) { projectId = projectId + '?members=true'; }
    return this.http.get(this.cincoApiUrl + '/projects/' + projectId)
      .map(res => res.json());
  }

  editProject(projectId, editProject) {
    return this.http.put(this.cincoApiUrl + '/projects/' + projectId, editProject)
      .map((res) => res.json());
  }

  getProjectConfig(projectId) {
    return this.http.get(this.cincoApiUrl + '/projects/' + projectId + '/config')
      .map(res => res.json());
  }

  updateProjectManagers(projectId, updatedManagers) {
    return this.http.put(this.cincoApiUrl + '/projects/' + projectId + '/managers', updatedManagers)
      .map((res) => res.json());
  }

  getProjectLogos(projectId) {
    return this.http.get(this.cincoApiUrl + '/projects/' + projectId + '/logos')
      .map(res => res.json());
  }

  /**
  * PUT /projects/{projectId}/logos/{classifier}
  * For logos, any program manager/admin should be able to submit a PUT to `/projects/{projectId}/logos/{classifier}`.
  * The request must have a Content-Type header of `image/**` (e.g. `image/png` or `image/jpeg`)
  * and can optionally include the image bytes in the body.
  * `{classifier}` can be anything, but it's meant to be something like "main" or "black-and-white" or "thumbnail".
  * On a successful PUT, the endpoint is going to respond with a 307 whose
  * Location header (and body) will contain a URL to S3 that the client can resubmit the PUT to.
  * That will actually upload the image.
  **/
  obtainS3URL(projectId, classifier, image) {
    return this.http.putS3URL(this.cincoApiUrl + '/projects/' + projectId + '/logos/' + classifier, image, image.contentType)
      .map(res => res);
  }

  uploadLogo(S3URL, file, contentType) {
    // let headers = new Headers();
    // headers.append('Content-Type', contentType);
    // return this.http.put(S3URL, file, {headers: headers})
    //TODO: Change headers in httpClient to ('Content-Type', contentType);
    return this.http.putUploadS3URL(S3URL, file, contentType)
      // .catch((e) => {
      //   console.log(e)
      //   return e;
      // })
      .map((res) => res );
      // .map(response => response.ok)
  }

  //////////////////////////////////////////////////////////////////////////////

  /**
  * Projects - Members
  * Resources for getting details about project members
  **/

  getProjectMembers(projectId) {
    return this.http.get(this.cincoApiUrl + '/projects/' + projectId + '/members')
      .map(res => res.json());
  }

  getMember(projectId, memberId) {
    return this.http.get(this.cincoApiUrl + '/projects/' + projectId + '/members/' + memberId)
      .map(res => res.json());
  }

  //////////////////////////////////////////////////////////////////////////////

  /**
  * Projects - Members - Contacts
  * Resources for getting and manipulating contacts of project members
  **/

  getMemberContactRoles() {
    return this.http.get(this.cincoApiUrl + '/project/members/contacts/types')
      .map(res => res.json());
  }

  getMemberContacts(projectId, memberId) {
    return this.http.get(this.cincoApiUrl + '/projects/' + projectId + '/members/' + memberId + '/contacts')
      .map(res => res.json());
  }

  addMemberContact(projectId, memberId, contactId, newContact) {
    return this.http.post(this.cincoApiUrl + '/projects/' + projectId + '/members/' + memberId + '/contacts/' + contactId, newContact)
      .map((res) => res.json());
  }

  removeMemberContact(projectId, memberId, contactId, roleId) {
    return this.http.delete(this.cincoApiUrl + '/projects/' + projectId + '/members/' + memberId + '/contacts/' + contactId + '/roles/' + roleId)
      .map((res) => res.json());
  }

  updateMemberContact(projectId, memberId, contactId, roleId, updatedContact) {
    return this.http.put(this.cincoApiUrl + '/projects/' + projectId + '/members/' + memberId + '/contacts/' + contactId + '/roles/' + roleId, updatedContact)
      .map((res) => res.json());
  }

  //////////////////////////////////////////////////////////////////////////////

  /**
  * Organizations - Contacts
  * Resources for getting and manipulating contacts of organizations
  **/

  getOrganizationContactTypes() {
    return this.http.get(this.cincoApiUrl + '/organizations/contacts/types')
      .map(res => res.json());
  }

  getOrganizationContacts(organizationId) {
    return this.http.get(this.cincoApiUrl + '/organizations/' + organizationId + '/contacts')
      .map(res => res.json());
  }

  createOrganizationContact(organizationId, newContact) {
    return this.http.post(this.cincoApiUrl + '/organizations/' + organizationId + '/contacts', newContact)
      .map((res) => res.json());
  }

  getOrganizationContact(organizationId, contactId) {
    return this.http.get(this.cincoApiUrl + '/organizations/' + organizationId + '/contacts/' + contactId)
      .map(res => res.json());
  }

  updateOrganizationContact(organizationId, contactId, contact) {
    return this.http.put(this.cincoApiUrl + '/organizations/' + organizationId + '/contacts/' + contactId, contact)
      .map((res) => res.json());
  }

  //////////////////////////////////////////////////////////////////////////////

  /**
  * Organizations - Projects:
  * Resources for getting details about an organizations project membership
  **/

  getOrganizationProjectMemberships(organizationId) {
    return this.http.get(this.cincoApiUrl + '/organizations/' + organizationId + '/projects_member')
      .map(res => res.json());
  }

  //////////////////////////////////////////////////////////////////////////////

  /**
  * Users
  * Resources to manage internal LF users and roles
  **/

  getCurrentUser() {
    return this.http.get(this.cincoApiUrl + '/users/current')
      .map(res => res.json());
  }

  getAllUsers() {
    return this.http.get(this.cincoApiUrl + '/users')
      .map(res => res.json());
  }

  createUser(user) {
    return this.http.post(this.cincoApiUrl + '/users', user)
      .map(res => res.json());
  }

  removeUser(userId) {
    return this.http.delete(this.cincoApiUrl + '/users/' + userId)
      .map(res => res.json());
  }

  getUser(userId) {
    return this.http.get(this.cincoApiUrl + '/users/' + userId)
      .map(res => res.json());
  }

  getUserRoles() {
    return this.http.get(this.cincoApiUrl + '/users/roles')
      .map(res => res.json());
  }

  updateUser(userId, user) {
    return this.http.put(this.cincoApiUrl + '/users/' + userId, user)
      .map(res => res.json());
  }

  addUserRole(userId, role) {
    return this.http.put(this.cincoApiUrl + '/users/' + userId + '/role/' + role, null)
      .map(res => res.json());
  }

  removeUserRole(userId, role) {
    return this.http.delete(this.cincoApiUrl + '/users/' + userId + '/role/' + role)
      .map(res => res.json());
  }

  //////////////////////////////////////////////////////////////////////////////

}
