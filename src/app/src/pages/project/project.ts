import { Component } from '@angular/core';
import { NavController, ModalController, NavParams, IonicPage } from 'ionic-angular';
import { CincoService } from '../../app/services/cinco.service';
import { AssetManagementModal } from '../asset-management/asset-management';

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
  project: {
    icon: string,
    name: string,
    description: string,
    datas: Array<{
      label: string,
      value: string,
    }>,
  };

  members: any;
  // members: Array<{
  //   id: string,
  //   alert?: string,
  //   name: string,
  //   level: string,
  //   status: string,
  //   annual_dues: string,
  //   renewal_date: string,
  // }>;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    private cincoService: CincoService,
    public modalCtrl: ModalController,
  ) {
    // If we navigated to this page, we will have an item available as a nav param
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
        this.project.name = response.name;
        this.project.icon = response.name;
        this.project.datas.push({
            label: "Status",
            value: response.status
        });
        this.project.datas.push({
            label: "Type",
            value: response.type
        });
        this.project.datas.push({
            label: "ID",
            value: response.id
        });
      }
    });
  }

  getProjectMembers(projectId) {
    console.log("getProjectMembers:");
    this.cincoService.getProjectMembers(projectId).subscribe(response => {
      if(response) {
        console.log(response);
        this.members = response;
      }
    });
  }

  memberSelected(event, memberId) {
    this.navCtrl.push('MemberPage', {
      projectId: this.projectId,
      memberId: memberId,
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
      icon: "https://dummyimage.com/600x250/ffffff/000.png&text=project+logo",
      name: "Project",
      description: "This project is a small, scalable, real-time operating system for use on resource-constraned systems supporting multiple architectures...",
      datas: [
        {
          label: "Budget",
          value: "$3,000,000",
        },
        {
          label: "Categories",
          value: "Embedded & IoT",
        },
        {
          label: "Launched",
          value: "5/1/2016",
        },
        {
          label: "Current",
          value: "$2,000,000 ($1,000,000)",
        },
        {
          label: "Members",
          value: "41",
        },
      ],
    };

  }
}
