// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component, OnInit } from '@angular/core';
import { environment } from 'src/environments/environment';
import { AppSettings } from './config/app-settings';
import { EnvConfig } from './config/cla-env-utils';
import { LfxHeaderService } from './core/services/lfx-header.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {
  hasExpanded: boolean;
  links: any[];

  constructor(
    private lfxHeaderService: LfxHeaderService
  ) {
    this.hasExpanded = true;
    this.links = [
      {
        title: 'Project Login',
        url: EnvConfig.default[AppSettings.PROJECT_CONSOLE_LINK] + '#/login'
      },
      {
        title: 'CLA Manager Login',
        url: EnvConfig.default[AppSettings.CORPORATE_CONSOLE_LINK] + '#/login'
      },
      {
        title: 'Developer',
        url: AppSettings.LEARN_MORE
      }
    ];
    this.mounted();
  }

  onToggled() {
    this.hasExpanded = !this.hasExpanded;
  }

  ngOnInit() {
    const element: any = document.getElementById('lfx-header');
    element.links = this.links;
  }

  mounted() {
    const script = document.createElement('script');
    script.setAttribute(
      'src',
      environment.LFX_HEADER_URL
    );
    document.head.appendChild(script);
  }
}



