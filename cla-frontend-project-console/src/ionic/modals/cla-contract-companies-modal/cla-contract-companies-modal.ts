// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { Events, IonicPage, ModalController, NavController, NavParams, ViewController } from 'ionic-angular';
import { ClaService } from '../../services/cla.service';

@IonicPage({
  segment: 'cla-contract-companies-modal'
})
@Component({
  selector: 'cla-contract-companies-modal',
  templateUrl: 'cla-contract-companies-modal.html'
})
export class ClaContractCompaniesModal {
  loading: any;

  claProjectId: string;
  companies: any;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public viewCtrl: ViewController,
    private claService: ClaService,
    public modalCtrl: ModalController,
    public events: Events
  ) {
    this.getDefaults();

    events.subscribe('modal:close', () => {
      this.dismiss();
    });
  }

  ngOnInit() {
    this.getProjectCompanies();
  }

  getDefaults() {
    this.claProjectId = this.navParams.get('claProjectId');
    this.companies = [];
    this.loading = {
      initial: true,
      companies: true
    };
  }

  getProjectCompanies() {
    this.loading.companies = true;
    this.claService.getProjectCompanies(this.claProjectId).subscribe(companies => {
      this.loading.initial = false;
      this.loading.companies = false;
      this.companies = companies;
    });
  }

  openClaCorporateMemberOptionsModal() {
    let modal = this.modalCtrl.create('ClaCorporateMemberOptionsModal');
    modal.present();
  }

  openClaCorporateManagementPage(companyId) {
    console.log(companyId);
  }

  // ContactUpdateModal modal dismiss
  dismiss() {
    this.viewCtrl.dismiss();
  }
}
