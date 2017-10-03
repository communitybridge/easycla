import { Component } from '@angular/core';
import { NavController, ModalController, NavParams, IonicPage } from 'ionic-angular';
import { CincoService } from '../../services/cinco.service';
import { KeycloakService } from '../../services/keycloak/keycloak.service';
import { SortService } from '../../services/sort.service';
import { ProjectModel } from '../../models/project-model';

@IonicPage({
  segment: 'project/:projectId'
})
@Component({
  selector: 'project',
  templateUrl: 'project.html',
  providers: [CincoService]
})
export class ProjectPage {
  selectedProject: any;
  projectId: string;

  project = new ProjectModel();

  membersCount: number;
  loading: any;
  sort: any;
  tab = 'membership';

  contracts: any;

  iclaUploadInfo: any;
  cclaUploadInfo: any;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    private cincoService: CincoService,
    private sortService: SortService,
    public modalCtrl: ModalController,
    private keycloak: KeycloakService
  ) {
    this.selectedProject = navParams.get('project');
    this.projectId = navParams.get('projectId');
    this.getDefaults();
  }

  ionViewCanEnter() {
    if(!this.keycloak.authenticated())
    {
      this.navCtrl.setRoot('LoginPage');
      this.navCtrl.popToRoot();
    }
    return this.keycloak.authenticated();
  }

  ionViewWillEnter() {
    if(!this.keycloak.authenticated())
    {
      this.navCtrl.push('LoginPage');
    }
  }

  ngOnInit() {
    this.getProject(this.projectId);
  }

  getProject(projectId) {
    let getMembers = true;
    this.cincoService.getProject(projectId, getMembers).subscribe(response => {
      if(response) {
        this.project = response;
        // This is to refresh an image that have same URL
        if(this.project.config.logoRef) { this.project.config.logoRef += "?" + new Date().getTime(); }
        this.membersCount = this.project.members.length;
        this.loading.project = false;
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

  openProjectUserManagementModal() {
    let modal = this.modalCtrl.create('ProjectUserManagementModal', {
      projectId: this.projectId,
      projectName: this.project.name,
    });
    modal.present();
  }

  openAssetManagementModal() {
    let modal = this.modalCtrl.create('AssetManagementModal', {
      projectId: this.projectId,
    });
    modal.onDidDismiss(newlogoRef => {
      if(newlogoRef){
        this.project.config.logoRef = newlogoRef;
      }
    });
    modal.present();
  }

  openClaContractConfigModal(contract) {
    let modal = this.modalCtrl.create('ClaContractConfigModal', {
      contract: contract,
    });
    modal.present();
  }

  openClaContractUploadModal(uploadInfo) {
    let modal = this.modalCtrl.create('ClaContractUploadModal', {
      uploadInfo: uploadInfo,
    });
    modal.present();
  }

  openProjectRepositoriesPage() {
    this.navCtrl.push('ProjectRepositoriesPage', {
      projectId: this.projectId
    });
  }

  getDefaults() {
    this.loading = {
      project: true,
    };
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
    this.sort = {
      alert: {
        arrayProp: 'alert',
        sortType: 'text',
        sort: null,
      },
      company: {
        arrayProp: 'org.name',
        sortType: 'text',
        sort: null,
      },
      product: {
        arrayProp: 'product',
        sortType: 'text',
        sort: null,
      },
      status: {
        arrayProp: 'invoices[0].status',
        sortType: 'text',
        sort: null,
      },
      dues: {
        arrayProp: 'annualDues',
        sortType: 'number',
        sort: null,
      },
      renewal: {
        arrayProp: 'renewalDate',
        sortType: 'date',
        sort: null,
      },
    };
  }

  sortMembers(prop) {
    this.sortService.toggleSort(
      this.sort,
      prop,
      this.project.members,
    );
  }

}
