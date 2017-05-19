import { Component } from '@angular/core';

import { NavController, NavParams, IonicPage } from 'ionic-angular';

import { CincoService } from '../../app/services/cinco.service'

@IonicPage({
  segment: 'project-details/:projectId'
})
@Component({
  selector: 'project-details',
  templateUrl: 'project-details.html'
})
export class ProjectDetailsPage {

  projectId: string;

  // This project definition is based on CINCO project class
  project: {
    id: string,
    name: string,
    description: string,
    managers: string,
    members: any,
    status: string,
    category: string,
    sector: string,
    url: string,
    startDate: string,
    logoRef: string,
    agreementRef: string,
    mailingListType: string,
    emailAliasType: string,
    address: string
  };

  membershipsCount: number;

  editProject: any;

  constructor(public navCtrl: NavController, public navParams: NavParams, private cincoService: CincoService) {
    this.editProject = {};
    this.projectId = navParams.get('projectId');
    this.getDefaults();
  }

  ngOnInit() {
    this.getProject(this.projectId);
  };

  getProject(projectId) {
    let getMembers = true;
    this.cincoService.getProject(projectId, getMembers).subscribe(response => {
      if(response) {
        this.project.id = response.id;
        this.project.name = response.name;
        this.project.description = response.description;
        this.project.managers = response.managers;
        this.project.members = response.members
        this.project.status = response.status;
        this.project.category = response.category;
        this.project.sector = response.sector;
        this.project.url = response.url;
        this.project.startDate = response.startDate;
        this.project.logoRef = response.logoRef;
        this.project.agreementRef = response.agreementRef;
        this.project.mailingListType = response.mailingListType;
        this.project.emailAliasType = response.emailAliasType;
        this.project.address = response.address.address;
      }
    });
  }

  submitEditProject() {
    this.editProject = {
      project_name: this.project.name,
      project_description: this.project.description,
      project_url: this.project.url
      // TODO: WIP
      // Based on PMC-70 criteria
      // categories, now sector
      // domain
      // billing address (free text)
      // country (picklist)
      // zip code (optional / don't do this if it's hard: validate based on country)
    };
    this.cincoService.editProject(this.projectId, this.editProject).subscribe(response => {
      this.navCtrl.push('ProjectPage', {
        projectId: this.projectId
      });
    });
  }

  cancelEditProject() {
    this.navCtrl.push('ProjectPage', {
      projectId: this.projectId
    });
  }

  changeLogo(){
    // TODO: WIP
    alert("Change Logo");
  }

  getDefaults() {
    this.project = {
      id: "",
      name: "",
      description: "",
      managers: "",
      members: "",
      status: "",
      category: "",
      sector: "",
      url: "",
      startDate: "",
      logoRef: "",
      agreementRef: "",
      mailingListType: "",
      emailAliasType: "",
      address: ""
    };
  }

}
