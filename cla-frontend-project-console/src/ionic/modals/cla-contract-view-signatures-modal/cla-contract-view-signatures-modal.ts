// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { PlatformLocation } from '@angular/common';
import {
  Events,
  IonicPage,
  ModalController,
  NavParams,
  PopoverController,
  ViewController,
} from 'ionic-angular';
import { ClaService } from '../../services/cla.service';
import { FormGroup, FormBuilder, Validators, FormControl } from '@angular/forms';

@IonicPage({
  segment: 'cla-contract-view-signatures-modal',
})
@Component({
  selector: 'cla-contract-view-signatures-modal',
  templateUrl: 'cla-contract-view-signatures-modal.html',
})
export class ClaContractViewSignaturesModal {
  selectedProject: any;
  claProjectId: string;
  claProjectName: string;
  loading: any;
  form: FormGroup;
  searchString: string;
  companies: any[];
  users: any[];
  filteredData: any[];
  data: any;
  // Pagination next/previous options
  limitPerPage: number = 100;
  resultCount: number = 0;
  nextKey = null;
  previousKeys = [];

  // Easy sort table variables
  columnData: any[] = [];
  column: any[] = [
    { head: 'Type', dataKey: 'signatureReferenceType' },
    { head: 'Name', dataKey: 'userName' },
    { head: 'Company Name', dataKey: 'companyName' },
    { head: 'GitHubID', dataKey: 'userGHID' },
    { head: 'LFID', dataKey: 'userLFID' },
    { head: 'Version', dataKey: 'version' },
    { head: 'Date', dataKey: 'signatureCreated' },
  ];

  constructor(
    public navParams: NavParams,
    private claService: ClaService,
    public viewCtrl: ViewController,
    public modalCtrl: ModalController,
    private popoverCtrl: PopoverController,
    public events: Events,
    private formBuilder: FormBuilder,
    private location: PlatformLocation
  ) {
    this.claProjectId = this.navParams.get('claProjectId');
    this.claProjectName = this.navParams.get('claProjectName');
    this.getDefaults();

    this.form = this.formBuilder.group({
      search: ['', Validators.compose([Validators.required, Validators.minLength(3)])],
      searchField: ['user'],
      fullMatch: [false],
    });

    this.location.onPopState(() => {
      this.viewCtrl.dismiss(false);
    });

    events.subscribe('modal:close', () => {
      this.dismiss();
    });
  }

  ngOnInit() {
    this.getSignatures();
  }

  get search(): FormControl {
    return <FormControl>this.form.get('search');
  }

  get searchField(): FormControl {
    return <FormControl>this.form.get('searchField');
  }

  get fullMatch(): FormControl {
    return <FormControl>this.form.get('fullMatch');
  }

  getDefaults() {
    this.data = {};
    this.loading = {
      signatures: true,
    };
  }

  filterDatatable() {
    if (this.form.valid) {
      this.searchString = this.search.value;
      this.getSignatures();
    }
  }

  resetFilter() {
    this.searchString = null;
    this.searchField.reset('user');
    this.fullMatch.setValue(false);
    this.search.reset();
    this.getSignatures();
  }

  // get all signatures
  getSignatures(lastKeyScanned = '') {
    this.loading.signatures = true;
    this.claService
      .getProjectSignaturesV3(
        this.claProjectId,
        this.limitPerPage,
        lastKeyScanned,
        this.searchString,
        this.searchField.value,
        null,
        this.fullMatch.value,
      )
      .subscribe((response) => {
        this.data = response;
        this.resultCount = this.data.resultCount;

        // Pagination Logic - add the key used to render this page to our previous keys
        if (lastKeyScanned) {
          this.previousKeys.push(lastKeyScanned);
        }
        // If we have a next key (usually we would unless there are no more records)
        if (this.data.lastKeyScanned) {
          this.nextKey = this.data.lastKeyScanned;
        } else {
          this.nextKey = null;
        }

        if (this.data && this.data.signatures) {
          this.columnData = this.data.signatures.map(e => ({
            ...e,
            signatureCreated: e.signatureCreated.split('T')[0],
            signatureReferenceType: this.getSignatureType(e),
            companyName: e.companyName ? this.trimCharacter(e.companyName, 25) : '',
            userName: e.userName ? this.trimCharacter(e.userName, 25) : '',
            icon: this.getSignatureType(e) === 'Company' ||
              this.getSignatureType(e) === 'Employee' ?
              { index: 0, iconName: 'briefcase' } :
              { index: 0, iconName: 'person' }
          }))
        } else {
          this.nextKey = null;
        }
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
      const previousLastScannedKey = this.previousKeys.length > 0 ? this.previousKeys.pop() : '';
      this.getSignatures(previousLastScannedKey);
    } else {
      this.getSignatures();
    }
  }

  previousButtonDisabled(): boolean {
    return this.previousKeys.length == 0;
  }

  nextButtonDisabled(): boolean {
    return this.nextKey == null && this.previousKeys.length >= 0;
  }

  signaturePopover(ev, signature) {
    let actions = {
      items: [
        {
          label: 'Details',
          callback: 'signatureDetails',
          callbackData: {
            signature: signature,
          },
        },
        {
          label: 'CLA',
          callback: 'signatureCla',
          callbackData: {
            signature: signature,
          },
        },
      ],
    };
    let popover = this.popoverCtrl.create('ActionPopoverComponent', actions);

    popover.present({
      ev: ev,
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



  dismiss() {
    this.viewCtrl.dismiss();
  }

  getSignatureType(signature: any): string {
    if (
      signature.signatureReferenceType === 'user' &&
      signature.signatureType === 'cla' &&
      signature.companyName == undefined
    ) {
      return 'Individual';
    } else if (
      signature.signatureReferenceType === 'user' &&
      signature.signatureType === 'cla' &&
      signature.companyName != undefined
    ) {
      return 'Employee';
    } else if (
      signature.signatureReferenceType === 'company' &&
      signature.signatureType === 'ccla' &&
      signature.companyName != undefined
    ) {
      return 'Company';
    } else {
      return 'Unknown';
    }
  }

  trimCharacter(text, length) {
    if (text !== undefined && text !== null) {
      return text.length > length ? text.substring(0, length) + '...' : text;
    }
  }
}
