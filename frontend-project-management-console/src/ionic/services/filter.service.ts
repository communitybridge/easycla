// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import { Injectable } from '@angular/core';

@Injectable()
export class FilterService {

  constructor() {}

  filterAllProjects(allProjects, projectProperty, keyword){
    return allProjects.filter((projects) => {
      if (projectProperty == 'managers') { return projects[projectProperty].indexOf(keyword.toLowerCase()) > -1; }
      else { return projects[projectProperty] == keyword; }
    });
  }

  resetFilter(data){
    return JSON.parse(JSON.stringify(data));
  }

}
