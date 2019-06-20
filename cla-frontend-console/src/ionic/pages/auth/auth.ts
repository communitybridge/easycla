// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import { Component } from '@angular/core';
import { NavController, NavParams } from 'ionic-angular';
import { AuthService } from "../../services/auth.service";
import { RolesService } from "../../services/roles.service";

 @Component({
  selector: 'page-auth',
  templateUrl: 'auth.html',
})
export class AuthPage {

  userRoles: any;
  gerritId: string;
  claType: string;

  constructor(public navCtrl: NavController, 
    public navParams: NavParams,
    public authService: AuthService, 
    public rolesService: RolesService) {
    this.gerritId = localStorage.getItem("gerritId");
    this.claType = localStorage.getItem("gerritClaType");
  }

  ionViewDidEnter() {
    setTimeout(() => {
      this.rolesService.getUserRolesPromise().then((userRoles) => {
        if (userRoles.isAuthenticated) { 
            if(this.claType == "ICLA") {
              this.navCtrl.setRoot('ClaGerritIndividualPage', {gerritId: this.gerritId});
            }
            else if(this.claType == "CCLA") { 
              this.navCtrl.setRoot('ClaGerritCorporatePage', {gerritId: this.gerritId});
            }
            localStorage.removeItem("gerritId");
            localStorage.removeItem("gerritClaType");

        }
        else { 
          this.navCtrl.setRoot('loginPage');
        }
      });
    }, 2000); 
  }

  

}
