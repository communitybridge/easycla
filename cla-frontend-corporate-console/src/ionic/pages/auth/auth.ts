// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component, OnInit } from '@angular/core';
import { NavParams } from 'ionic-angular';
import { NavController } from 'ionic-angular';
import { AuthService } from '../../services/auth.service';
import { EnvConfig } from '../../services/cla.env.utils';
import { RolesService } from '../../services/roles.service';

/**
 * Generated class for the AuthPage page.
 *
 * See https://ionicframework.com/docs/components/#navigation for more info on
 * Ionic pages and navigation.
 */
@Component({
  selector: 'auth-page',
  templateUrl: 'auth.html'
})
export class AuthPage implements OnInit {
  userRoles: any;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public authService: AuthService,
    public rolesService: RolesService
  ) { }

  ngOnInit() {
    setTimeout(() => {
      console.log(this.authService.loggedIn);
    if (this.authService.loggedIn) {
      this.navCtrl.setRoot('CompaniesPage');
    } else {
      this.redirectToLogin();
    }
    }, 2000);
  }

  redirectToLogin() {
    if (EnvConfig['lfx-header-enabled'] === "true") {
      window.open(EnvConfig['landing-page'], '_self');
    } else {
      this.navCtrl.setRoot('LoginPage');
    }
  }
}
