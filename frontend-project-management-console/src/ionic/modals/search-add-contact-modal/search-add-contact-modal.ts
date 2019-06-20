// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import { Component, ChangeDetectorRef } from '@angular/core';
import { NavController, NavParams, ModalController, ViewController, AlertController, IonicPage } from 'ionic-angular';
import { CincoService } from '../../services/cinco.service';
import { SortService } from '../../services/sort.service';

@IonicPage({
  segment: 'search-add-contact-modal'
})
@Component({
  selector: 'search-add-contact-modal',
  templateUrl: 'search-add-contact-modal.html',
  providers: [
    CincoService,
    SortService,
  ]
})
export class SearchAddContactModal {
  projectId: string;
  memberId: string;
  org: any;
  enteredEmail: string;
  organizationContacts: any;
  orgContactRoles: any;
  loading: any;
  sort: any;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public modalCtrl: ModalController,
    public viewCtrl: ViewController,
    public alertCtrl: AlertController,
    private changeDetectorRef: ChangeDetectorRef,
    private cincoService: CincoService,
    private sortService: SortService,
  ) {
    this.getDefaults();
    this.projectId = this.navParams.get('projectId');
    this.memberId = this.navParams.get('memberId');
    this.org = this.navParams.get('org');
  }

  getDefaults() {
    this.loading = {
      contacts: true,
    }
    this.orgContactRoles = {};
    this.sort = {
      type: {
        arrayProp: 'type',
        sortType: 'text',
        sort: null,
      },
      name: {
        arrayProp: 'givenName',
        sortType: 'text',
        sort: null,
      },
      email: {
        arrayProp: 'email',
        sortType: 'text',
        sort: null,
      },
    };
  }

  ngOnInit() {
    let orgId = this.org.id;
    this.getOrgContactRoles();
    this.getOrganizationContacts(orgId);
  }

  getOrganizationContacts(orgId) {
    this.cincoService.getOrganizationContacts(orgId).subscribe(response => {
      if(response) {
        this.organizationContacts = response;
        this.loading.contacts = false;
      }
    });
  }

  getOrgContactRoles() {
    this.cincoService.getOrganizationContactTypes().subscribe(response => {
      if(response) {
        this.orgContactRoles = response;
      }
    });
  }

  // ContactUpdateModal modal dismiss
  dismiss() {
    this.viewCtrl.dismiss();
  }

  addContact(contact) {
    let modal = this.modalCtrl.create('ContactUpdateModal', {
      projectId: this.projectId,
      memberId: this.memberId,
      org: this.org,
      contact: contact,
    });
    modal.onDidDismiss(data => {
      // A refresh of data anytime the modal is dismissed
      let orgId = this.org.id;
      this.getOrganizationContacts(orgId);
    });
    modal.present();
  }

  filterContactsByEmail() {

  }

  sortContacts(prop) {
    this.sortService.toggleSort(
      this.sort,
      prop,
      this.organizationContacts,
    );
  }

}
