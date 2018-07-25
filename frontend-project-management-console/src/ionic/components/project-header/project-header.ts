import { Input, Component } from '@angular/core';
import { NavController, ModalController } from 'ionic-angular';
import { CincoService } from '../../services/cinco.service';

@Component({
  selector: 'project-header',
  templateUrl: 'project-header.html',
})
export class ProjectHeaderComponent {

  @Input('projectId')
  private projectId: string;

  project: any;

  projectSectors: any;

  constructor (
    private navCtrl: NavController,
    private cincoService: CincoService,
    public modalCtrl: ModalController,
  ) {
    this.getDefaults();
  }

  ngOnInit() {
    this.getProjectSectors();
    this.getProject(this.projectId);
  }

  getDefaults() {
    this.projectSectors = {};
    this.project = {
      id: "",
      name: "Project",
      description: "Description",
      managers: "",
      members: [],
      status: "",
      category: "",
      sector: "",
      url: "",
      logoRef: "",
      startDate: "",
      agreementRef: "",
      mailingListType: "",
      emailAliasType: "",
      address: {
        address: {
          administrativeArea: "",
          country: "",
          localityName: "",
          postalCode: "",
          thoroughfare: ""
        },
        type: ""
      },
      config: {
        logoRef: ""
      }
    };
  }

  viewProjectDetails(projectId){
    this.navCtrl.push('ProjectDetailsPage', {
      projectId: projectId
    });
  }

  getProject(projectId) {
    let getMembers = true;
    this.cincoService.getProject(projectId, getMembers).subscribe(response => {
      if(response) {
        this.project = response;
      }
    });
  }

  getProjectSectors() {
    this.cincoService.getProjectSectors().subscribe(response => {
      this.projectSectors = response;
    });
  }

  openProjectUserManagementModal() {
    let modal = this.modalCtrl.create('ProjectUserManagementModal', {
      projectId: this.projectId,
      projectName: this.project.name
    });
    modal.present();
  }

  openAssetManagementModal() {
    let modal = this.modalCtrl.create('AssetManagementModal', {
      projectId: this.projectId,
      projectName: this.project.name
    });
    modal.present();
  }
}
