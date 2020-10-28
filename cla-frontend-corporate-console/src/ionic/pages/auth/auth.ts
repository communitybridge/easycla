// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { AfterViewInit, Component } from '@angular/core';
import { NavController } from 'ionic-angular';
import { AuthService } from '../../services/auth.service';
import { EnvConfig } from '../../services/cla.env.utils';
import { LfxHeaderService } from '../../services/lfx-header.service';

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
export class AuthPage implements AfterViewInit {
  constructor(
    public navCtrl: NavController,
    public authService: AuthService,
    private lfxHeaderService: LfxHeaderService
  ) { }

  ngAfterViewInit() {

    this.authService.redirectRoot.subscribe((target) => {
      window.history.replaceState(null, null, window.location.pathname);
      this.navCtrl.setRoot('CompaniesPage');
    });

    this.authService.userProfile$.subscribe(user => {
      if (user !== undefined) {
        if (user) {
          this.lfxHeaderService.setUserInLFxHeader();
          this.navCtrl.setRoot('CompaniesPage');
        } else {
          window.open(EnvConfig['landing-page'], '_self');
        }
      }
    });
  }

  onClickToggle(toggle) {

  }
}


