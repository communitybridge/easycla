import { Component } from '@angular/core';
import { NavController, NavParams, IonicPage } from 'ionic-angular';
import { CincoService } from '../../app/services/cinco.service'

@IonicPage({
  segment: 'project-page/:projectId'
})
@Component({
  selector: 'page-project',
  templateUrl: 'project.html'
})
export class ProjectPage {
  selectedProject: any;
  projectId: string;

  // This project definition is based on CINCO project class
  project: {
    id: string,
    name: string,
    description: string,
    managers: string,
    projectStatus: string,
    projectCategory: string,
    url: string,
    startDate: string,
    logoRef: string,
    agreementRef: string,
    mailingListType: string,
    emailAliasType: string
  };

  members: any;
  membersCount: number;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    private cincoService: CincoService
  ) {
    this.selectedProject = navParams.get('project');
    this.projectId = navParams.get('projectId');
    this.getDefaults();
  }

  ngOnInit() {
    this.getProject(this.projectId);
    this.getProjectMembers(this.projectId);
  };

  getProject(projectId) {
    this.cincoService.getProject(projectId).subscribe(response => {
      if(response) {
        this.project.id = response.id;
        this.project.name = response.name;
        this.project.description = response.description;
        this.project.managers = response.managers;
        this.project.projectStatus = response.projectStatus;
        this.project.projectCategory = response.projectCategory;
        this.project.url = response.url;
        this.project.startDate = response.startDate;
        this.project.logoRef = response.logoRef;
        this.project.agreementRef = response.agreementRef;
        this.project.mailingListType = response.mailingListType;
        this.project.emailAliasType = response.emailAliasType;
      }
    });
  }

  getProjectMembers(projectId) {
    this.cincoService.getProjectMembers(projectId).subscribe(response => {
      if(response) {
        this.members = response;
        this.membersCount = this.members.length;
      }
    });
  }

  memberSelected(event, memberId) {
    this.navCtrl.push('MemberPage', {
      projectId: this.projectId,
      memberId: memberId,
    });
  }

  getDefaults() {

    this.project = {
      id: "",
      name: "Project",
      description: "This project is a small, scalable, real-time operating system for use on resource-constraned systems supporting multiple architectures...",
      managers: "",
      projectStatus: "",
      projectCategory: "",
      url: "",
      startDate: "",
      logoRef: "https://dummyimage.com/600x250/ffffff/000.png&text=project+logo",
      agreementRef: "",
      mailingListType: "",
      emailAliasType: ""
    };
  }
}
