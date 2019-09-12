// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { DatePipe } from '@angular/common';
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
  rows: any[];
  selectedSize: any[];
  allSizes: any[];

  companies: any[];
  users: any[];
  filteredData: any[];
  data: any;
  page: any;
  lastScannedKey: any[];
  previousLastScannedKeys: any[];
  pageSize: any

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    private claService: ClaService,
    private sortService: SortService,
    public viewCtrl: ViewController,
    public modalCtrl: ModalController,
    private popoverCtrl: PopoverController,
    private keycloak: KeycloakService,
    private datePipe: DatePipe,
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
    this.page = {
      size: 10,
      pageNumber: 0

    }
    this.pageSize = 50
    this.lastScannedKey = [];
    this.previousLastScannedKeys = [];
    this.data = {};
    this.loading = {
      signatures: true
    };
    this.allSizes = [
      1,
      10,
      50,
      100,
      150,
      200
    ]
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

  sizeSelectedChanged() {
    this.pageSize = this.selectedSize;
    this.getSignatures();

  }

  // get all signatures
  getSignatures(lastKeyScanned = "") {
    this.claService.getProjectSignaturesV3(this.claProjectId, this.pageSize, lastKeyScanned).subscribe((response) => {
      this.data = response;
      if (this.data.lastKeyScanned) {
        // push next keys to a stack
        this.lastScannedKey.push(this.data.lastKeyScanned)
      }
      this.page.totalCount = this.data.resultCount
      this.signatures = this.data.signatures;
      this.rows = this.mapSignatures();
      this.loading.signatures = false;
    });
  }

  getNextPage() {
    this.loading.signatures = true;
    let lastKeyScanned = this.lastScannedKey.pop();
    this.previousLastScannedKeys.push(lastKeyScanned);
    this.getSignatures(lastKeyScanned)
  }

  getPreviousPage() {
    this.loading.signatures = true;
    const previousLastScannedKeys = this.previousLastScannedKeys.shift()
    this.getSignatures()
  }

  sortMembers(prop) {
   
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
    return this.signatures && this.signatures.map((signature) => {
      let date = this.datePipe.transform(signature.signatureCreated, 'yyyy-MM-dd');
      const formattedSignature = {
        'Entity Type': signature.signatureReferenceType,
        'Name': signature.userName && signature.userName,
        'Company': signature.companyName && signature.companyName,
        'GithubID': signature.userGHID && signature.userGHID,
        "LFID": signature.LFID && signature.LFID,
        'Version': `v${signature.version}`,
        'Date': date
      }
      return formattedSignature
    })
  }

  onSearch($event) {
    this.filteredData = this.rows;
    let val = $event.value.trim().toLowerCase();
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
