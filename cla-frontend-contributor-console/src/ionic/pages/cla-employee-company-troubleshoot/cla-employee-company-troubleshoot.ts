// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { IonicPage, ModalController, NavParams } from 'ionic-angular';
import { ClaService } from '../../services/cla.service';
import { ClaSignatureModel } from '../../../../../cla-frontend-corporate-console/src/ionic/models/cla-signature';
import { generalConstants } from '../../constants/general';

@IonicPage({
  segment: 'cla/project/:projectId/user/:userId/employee/company/:companyId/troubleshoot'
})
@Component({
  selector: 'cla-employee-company-troubleshoot',
  templateUrl: 'cla-employee-company-troubleshoot.html'
})
export class ClaEmployeeCompanyTroubleshootPage {
  loading: any;
  projectId: string;
  repositoryId: string;
  userId: string;
  companyId: string;
  authenticated: boolean;
  cclaSignature: any;
  project: any;
  company: any;
  gitService: string;

  constructor(
    private modalCtrl: ModalController,
    public navParams: NavParams,
    private claService: ClaService
  ) {
    this.getDefaults();
    this.projectId = navParams.get('projectId');
    this.repositoryId = navParams.get('repositoryId');
    this.userId = navParams.get('userId');
    this.companyId = navParams.get('companyId');
    this.gitService = navParams.get('gitService');
    this.authenticated = navParams.get('authenticated');
  }

  getDefaults() {
    this.loading = {};
    this.project = {
      project_name: '',
      logoUrl: ''
    };
    this.company = {
      company_name: ''
    };
    this.cclaSignature = new ClaSignatureModel();
  }

  ngOnInit() {
    this.getProject(this.projectId);
    this.getCompany(this.companyId);
    this.getProjectSignatures(this.projectId, this.companyId);
  }

  getProject(projectId: string) {
    this.loading.projects = true;
    this.claService.getProject(projectId).subscribe((response) => {
      this.loading.projects = false;
      this.project = response;
    });
  }

  getCompany(companyId: string) {
    this.loading.companies = true;
    this.claService.getCompany(companyId).subscribe((response) => {
      this.loading.companies = true;
      this.company = response;
    });
  }

  getProjectSignatures(projectId: string, companyId: string) {
    // Get CCLA Company Signatures - should just be one
    this.loading.signatures = true;
    this.claService.getCompanyProjectSignatures(companyId, projectId).subscribe(
      (response) => {
        this.loading.signatures = false;
        if (response.signatures) {
          let cclaSignatures = response.signatures.filter((sig) => sig.signatureType === 'ccla');
          if (cclaSignatures.length) {
            this.cclaSignature = cclaSignatures[0];
            // Sort the values
            if (this.cclaSignature.githubOrgWhitelist) {
              const sortedList: string[] = this.cclaSignature.githubOrgWhitelist.sort((a, b) => {
                return a.trim().localeCompare(b.trim());
              });
              // Remove duplicates - set doesn't allow dups
              this.cclaSignature.githubOrgWhitelist = Array.from(new Set(sortedList));
            }
          }
        }
      },
      (exception) => {
        this.loading.signatures = false;
      }
    );
  }

  openGitServiceEmailSettings() {
    window.open(generalConstants.githubEmailURL, '_blank');
  }

  openClaEmployeeRequestAccessModal() {
    let modal = this.modalCtrl.create('ClaEmployeeRequestAccessModal', {
      projectId: this.projectId,
      repositoryId: this.repositoryId,
      userId: this.userId,
      companyId: this.companyId,
      authenticated: this.authenticated
    });
    modal.present();
  }
}
