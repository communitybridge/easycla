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
  editProject;
  project_name: String;
  project_type: String;

  projectId: string;

  // This project definition is based on CINCO project class
  project: {
    id: string,
    name: string,
    description: string,
    managers: string,
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

  memberships: any;
  membershipsCount: number;

  constructor(public navCtrl: NavController, public navParams: NavParams, private cincoService: CincoService) {
    this.editProject = {};
    this.projectId = navParams.get('projectId');
    this.getDefaults();
  }

  ngOnInit() {
    this.getProject(this.projectId);
  };

  getProject(projectId) {
    this.cincoService.getProject(projectId).subscribe(response => {
      if(response) {
        this.project.id = response.id;
        this.project.name = response.name;
        this.project.description = response.description;
        this.project.managers = response.managers;
        this.project.status = response.status;
        this.project.category = response.category;
        this.project.sector = response.sector;
        this.project.url = response.url;
        this.project.startDate = response.startDate;
        this.project.logoRef = response.logoRef;
        this.project.agreementRef = response.agreementRef;
        this.project.mailingListType = response.mailingListType;
        this.project.emailAliasType = response.emailAliasType;
        this.project.address = response.address;
      }
    });
  }

  submitEditProject() {
    this.editProject = {
      project_name: this.project_name,
      project_type: this.project_type
    };
    this.cincoService.editProject(this.editProject).subscribe(response => {
      this.navCtrl.push('ProjectPage', {
        projectId: response
      });
    });
  }

  getDefaults() {
    this.project = {
      id: "",
      name: "",
      description: "This project is a small, scalable, real-time operating system for use on resource-constraned systems supporting multiple architectures...",
      managers: "",
      status: "",
      category: "",
      sector: "",
      url: "",
      startDate: "",
      logoRef: "https://dummyimage.com/600x250/ffffff/000.png&text=project+logo",
      agreementRef: "",
      mailingListType: "",
      emailAliasType: "",
      address: ""
    };

    this.memberships = [
      {
        tier: "Gold",
        numberOfMembers: "2",
        annualCost: "$200,000",
        boardSeat: "Yes"
      },{
        tier: "Platinum",
        numberOfMembers: "8",
        annualCost: "$600,000",
        boardSeat: "No"
      }
    ];
  }

}
