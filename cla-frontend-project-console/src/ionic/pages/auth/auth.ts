// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component, OnInit } from '@angular/core';
import { NavController } from 'ionic-angular';
import { AuthService } from '../../services/auth.service';
import { EnvConfig } from '../../services/cla.env.utils';

/**
 * Generated class for the AuthPage page.
 *
 * See https://ionicframework.com/docs/components/#navigation for more info on
 * Ionic pages and navigation.
 */

@Component({
  selector: 'page-auth',
  templateUrl: 'auth.html'
})
export class AuthPage implements OnInit {
  constructor(
    public navCtrl: NavController,
    public authService: AuthService
  ) { }

  ngOnInit() {
    this.authService.redirectRoot.subscribe((target) => {
      window.history.replaceState(null, null, window.location.pathname);
      this.navCtrl.setRoot('AllProjectsPage');
    });

    this.authService.checkSession.subscribe((loggedIn) => {
      if (loggedIn) {
        this.navCtrl.setRoot('AllProjectsPage');
      } else {
        this.redirectToLogin();
      }
    });

  }

  redirectToLogin() {
    if (EnvConfig['lfx-header-enabled'] === "true") {
      window.open(EnvConfig['landing-page'], '_self');
    } else {
      this.navCtrl.setRoot('LoginPage');
    }
  }
}
