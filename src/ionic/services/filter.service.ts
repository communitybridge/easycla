import { Injectable } from '@angular/core';

@Injectable()
export class FilterService {

  constructor() {}

  filterAllProjects(allProjects, projectProperty, keyword){
    return allProjects.filter((projects) => {
      return projects[projectProperty] == keyword;
    });
  }

  resetFilter(data){
    return JSON.parse(JSON.stringify(data));
  }

}
