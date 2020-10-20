// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { NavController, NavParams } from 'ionic-angular';
import { AuthService } from '../../services/auth.service';
import { EnvConfig } from '../../services/cla.env.utils';
import { RolesService } from '../../services/roles.service';

@Component({
  selector: 'page-auth',
  templateUrl: 'auth.html'
})
export class AuthPage {
  userRoles: any;
  gerritId: string;
  claType: string;

  constructor(
    public navCtrl: NavController,
    public authService: AuthService,
  ) {
    this.gerritId = localStorage.getItem('gerritId');
    this.claType = localStorage.getItem('gerritClaType');
  }

  ngOnInit() {
    this.authService.redirectRoot.subscribe((target) => {
      if (this.authService.loggedIn) {
        if (this.claType == 'ICLA') {
          this.navCtrl.setRoot('ClaGerritIndividualPage', { gerritId: this.gerritId });
        } else if (this.claType == 'CCLA') {
          this.navCtrl.setRoot('ClaGerritCorporatePage', { gerritId: this.gerritId });
        }
      } else {
        console.log('Redirect to login');
      }
    });
  }

  redirectAsPerType() {
    let url = `${window.location.origin}`;
    if (this.claType == 'ICLA') {
      url += '/#/cla/gerrit/project/' + this.gerritId + '/individual';
    } else if (this.claType == 'CCLA') {
      url += '/#/cla/gerrit/project/' + this.gerritId + '/corporate';
    }
    window.open(url, '_self');
  }
}
