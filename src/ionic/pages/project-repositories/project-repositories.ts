import { Component } from '@angular/core';
import { NavController, ModalController, NavParams, IonicPage } from 'ionic-angular';
import { CincoService } from '../../services/cinco.service';
import { KeycloakService } from '../../services/keycloak/keycloak.service';
import { SortService } from '../../services/sort.service';
import { ProjectModel } from '../../models/project-model';

@IonicPage({
  segment: 'project/:projectId/repositories'
})
@Component({
  selector: 'project-repositories',
  templateUrl: 'project-repositories.html',
  providers: [CincoService]
})
export class ProjectRepositoriesPage {
  projectId: string;
  project = new ProjectModel();

  membersCount: number;
  loading: any;
  sort: any;
  tab = 'repositories';

  contracts: any;

  iclaUploadInfo: any;
  cclaUploadInfo: any;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    private cincoService: CincoService,
    private sortService: SortService,
    public modalCtrl: ModalController,
    private keycloak: KeycloakService,
  ) {
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
    let getMembers = false;
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
        this.loading.project = false;
      }
    });
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

  openClaContractVersionModal(uploadInfo) {
    let modal = this.modalCtrl.create('ClaContractVersionModal', {
    });
    modal.present();
  }

  openClaOrganizationProviderModal() {
    let modal = this.modalCtrl.create('ClaOrganizationProviderModal', {
    });
    modal.present();
  }

  openProjectPage() {
    this.navCtrl.push('ProjectPage', {
      projectId: this.projectId,
    });
  }

  getDefaults() {
    this.loading = {
      project: true,
    };
    // this.project = {
    //   id: "",
    //   name: "Project",
    //   description: "Description",
    //   managers: "",
    //   members: [],
    //   status: "",
    //   category: "",
    //   sector: "",
    //   url: "",
    //   startDate: "",
    //   logoRef: "",
    //   agreementRef: "",
    //   mailingListType: "",
    //   emailAliasType: "",
    //   address: {
    //     address: {
    //       administrativeArea: "",
    //       country: "",
    //       localityName: "",
    //       postalCode: "",
    //       thoroughfare: ""
    //     },
    //     type: ""
    //   }
    // };

    this.contracts = [
      {
        id: '000001',
        name: 'Zephyr Contract',
        ccla: true,
        cclaAndIcla: true,
        icla: true,
        contracts: {
          ccla: {
            name: 'zephyr_CLA_corporate.pdf',
            src: 'https://example.com/something.pdf',
            uploadDate: '6/30/2017, 11:31 PST',
            version: '1.0',
          },
          icla: {
            name: 'zephyr_CLA_corporate.pdf',
            src: 'https://example.com/something.pdf',
            uploadDate: '6/30/2017, 11:31 PST',
            version: '1.0',
          }
        },
        organizations: [

        ],
      },
    ];
  }

  sortMembers(prop) {
    this.sortService.toggleSort(
      this.sort,
      prop,
      this.project.members,
    );
  }

}
