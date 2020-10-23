// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component, OnInit } from '@angular/core';
import { IonicPage, NavController } from 'ionic-angular';
import { AuthService } from '../../services/auth.service';
import { AUTH_ROUTE } from '../../services/auth.utils';
import { EnvConfig } from '../../services/cla.env.utils';

@IonicPage({
  name: 'LoginPage',
  segment: 'login',
})
@Component({
  selector: 'login',
  templateUrl: 'login.html',
})
export class LoginPage implements OnInit {

  constructor(
    public authService: AuthService,
    private navCtrl: NavController
  ) { }

  ngOnInit() {
    if (EnvConfig['lfx-header-enabled'] === "true") {
      window.open(EnvConfig['landing-page'], '_self');
    }
  }

  login() {
    if (this.authService.loggedIn) {
      this.navCtrl.setRoot('AllProjectsPage');
    } else {
      this.authService.login(AUTH_ROUTE);
    }
  }
}
