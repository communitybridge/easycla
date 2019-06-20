// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import { Component } from '@angular/core';
import { NavController, NavParams, ViewController, IonicPage, ModalController } from 'ionic-angular';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { CincoService } from '../../services/cinco.service'

@IonicPage({
  segment: 'cla-corporate-whitelist-modal'
})
@Component({
  selector: 'cla-corporate-whitelist-modal',
  templateUrl: 'cla-corporate-whitelist-modal.html',
})
export class ClaCorporateWhitelistModal {
  domains: any;

  _form: FormGroup;
  submitAttempt: boolean = false;
  currentlySubmitting: boolean = false;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public viewCtrl: ViewController,
    private cincoService: CincoService,
    public modalCtrl: ModalController,
    private formBuilder: FormBuilder,
  ) {
    this.getDefaults();

    this._form = formBuilder.group({
      domain:['', Validators.compose([Validators.required])],
    });
  }

  ngOnInit() {

  }

  getDefaults() {
    this.domains = [
      "*@opensource.google.com",
      "*@kubernetes.io",
    ];
  }

  // ContactUpdateModal modal dismiss
  dismiss() {
    this.viewCtrl.dismiss();
  }

  addWhitelistDomain() {
    let domain = this._form.value.domain;
    this._form.patchValue({domain:''});
    this.domains.unshift(domain);
  }

}
