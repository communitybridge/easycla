import { Component } from '@angular/core';
import { NavController, ModalController, NavParams, IonicPage } from 'ionic-angular';
import { ClaService } from 'cla-service';
import { KeycloakService } from '../../services/keycloak/keycloak.service';
import { SortService } from '../../services/sort.service';
import { PopoverController } from 'ionic-angular';

@IonicPage({
  segment: 'cla-contracts-contributors/:claProjectId'
})
@Component({
  selector: 'cla-contracts-contributors',
  templateUrl: 'cla-contracts-contributors.html',
})
export class ClaContractsContributorsPage {
  selectedProject: any;
  claProjectId: string;

  loading: any;
  sort: any;
  signatures: any;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    private claService: ClaService,
    private sortService: SortService,
    public modalCtrl: ModalController,
    private popoverCtrl: PopoverController,
    private keycloak: KeycloakService
  ) {
    this.claProjectId = this.navParams.get('claProjectId');
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
    this.claService.getProjectSignatures(this.claProjectId).subscribe((signatures) => {
      console.log("signatures");
      console.log(signatures);
    });
  }

  getDefaults() {
    this.loading = {
      // project: true,
    };
    this.sort = {
      entity: {
        arrayProp: 'entity',
        sortType: 'text',
        sort: null,
      },
      company: {
        arrayProp: 'company',
        sortType: 'text',
        sort: null,
      },
      name: {
        arrayProp: 'name',
        sortType: 'text',
        sort: null,
      },
      email: {
        arrayProp: 'email',
        sortType: 'text',
        sort: null,
      },
      version: {
        arrayProp: 'version',
        sortType: 'semver',
        sort: null,
      },
    };

    this.signatures = [
      {
        entity: "corporate",
        company: "E Corp.",
        name: "Aname Lastname",
        email: "cname@example.com",
        version: "1.0",
      },
      {
        entity: "individual",
        company: "F Corp.",
        name: "Bname Lastname",
        email: "gname@example.com",
        version: "1.1",
      },
      {
        entity: "employee",
        company: "A Corp.",
        name: "Cname Lastname",
        email: "bname@example.com",
        version: "1.5",
      },
      {
        entity: "employee",
        company: "G Corp.",
        name: "Dname Lastname",
        email: "fname@example.com",
        version: "1.10",
      },
      {
        entity: "individual",
        company: "D Corp.",
        name: "Ename Lastname",
        email: "bname@example.com",
        version: "2.0",
      },
      {
        entity: "individual",
        company: "B Corp.",
        name: "Gname Lastname",
        email: "aname@example.com",
        version: "2.5",
      },
      {
        entity: "employee",
        company: "C Corp.",
        name: "Fname Lastname",
        email: "dname@example.com",
        version: "2.10",
      },
    ];
  }

  sortMembers(prop) {
    this.sortService.toggleSort(
      this.sort,
      prop,
      this.signatures,
    );
  }

  signaturePopover(ev, signature) {
    let actions = {
      items: [
        {
          label: 'Details',
          callback: 'signatureDetails',
          callbackData: {
            signature: signature,
          }
        },
        {
          label: 'CLA',
          callback: 'signatureCla',
          callbackData: {
            signature: signature,
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

  signatureDetails(data) {
    console.log('signature details');
  }

  signatureCla(data) {
    console.log('signature cla');
  }

}
