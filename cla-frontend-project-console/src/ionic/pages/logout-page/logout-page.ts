// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import {Component} from '@angular/core';
import {IonicPage, NavController} from 'ionic-angular';
import {AuthService} from "../../services/auth.service";

@IonicPage({
  name: 'LogoutPage',
  segment: 'logout'
})
@Component({
  selector: 'logout-page',
  templateUrl: 'logout-page.html'
})
export class LogoutPage {
  constructor(
    public navCtrl: NavController,
    public authService: AuthService
  ) {
  }

  ngOnInit() {
    // Will redirect user back to the login page after logging out of auth0
    this.authService.logout();
  }
}
