import { Component } from '@angular/core';
import { NavController, ModalController, NavParams, IonicPage } from 'ionic-angular';
import { CincoService } from '../../services/cinco.service';
import { KeycloakService } from '../../services/keycloak/keycloak.service';
import { SortService } from '../../services/sort.service';
import { ProjectModel } from '../../models/project-model';
import { PopoverController } from 'ionic-angular';

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
    private popoverCtrl: PopoverController,
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

  openClaOrganizationAppModal(orgName) {
    let modal = this.modalCtrl.create('ClaOrganizationAppModal', {
      orgName: orgName,
    });
    modal.present();
  }

  openProjectPage() {
    this.navCtrl.push('ProjectPage', {
      projectId: this.projectId,
    });
  }

  openClaContractsContributorsPage(contractId) {
    this.navCtrl.push('ClaContractsContributorsPage', {
      contractId: contractId,
    });
  }

  getDefaults() {
    this.loading = {
      project: true,
    };

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
          {
            id: "000001",
            name: "Zephyr Project",
            gitUrl: "https://github.com/zephyrproject-rtos",
            description: "https://www.zephyrproject.org",
            appConnected: true,
            repositories: [
              {
                name: "zephyr",
              },
              {
                name: "zephyr-bluetooth",
              },
              {
                name: "zephyr-wifi",
              },
            ],
          },
        ],
      },
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
          {
            id: "000001",
            name: "Zephyr Project",
            gitUrl: "https://github.com/zephyrproject-rtos",
            description: "https://www.zephyrproject.org",
            appConnected: false,
            repositories: [
              {
                name: "zephyr",
              },
              {
                name: "zephyr-bluetooth",
              },
              {
                name: "zephyr-wifi",
              },
            ],
          },
        ],
      },
    ];
  }

  organizationPopover(ev, organization) {
    let actions = {
      items: [
        {
          label: 'Delete',
          callback: 'organizationDelete',
          callbackData: {
            organization: organization,
          }
        },
      ]
    };
    let popover = this.popoverCtrl.create(
      'ActionPopoverComponent',
      actions,
    );

    popover.present({
      ev: ev
    });

    popover.onDidDismiss((popoverData) => {
      if(popoverData) {
        this.popoverResponse(popoverData);
      }
    });
  }

  /**
   * Called if popover dismissed with data. Passes data to a callback function
   * @param  {object} popoverData should contain .callback and .callbackData
   */
  popoverResponse(popoverData) {
    let callback = popoverData.callback;
    if(this[callback]) {
      this[callback](popoverData.callbackData);
    }
  }

  organizationDelete(data) {
    console.log('organization delete');
    console.log(data.organization);
  }

}
