// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component, Input } from '@angular/core';
import { AuthService } from 'src/app/core/services/auth.service';
import { AppSettings } from 'src/app/config/app-settings';
import { StorageService } from 'src/app/core/services/storage.service';
import { EnvConfig } from 'src/app/config/cla-env-utils';

@Component({
  selector: 'app-cla-console-section',
  templateUrl: './cla-console-section.component.html',
  styleUrls: ['./cla-console-section.component.scss']
})
export class ClaConsoleSectionComponent {
  @Input() consoleMetadata: any;

  constructor(
    private authService: AuthService,
    private storageService: StorageService
  ) { }

  onClickProceed(type: string) {
    this.storageService.setItem('type', type);
    const redirectConsole = (type === 'Projects') ? AppSettings.PROJECT_CONSOLE_LINK : AppSettings.CORPORATE_CONSOLE_LINK;
    window.open(EnvConfig.default[redirectConsole] + '#/authorize', '_self');
  }

  onClickRequestAccess() {
    window.open(AppSettings.REQUEST_ACCESS_LINK, '_blank');
  }

  onClickLearnMore() {
    window.open(AppSettings.LEARN_MORE, '_blank');
  }
}
