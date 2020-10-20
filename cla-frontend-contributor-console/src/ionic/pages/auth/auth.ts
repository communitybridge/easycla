// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Location } from '@angular/common';
import { Component } from '@angular/core';
import { NavController, NavParams } from 'ionic-angular';
import { AuthService } from '../../services/auth.service';

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
    private location: Location
  ) {
    this.gerritId = localStorage.getItem('gerritId');
    this.claType = localStorage.getItem('gerritClaType');
  }

  ngOnInit() {
    this.authService.redirectRoot.subscribe((target) => {
      if (this.authService.loggedIn) {
        window.history.replaceState(null, null, window.location.pathname);
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

}
