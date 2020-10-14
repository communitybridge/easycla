// Copyright The Linux Foundation and each contributor to CommunityBridge.	
// SPDX-License-Identifier: MIT	

import { Component } from '@angular/core';
import { NavController } from 'ionic-angular';
import { IonicPage } from 'ionic-angular';
import { AuthService } from '../../services/auth.service';
import { AUTH_REDIRECT_STATE } from '../../services/auth.utils';

@IonicPage({
  name: 'LoginPage',
  segment: 'login'
})
@Component({
  selector: 'login-page',
  templateUrl: 'login-page.html'
})
export class LoginPage {
  constructor(
    private authService: AuthService,
    private navCtrl: NavController
  ) { }

  login() {
    if (this.authService.loggedIn) {
      this.navCtrl.setRoot('CompaniesPage');
    } else {
      this.authService.login(AUTH_REDIRECT_STATE);
    }
  }
}