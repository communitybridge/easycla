import { Component } from '@angular/core';
import {
  NavController,
  NavParams,
  ViewController,
  IonicPage,
  Events,
  ModalController,
  PopoverController,
} from 'ionic-angular';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ClaService } from '../../services/cla.service'
import {SortService} from "../../services/sort.service";
import {KeycloakService} from "../../services/keycloak/keycloak.service";
import {RolesService} from "../../services/roles.service";

@IonicPage({
  segment: 'cla-contract-view-signatures-modal'
})
@Component({
  selector: 'cla-contract-view-signatures-modal',
  templateUrl: 'cla-contract-view-signatures-modal.html',
})
export class ClaContractViewSignaturesModal {
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
    public viewCtrl: ViewController,
    public modalCtrl: ModalController,
    private popoverCtrl: PopoverController,
    private keycloak: KeycloakService,
    public rolesService: RolesService,
    public events: Events
  ) {
    this.claProjectId = this.navParams.get('claProjectId');
    this.getDefaults();

    events.subscribe('modal:close', () => {
      this.dismiss();
    });
  }

  ngOnInit() {
    this.getSignatures();
  }

  getDefaults() {
    this.loading = {
      signatures: true,
    };
    this.sort = {
      signatureType: {
        arrayProp: 'signatureType',
        sortType: 'text',
        sort: null,
      },
      name: {
        arrayProp: 'referenceEntity.user_name',
        sortType: 'text',
        sort: null,
      },
      company: {
        arrayProp: 'signature_user_ccla_company_id',
        sortType: 'text',
        sort: null,
      },
      githubId: {
        arrayProp: 'referenceEntity.user_github_id',
        sortType: 'number',
        sort: null,
      },
      version: {
        arrayProp: 'documentVersion',
        sortType: 'semver',
        sort: null,
      },
      date: {
        arrayProp: 'date_modified',
        sortType: 'date',
        sort: null,
      },
    };
    this.signatures = [];
  }

  getSignatures() {
    this.claService.getProjectSignatures(this.claProjectId).subscribe((signatures) => {
      for (let signature of signatures) {
        signature.documentVersion = `${signature.signature_document_major_version}.${signature.signature_document_minor_version}`;
        if (signature.signature_reference_type === "user") {
          // ICLA, Employee CCLA
          this.claService.getUser(signature.signature_reference_id).subscribe(user => {
            if (user)  {
              signature.user = user;
              //Employee CCLA if signature includes a ccla_company_id
              if (signature.signature_user_ccla_company_id) {
                this.claService.getCompany(signature.signature_user_ccla_company_id).subscribe(company => {
                  if (company) {
                    signature.company = company;
                  }
                })
              }
            }
          });
        }
        else if (signature.signature_reference_type === "company") {
          // CCLA signed by company authority
          this.claService.getCompany(signature.signature_reference_id).subscribe(company => {
            if(company) {
              signature.company = company;
              this.claService.getUser(company.company_manager_id).subscribe(user => {
                if (user) {
                  signature.user = user; 
                }
              })
            }
          });
        }
      }
      this.signatures = signatures;
      this.loading.signatures = false;
    });
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

  dismiss() {
    this.viewCtrl.dismiss();
  }

}
