// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component, OnInit } from '@angular/core';
import { environment } from 'src/environments/environment';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {
  hasExpanded: boolean;
  links: any[];

  constructor() {
    this.hasExpanded = true;
    this.links = [
      {
        title: 'Project Login',
        url: environment.PROJECT_LOGIN_URL
      },
      {
        title: 'CLA Manager Login',
        url: environment.CORPORATE_LOGIN_URL
      },
      {
        title: 'Developer',
        url: environment.CONTRIBUTOR_LOGIN_URL
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



