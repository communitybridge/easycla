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

  getMember(projectId, memberId) {

  }

}
