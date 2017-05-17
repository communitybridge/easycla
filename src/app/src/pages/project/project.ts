import { Component } from '@angular/core';
import { NavController, ModalController, NavParams, IonicPage } from 'ionic-angular';
import { CincoService } from '../../app/services/cinco.service';
import { AssetManagementModal } from '../asset-management/asset-management';

@IonicPage({
  segment: 'project/:projectId'
})
@Component({
  selector: 'project',
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

  membersCount: number;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    private cincoService: CincoService,
    public modalCtrl: ModalController,
  ) {
    this.selectedProject = navParams.get('project');
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
        this.project.members = response.members;
        this.membersCount = this.project.members.length;
      }
    });
  }

  memberSelected(event, memberId) {
    this.navCtrl.push('MemberPage', {
      projectId: this.projectId,
      memberId: memberId,
    });
  }

  viewProjectDetails(projectId){
    this.navCtrl.push('ProjectDetailsPage', {
      projectId: projectId
    });
  }

  openAssetManagementModal() {
    let modal = this.modalCtrl.create(AssetManagementModal, {
      projectId: this.projectId,
    });
    modal.present();
  }

  getDefaults() {
    this.project = {
      id: "",
      name: "Project",
      description: "Description",
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
