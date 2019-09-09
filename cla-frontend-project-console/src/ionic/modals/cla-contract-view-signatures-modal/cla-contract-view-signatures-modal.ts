// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

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
import { SortService } from "../../services/sort.service";
import { KeycloakService } from "../../services/keycloak/keycloak.service";
import { RolesService } from "../../services/roles.service";


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
  signatures: any[];
  searchTerm: string;
  columns: any[];
  rows: any[]

  companies: any[];
  users: any[];
  filteredData: any[];

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
      signatures: true
    };
    this.searchTerm = '',
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
    this.filteredData = this.rows
    this.columns = [
      { prop: 'Entity Type' },
      { prop: 'Name' },
      { prop: 'Company' },
      { prop: 'GithubID' },
      { prop: 'LFID' },
      { prop: 'Version' },
      { prop: 'Date' }
    ];
  }

  async getUser(signatureReferenceId) {
    return await this.claService.getUser(signatureReferenceId).toPromise();
  }

  async getCompany(referenceId) {
    return await this.claService.getCompany(referenceId).toPromise();
  }

  getSignatures() {
    this.claService.getProjectSignatures(this.claProjectId).subscribe((signatures) => {
      for (let signature of signatures) {
        signature.documentVersion = `${signature.signature_document_major_version}.${signature.signature_document_minor_version}`;
        // Include only signed signatures into SignedSignatures
        if (signature.signature_signed === true) {
          if (signature.signature_reference_type === "user") {
            // ICLA, Employee CCLA
            this.getUser(signature.signature_reference_id).then((user) => {
              if (user) {
                signature.user = { ...user };
                if (signature.signature_user_ccla_company_id) {
                  this.getCompany(signature.signature_user_ccla_company_id).then((company) => {
                    if (company) {
                      signature.company = { ...company };
                    }
                  }).catch((err) => {
                    console.log(err, 'this is company service error')
                  })
                }
              }
            });
          }
          else if (signature.signature_reference_type === "company") {
            this.getCompany(signature.signature_reference_id).then((company) => {
              if (company) {
                signature.company = { ...company };
                this.getUser(company.company_manager_id).then((user) => {
                  if (user) {
                    signature.user = { ...user }
                  }
                }).catch((err) => {
                  console.log(err, 'this is getUset service')
                })
              }
            });
          }
          this.signatures.push(signature);
          setTimeout(() => {
            this.loading.signatures = false;
            this.rows = this.mapSignatures();
          }, 5000)
        }
      }
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
      if (popoverData) {
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
    if (this[callback]) {
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

  mapSignatures() {
    this.loading.signatures = false;
    return this.signatures.map((signature) => {
      const formattedSignature = {
        'Entity Type': signature.signature_reference_type,
        'Name': signature.user && signature.user.user_name,
        'Company': signature.company && signature.company.company_name,
        'GithubID': signature.user && signature.user.user_github_id,
        "LFID": signature.user && signature.user.lf_username,
        'Version': `v${signature.documentVersion}`,
        'Date': signature.date_modified
      }
      return formattedSignature
    })
  }

  onSearch($event) {
    this.filteredData = this.rows;
    let val = $event.value.trim().toLowerCase();
    console.log(val, 'this is va;')
    if (val.length > 0) {
      let colsAmt = this.columns.length;
      let keys = Object.keys(this.rows[0]);
      this.rows = this.filteredData.filter(function (item) {
        for (let i = 0; i < colsAmt; i++) {
          if (item[keys[i]] !== null && item[keys[i]] !== undefined && item[keys[i]].toString().toLowerCase().indexOf(val) !== -1 || !val) {
            // found match, return true to add to result set
            return true;
          }
        }
      });
    }
    else {
      this.rows = this.mapSignatures();
    }
  }
}
