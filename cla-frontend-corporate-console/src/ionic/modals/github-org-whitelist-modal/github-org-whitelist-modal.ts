// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import {Component} from '@angular/core';
import {NavParams, ViewController, IonicPage, Events} from 'ionic-angular';
import {ClaService} from '../../services/cla.service';

@IonicPage({
  segment: 'github-org-whitelist-modal'
})
@Component({
  selector: 'github-org-whitelist-modal',
  templateUrl: 'github-org-whitelist-modal.html'
})
export class GithubOrgWhitelistModal {
  responseErrors: string[] = [];
  organizations: any[] = [];
  corporateClaId: any;
  companyId: any;

  constructor(
    public navParams: NavParams,
    public viewCtrl: ViewController,
    public claService: ClaService,
    public events: Events
  ) {
    this.corporateClaId = this.navParams.get('corporateClaId');
    this.companyId = this.navParams.get('companyId');

    events.subscribe('modal:close', () => {
      this.dismiss();
    });
  }

  ngOnInit() {
    this.getOrgWhitelist();
  }

  getOrgWhitelist () {
    this.claService.getGithubOrganizationWhitelist(this.companyId, this.corporateClaId)
      .subscribe(organizations => {
        this.organizations = organizations;
      })
  }

  addOrganization(organizationId) {
    this.claService.addGithubOrganizationWhitelist(this.companyId, this.corporateClaId, organizationId)
      .subscribe(() => this.getOrgWhitelist());
  }

  removeOrganization(organizationId) {
    this.claService.removeGithubOrganizationWhitelist(this.companyId, this.corporateClaId, organizationId)
      .subscribe(() => this.getOrgWhitelist());
  }

  dismiss() {
    this.viewCtrl.dismiss();
  }

  connectGithub (){
    this.claService.githubLogin(this.companyId, this.corporateClaId);
  }
}
