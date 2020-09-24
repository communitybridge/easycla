// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { environment } from 'src/environments/environment';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  title = 'EasyCLA';
  projectSection;
  organizationSection;
  developerSection;
  hasExapnded;

  constructor(
  ) {
    this.initializeData();
    this.mounted();
    this.hasExapnded = true;
  }

  initializeData() {
    this.projectSection = {
      title: 'Projects',
      subtitle: 'Reduce administrative hassles of supporting the CLA for your project.',
      highlights: [
        'Look in one place to find the companies and individuals who have signed the CLA',
        'Enable Companies to manage authorization of their own developers',
        'Support both Individual and Corporate Contributors within a single portal'
      ],
      footerTitle: 'Are you a Community Manager or Project Manager?'
    };

    this.organizationSection = {
      title: 'Organizations',
      subtitle: 'Enable all your developers to contribute code easily and quickly while remaining compliant with contribution policies.',
      highlights: [
        'Approved developers based on email, domain, GitHub handle, or GitHub organization',
        'Enable your signatories and contributors to sign CLAs using DocuSign® electronic signatures',
        'Enforce signing of the Corporate CLA by your developers without slowing them down with manual bureaucracy'
      ],
      footerTitle: 'Are you a CLA Manager for your organization?'
    };

    this.developerSection = {
      title: 'Developers',
      subtitle: 'Get started contributing code faster and with less friction.',
      highlights: [
        'Receive an automatic notification in GitHub or Gerrit if you need to be authorized',
        'Sign your Individual CLA with an e-signature',
        'Start contributing faster with a streamlined authorization workflow for Corporate CLAs'
      ],
      footerTitle: 'Do you contribute code to a project that uses CLAs?'
    };
  }


  onToggled() {
    this.hasExapnded = !this.hasExapnded;
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


