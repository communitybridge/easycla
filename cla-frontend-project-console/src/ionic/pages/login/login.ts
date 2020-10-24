// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { IonicPage, NavController } from 'ionic-angular';
import { AuthService } from '../../services/auth.service';
import { AUTH_ROUTE } from '../../services/auth.utils';

@IonicPage({
  name: 'LoginPage',
  segment: 'login',
})
@Component({
  selector: 'login',
  templateUrl: 'login.html',
})
export class LoginPage {

  constructor(
    public authService: AuthService,
    private navCtrl: NavController
  ) { }

  login() {
    if (this.authService.loggedIn) {
      this.navCtrl.setRoot('AllProjectsPage');
    } else {
      this.authService.login(AUTH_ROUTE);
    }
  }
}
