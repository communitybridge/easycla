import { Injectable } from '@angular/core';
import { Http } from '@angular/http';
import 'rxjs/Rx';

@Injectable()
export class CincoService{
  http: any;
  baseUrl: String;

  constructor(http: Http){
    this.http = http;
    this.baseUrl = '';
  }

  getAllProjects(){
    return this.http.get(this.baseUrl + '/get_all_projects')
            .map(res => res.json());
  }

  getProject(projectId){
    return this.http.get(this.baseUrl + '/get_project/' + projectId)
            .map(res => res.json());
  }

}
