// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import {Component} from '@angular/core';
import {DatePipe} from '@angular/common';
import {
  Events,
  IonicPage,
  ModalController,
  NavController,
  NavParams,
  PopoverController,
  ViewController,
} from 'ionic-angular';
import {ClaService} from '../../services/cla.service'
import {SortService} from "../../services/sort.service";
import {KeycloakService} from "../../services/keycloak/keycloak.service";
import {RolesService} from "../../services/roles.service";
import {ColumnMode, SortType} from "@swimlane/ngx-datatable";


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
  claProjectName: string;

  ColumnMode = ColumnMode;
  SortType = SortType;

  loading: any;
  //sort: any;
  columns: any[];
  rows: any[];

  companies: any[];
  users: any[];
  filteredData: any[];
  data: any;
  page: any;

  // Pagination next/previous options
  nextKey: string;
  previousKeys: any[];

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
    this.claProjectName = this.navParams.get('claProjectName');
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
      pageNumber: 0
    };

    // Pagination initialization
    this.nextKey = null;
    this.previousKeys = [];

    this.data = {};
    this.loading = {
      signatures: true
    };

    /*
    this.sort = {
      Type: {
        arrayProp: 'Type',
        sortType: 'text',
        sort: null,
      },
      Name: {
        arrayProp: 'Name',
        sortType: 'text',
        sort: null,
      },
      Company: {
        arrayProp: 'Company',
        sortType: 'text',
        sort: null,
      },
      GitHubID: {
        arrayProp: 'GitHubID',
        sortType: 'number',
        sort: null,
      },
      LFID: {
        arrayProp: 'LFID',
        sortType: 'number',
        sort: null,
      },
      Version: {
        arrayProp: 'documentVersion',
        sortType: 'semver',
        sort: null,
      },
      Date: {
        arrayProp: 'Date',
        sortType: 'date',
        sort: null,
      },
    };
     */

    this.filteredData = this.rows;
    this.columns = [
      {prop: 'Type'},
      {prop: 'Name'},
      {prop: 'Company'},
      {prop: 'GitHubID'},
      {prop: 'LFID'},
      {prop: 'Version'},
      {prop: 'Date'}
    ];
  }

  async getUser(signatureReferenceId) {
    return await this.claService.getUser(signatureReferenceId).toPromise();
  }

  async getCompany(referenceId) {
    return await this.claService.getCompany(referenceId).toPromise();
  }

  // get all signatures
  getSignatures(lastKeyScanned = '') {
    this.loading.signatures = true;
    this.claService.getProjectSignaturesV3(this.claProjectId, 100, lastKeyScanned).subscribe((response) => {
      this.data = response;

      // Pagination Logic - add the key used to render this page to our previous keys
      this.previousKeys.push(lastKeyScanned);
      // If we have a next key (usually we would unless there are no more records)
      if (this.data.lastKeyScanned) {
        this.nextKey = this.data.lastKeyScanned;
      } else {
        this.nextKey = null;
      }

      this.page.totalCount = this.data.resultCount;
      this.rows = this.mapSignatures(this.data.signatures);
      this.loading.signatures = false;
    });
  }

  getNextPage() {
    if (this.nextKey) {
      this.getSignatures(this.nextKey);
    } else {
      this.getSignatures();
    }
  }

  getPreviousPage() {
    if (this.previousKeys.length > 0) {
      // Since the most recent previous key is the current page - we want to go one more back to the one before
      // so we pop two and use the second key as the key to use to render the page
      this.previousKeys.pop();
      const previousLastScannedKey = this.previousKeys.pop();
      this.getSignatures(previousLastScannedKey);
    } else {
      this.getSignatures();
    }
  }

  previousButtonDisabled(): boolean {
    return !(this.previousKeys.length > 1);
  }

  nextButtonDisabled(): boolean {
    return (this.nextKey == null && this.previousKeys.length >= 0);
  }

  previousButtonColor(): string {
    if (this.previousKeys.length <= 1) {
      return 'gray';
    } else {
      return 'secondary';
    }
  }

  nextButtonColor(): string {
    if (this.nextKey == null && this.previousKeys.length >= 0) {
      return 'gray';
    } else {
      return 'secondary';
    }
  }

  /**
   * Helper function to dump the pagination details.
   */
  debugShowPaginationReport() {
    console.log('NextKey: ' + this.nextKey);
    console.log('PreviousKeys:');
    console.log(this.previousKeys);
    console.log('------------------------------');
  }

  /**
   * Default sorting.
   *
   * @param prop
   */
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

  mapSignatures(signatures: any[]) {
    // If no records
    if (signatures == null || signatures.length == 0) {
      return [];
    } else {
      return signatures && signatures.map((signature) => {
        let date = this.datePipe.transform(signature.signatureCreated, 'yyyy-MM-dd');
        return {
          /**
           * | Type                   | Reference Type | Signature Type | Company Name |
           * |------------------------|----------------|----------------|--------------|
           * | ICLA (individual icon) | user           | cla            | empty        |
           * | CCLA (employee icon)   | user           | cla            | not empty    |
           * | CCLA (company icon)    | company        | ccla           | not empty    |
           */
          'Type': this.getSignatureType(signature),
          'Name': signature.userName && signature.userName,
          'Company': signature.companyName && signature.companyName,
          'GitHubID': signature.userGHID && signature.userGHID,
          "LFID": signature.userLFID && signature.userLFID,
          'Version': `v${signature.version}`,
          'Date': date
        }
      })
    }
  }

  getSignatureType(signature: any): string {
    if (signature.signatureReferenceType === 'user' &&
      signature.signatureType === 'cla' &&
      signature.companyName == undefined) {
      return 'individual';
    } else if (signature.signatureReferenceType === 'user' &&
      signature.signatureType === 'cla' &&
      signature.companyName != undefined) {
      return 'employee';
    } else if (signature.signatureReferenceType === 'company' &&
      signature.signatureType === 'ccla' &&
      signature.companyName != undefined) {
      return 'company';
    } else {
      return 'unknown';
    }
  }
}
