// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { NavController, ModalController, NavParams, IonicPage } from 'ionic-angular';
import { CincoService } from '../../../services/cinco.service';
import { KeycloakService } from '../../../services/keycloak/keycloak.service';
import { SortService } from '../../../services/sort.service';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { CheckboxValidator } from  '../../../validators/checkbox';
import { RolesService } from '../../../services/roles.service';
import { Restricted } from '../../../decorators/restricted';

@Restricted({
  roles: ['isAuthenticated', 'isPmcUser'],
})
// @IonicPage({
//   segment: 'project/:projectId/member/:memberId/details'
// })
@Component({
  selector: 'member-details',
  templateUrl: 'member-details.html',
})
export class MemberDetailsPage {
  projectId: any;
  memberId: any;
  details: any;

  loading: any;
  sort: any;

  contracts: any;

  iclaUploadInfo: any;
  cclaUploadInfo: any;

  form: FormGroup;
  submitAttempt: boolean = false;
  currentlySubmitting: boolean = false;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    private cincoService: CincoService,
    private sortService: SortService,
    public modalCtrl: ModalController,
    private formBuilder: FormBuilder,
    private keycloak: KeycloakService,
    public rolesService: RolesService,
  ) {
    this.projectId = navParams.get('projectId');
    this.memberId = navParams.get('memberId');
    this.form = formBuilder.group({
      confirm: [false, Validators.compose([CheckboxValidator.isChecked])],
    });
    this.getDefaults();
  }

  ngOnInit() {
    this.getMemberDetails();
  }

  getMemberDetails() {
    // this.cincoService.getMemberDetails(this.projectId).subscribe(response => {
    //   if (response) {
    //     this.details = response;
    //   }
    //   this.loading.details = false;
    // });
    setTimeout((function(){
      this.details = {
        product: "Gold Membership",
        startTier: "",
        endTier: "",
        level: "End User",
        startDate: "1/1/2017",
        endDate: "12/31/2017",
        poNumber: "0123456789",
        description: "Gold Membership, Sample Project",
        netsuiteMemo:"400 character memo field goes here. This is a memo field where someone might choose to write some sort of memo or other information that relates to specifically how the member should receive their invoice, for instance one time Jen had to write out a detailed fee schedule.",
        amount: "$500,000",
      };
      this.loading.details = false;
    }).bind(this),2000);

  }

  getDefaults() {
    this.loading = {
      details: true,
    };
    this.details = {
      product: "",
      startTier: "",
      endTier: "",
      level: "",
      startDate: "",
      endDate: "",
      poNumber: "",
      description: "",
      netsuiteMemo:"",
      ammount:"",
    };

  }

  submit() {

  }
}
