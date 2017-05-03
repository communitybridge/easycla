import { Injectable } from '@angular/core';
import { Http, Headers } from '@angular/http';

import 'rxjs/Rx';

@Injectable()
export class CincoService{
  http: any;
  baseUrl: String;

  constructor(http: Http) {
    this.http = http;
    this.baseUrl = '';
  }

  getAllProjects() {
    return this.http.get(this.baseUrl + '/get_all_projects')
            .map(res => res.json());
  }

  getProject(projectId) {
    return this.http.get(this.baseUrl + '/get_project/' + projectId)
            .map(res => res.json());
  }

  postProject(newProject) {
    let headers = new Headers({ 'Content-Type': 'application/json' });
    let body = new FormData();
    body.append('project_name', newProject.project_name);
    // body.append('project_type', newProject.project_type);
    return this.http.post('/post_project', body, headers)
                .map((res) => res.json());
  }

  getProjectMembers(projectId) {
    return this.http.get(this.baseUrl + '/members/' + projectId)
            .map(res => res.json());
  }

  getMember(projectId, memberId) {
    console.log('getMember called');
    var response = this.http.get(this.baseUrl + '/member/' + projectId + '/' + memberId)
            .map(res => res.json());
    console.log(response);
    return response;
  }

  updateMemberContact(projectId, memberId, contactId, contact) {
    console.log('updateMemberContact called');
    let headers = new Headers({ 'Content-Type': 'application/json' });
    let body = new FormData();
    body.append('projectId', projectId);
    body.append('memberId', memberId);
    body.append('contactId', contactId);
    body.append('contactEmail', contact.email);
    body.append('contactBio', contact.bio);
    body.append('contactPhone', contact.phone);
    console.log(body);
    return this.http.post('/update_member_contact', body, headers)
                .map((res) => res.json());
  }

  addMemberContact(projectId, memberId, contact) {
    console.log('updateMemberContact called');
    let headers = new Headers({ 'Content-Type': 'application/json' });
    let body = new FormData();
    body.append('projectId', projectId);
    body.append('memberId', memberId);
    body.append('contactEmail', contact.email);
    body.append('contactBio', contact.bio);
    body.append('contactPhone', contact.phone);
    console.log(body);
    return this.http.post('/add_member_contact', body, headers)
                .map((res) => res.json());
  }

  removeMemberContact(projectId, memberId, contactId) {
    console.log('updateMemberContact called');
    let headers = new Headers({ 'Content-Type': 'application/json' });
    let body = new FormData();
    body.append('projectId', projectId);
    body.append('memberId', memberId);
    body.append('contactId', contactId);
    console.log(body);
    return this.http.post('/remove_member_contact', body, headers)
                .map((res) => res.json());
  }

}
