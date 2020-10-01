// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component, Input } from '@angular/core';
import { AuthService } from 'src/app/core/services/auth.service';
import { AppSettings } from 'src/app/config/app-settings';

@Component({
  selector: 'app-cla-console-section',
  templateUrl: './cla-console-section.component.html',
  styleUrls: ['./cla-console-section.component.scss']
})
export class ClaConsoleSectionComponent {
  @Input() consoleMetadata: any;

  constructor(
    private authService: AuthService
  ) { }

  onClickProceed(type: string) {
    this.authService.login(type);
  }

  onClickRequestAccess() {
    window.open(AppSettings.REQUEST_ACCESS_LINK, '_blank');
  }

  onClickLearnMore() {
    window.open(AppSettings.LEARN_MORE, '_blank');
  }
}
