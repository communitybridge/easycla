// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { IonicPage } from 'ionic-angular';
import { AuthService } from '../../services/auth.service';

@IonicPage({
  name: 'Authorize',
  segment: 'authorize',
})
@Component({
  selector: 'authorize',
  templateUrl: 'authorize.html',
})
export class Authorize {

  constructor(
    public authService: AuthService
  ) {
    console.log('Entered');
    this.authService.login();
  }

  onClickToggle(hasExpanded) {
  }
}
