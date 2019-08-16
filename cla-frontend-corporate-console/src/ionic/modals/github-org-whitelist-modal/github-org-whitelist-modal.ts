// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import {Component} from '@angular/core';
import {Events, IonicPage, NavParams, ViewController} from 'ionic-angular';
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
  signatureId: string;

  constructor(
    public navParams: NavParams,
    public viewCtrl: ViewController,
    public claService: ClaService,
    public events: Events
  ) {
    this.corporateClaId = this.navParams.get('corporateClaId');
    this.companyId = this.navParams.get('companyId');
    this.signatureId = this.navParams.get('signatureId');
    // console.log('constructor - loaded signature id: ' + this.signatureId);

    events.subscribe('modal:close', () => {
      this.dismiss();
    });
  }

  ngOnInit() {
    this.getGitHubOrgWhiteList();
  }

  getGitHubOrgWhiteList() {
    //console.log('getGitHubOrgWhiteList - using signature id: ' + this.signatureId);
    this.claService.getGithubOrganizationWhitelistEntries(this.signatureId)
      .subscribe(organizations => {
        this.organizations = organizations;
      });
  }

  addGitHubOrgWhiteList(organizationId) {
    //console.log('addGitHubOrgWhiteList - using signature id: ' + this.signatureId +
    //  ' with organization id: ' + organizationId);
    this.claService.addGithubOrganizationWhitelistEntry(this.signatureId, organizationId)
      .subscribe(() => this.getGitHubOrgWhiteList());
  }

  removeGitHubOrgWhiteList(organizationId) {
    //console.log('removeGitHubOrgWhiteList - using signature id: ' + this.signatureId +
    //  ' with organization id: ' + organizationId);
    this.claService.removeGithubOrganizationWhitelistEntry(this.signatureId, organizationId)
      .subscribe(() => this.getGitHubOrgWhiteList());
  }

  dismiss() {
    this.viewCtrl.dismiss();
  }

  connectGithub() {
    this.claService.githubLogin(this.companyId, this.corporateClaId);
  }
}
