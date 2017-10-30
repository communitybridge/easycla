import { Component } from '@angular/core';
import { NavController, ModalController, NavParams, IonicPage } from 'ionic-angular';
import { CincoService } from '../../../services/cinco.service';
import { KeycloakService } from '../../../services/keycloak/keycloak.service';
import { SortService } from '../../../services/sort.service';
import { ProjectModel } from '../../../models/project-model';
import { PopoverController } from 'ionic-angular';

@IonicPage({
  segment: 'project/:projectId/cla'
})
@Component({
  selector: 'project-cla',
  templateUrl: 'project-cla.html',
  providers: [CincoService]
})
export class ProjectClaPage {
  projectId: string;

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

  getDefaults() {

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

  ngOnInit() {

  }

  ionViewCanEnter() {
    if (!this.keycloak.authenticated()) {
      this.navCtrl.setRoot('LoginPage');
      this.navCtrl.popToRoot();
    }
    return this.keycloak.authenticated();
  }

  ionViewWillEnter() {
    if (!this.keycloak.authenticated()) {
      this.navCtrl.push('LoginPage');
    }
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

  openClaContractsContributorsPage(contractId) {
    this.navCtrl.push('ClaContractsContributorsPage', {
      contractId: contractId,
    });
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
