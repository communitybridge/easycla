import { Component } from '@angular/core';

import { NavController, IonicPage } from 'ionic-angular';

import { CincoService } from '../../app/services/cinco.service'
import { ProjectPage } from '../project/project';


@IonicPage({
  segment: 'project-add'
})
@Component({
  selector: 'add-project',
  templateUrl: 'add-project.html'
})
export class AddProjectPage {
  newProject;
  project_name: String;
  project_type: String;

  constructor(public navCtrl: NavController, private cincoService: CincoService) {
    this.newProject = {};
  }

  submitNewProject() {
    this.newProject = {
      project_name: this.project_name,
      project_type: this.project_type
    };
    this.cincoService.postProject(this.newProject).subscribe(response => {
      console.log(response)
      this.navCtrl.push(ProjectPage, {
        projectId: response
      });
    });
  }

}
