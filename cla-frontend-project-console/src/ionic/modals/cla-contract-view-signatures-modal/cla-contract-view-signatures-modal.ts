// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { DatePipe } from '@angular/common';
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
import { ColumnMode, SortType } from '@swimlane/ngx-datatable';
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

  ColumnMode = ColumnMode;
  SortType = SortType;

  loading: any;
  //sort: any;
  columns: any[];
  rows: any[];

  form: FormGroup;
  searchString: string;

  companies: any[];
  users: any[];
  filteredData: any[];
  data: any;
  page: any;

  // Pagination next/previous options
  nextKey: string;
  previousKeys: any[];

  // Easy sort table variables
  columnData = [];
  column = [
    {head:'Type', dataKey:'signatureReferenceType'}, 
    {head:'Name', dataKey:'userName,'}, 
    {head:'Company Name', dataKey:'companyName'}, 
    {head:'GitHubID', dataKey:'userGHID'},
    {head:'LFID', dataKey:'userLFID'},
    {head:'Version', dataKey:'version'},
    {head:'Date', dataKey:'signatureCreated'},
  ];

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
    public events: Events,
    private formBuilder: FormBuilder,
  ) {
    this.claProjectId = this.navParams.get('claProjectId');
    this.claProjectName = this.navParams.get('claProjectName');
    this.getDefaults();

    this.form = this.formBuilder.group({
      search: ['', Validators.compose([Validators.required, Validators.minLength(3)])],
      searchField: ['user'],
      fullMatch: [false],
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

  async getUser(signatureReferenceId) {
    return await this.claService.getUser(signatureReferenceId).toPromise();
  }

  async getCompany(referenceId) {
    return await this.claService.getCompany(referenceId).toPromise();
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
        100,
        lastKeyScanned,
        this.searchString,
        this.searchField.value,
        null,
        this.fullMatch.value,
      )
      .subscribe((response) => {
        this.data = response;
        this.page.totalCount = this.data.resultCount;
        console.log('data', this.data);
        if(this.data && this.data.signatures){
          this.columnData = this.data.signatures.map(e => ({
            ...e, 
            signatureCreated:  e.signatureCreated.split('T')[0],
            signatureReferenceType: this.getSignatureType(e),
            icon: this.getSignatureType(e) === 'Company' || 
            this.getSignatureType(e) === 'Employee' ? 
            {index: 0, iconName: 'briefcase'}: 
            {index: 0, iconName: 'person'}
          }))
        } else {
          this.columnData = []
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
    return this.nextKey == null && this.previousKeys.length >= 0;
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

  signatureDetails(data) {
    console.log('signature details');
  }

  signatureCla(data) {
    console.log('signature cla');
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
}
