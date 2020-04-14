// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component, ViewChild } from '@angular/core';
import { DatePipe, PlatformLocation } from '@angular/common';
import {
  Events,
  IonicPage,
  ModalController,
  NavController,
  NavParams,
  PopoverController,
  ViewController,
} from 'ionic-angular';
import { ClaService } from '../../services/cla.service';
import { SortService } from '../../services/sort.service';
import { KeycloakService } from '../../services/keycloak/keycloak.service';
import { RolesService } from '../../services/roles.service';
import { FormGroup, FormBuilder, Validators, FormControl } from '@angular/forms';

@IonicPage({
  segment: 'cla-contract-view-companies-signatures-modal',
})
@Component({
  selector: 'cla-contract-view-companies-signatures-modal',
  templateUrl: 'cla-contract-view-companies-signatures-modal.html',
})
export class ClaContractViewCompaniesSignaturesModal {
  @ViewChild('companiesTable') table: any;
  selectedProject: any;
  claProjectId: string;
  claProjectName: string;
  loading: any;
  form: FormGroup;
  searchString: string = '';
  companies: any[];
  users: any[];
  data: any;
  page: any;
  // Pagination next/previous options
  limitPerPage: number = 100;
  resultCount: number = 0;
  nextKey = null;
  previousKeys = [];
  errorMsg = '';

  // Easy sort table variables
  columnData = [];
  column = [{ head: 'Company', dataKey: 'companyName' }, { head: 'Version', dataKey: 'version' }, { head: 'Date Signed', dataKey: 'signatureCreated' }];
  childColumn = [{ head: 'Name', dataKey: 'username' }, { head: 'Email', dataKey: 'lfEmail' }, { head: 'username/LFID', dataKey: 'lfUsername' }];

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    private claService: ClaService,
    public viewCtrl: ViewController,
    public modalCtrl: ModalController,
    private popoverCtrl: PopoverController,
    public rolesService: RolesService,
    public events: Events,
    private formBuilder: FormBuilder,
    private location: PlatformLocation
  ) {
    this.claProjectId = this.navParams.get('claProjectId');
    this.claProjectName = this.navParams.get('claProjectName');
    this.getDefaults();

    this.form = this.formBuilder.group({
      search: ['', Validators.compose([Validators.required, Validators.minLength(3)])],
      searchField: ['company'],
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
    this.page = {
      pageNumber: 0,
    };

    // Pagination initialization
    this.nextKey = null;
    this.previousKeys = [];

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
    this.searchField.reset('company');
    this.fullMatch.setValue(false);
    this.search.reset();
    this.getSignatures();
  }

  // get all signatures
  getSignatures(lastKeyScanned = '') {
    this.errorMsg = '';
    this.loading.signatures = true;
    this.claService
      .getProjectSignaturesV3(
        this.claProjectId,
        this.limitPerPage,
        lastKeyScanned,
        this.searchString,
        this.searchField.value,
        'ccla',
        this.fullMatch.value,
      )
      .subscribe((response) => {
        this.data = response;

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
            companyName: this.trimCharacter(e.companyName, 30)
          }))
        } else {
          this.columnData = []
        }

        this.page.totalCount = this.data.resultCount;
        this.loading.signatures = false;
      }, err => {
        this.errorMsg = 'Could not retrieve company signature data. Please try again';
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
      return 'individual';
    } else if (
      signature.signatureReferenceType === 'user' &&
      signature.signatureType === 'cla' &&
      signature.companyName != undefined
    ) {
      return 'employee';
    } else if (
      signature.signatureReferenceType === 'company' &&
      signature.signatureType === 'ccla' &&
      signature.companyName != undefined
    ) {
      return 'company';
    } else {
      return 'unknown';
    }
  }

  trimCharacter(text, length) {
    return text.length > length ? text.substring(0, length) + '...' : text;
  }
}
