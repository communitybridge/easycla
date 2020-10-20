// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { IonicPage, NavController } from 'ionic-angular';
import { RolesService } from '../../services/roles.service';
import { AuthService } from '../../services/auth.service';
import { AUTH_ROUTE } from '../../services/auth.utils';

@IonicPage({
  name: 'LoginPage',
  segment: 'login'
})
@Component({
  selector: 'login',
  templateUrl: 'login.html'
})
export class LoginPage {
  canAccess: boolean;

  constructor(
    public navCtrl: NavController,
    public rolesService: RolesService,
    public authService: AuthService
  ) { }

  login() {
    if (this.authService.loggedIn) {
      const gerritId = localStorage.getItem('gerritId');
      const claType = localStorage.getItem('gerritClaType');
      if (claType == 'ICLA') {
        this.navCtrl.setRoot('ClaGerritIndividualPage', { gerritId: gerritId });
      } else if (claType == 'CCLA') {
        this.navCtrl.setRoot('ClaGerritCorporatePage', { gerritId: gerritId });
      } else {
        console.log('Invalid URL : Login flow work only for Gerrit.');
      }
      localStorage.removeItem('gerritId');
      localStorage.removeItem('gerritClaType');
    } else {
      this.authService.login(AUTH_ROUTE);
    }
  }
}
