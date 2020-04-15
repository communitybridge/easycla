// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { IonicPage, NavController, NavParams, ViewController } from 'ionic-angular';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ClaService } from '../../services/cla.service';
import { ClaCompanyModel } from '../../models/cla-company';
import { PlatformLocation } from '@angular/common';

@IonicPage({
  segment: 'projects-ccla-select-modal'
})
@Component({
  selector: 'projects-ccla-select-modal',
  templateUrl: 'projects-ccla-select-modal.html'
})
export class ProjectsCclaSelectModal {
  projectId: any;
  form: FormGroup;
  projects: any;
  projectsFiltered: any;
  loading: any = true;
  submitDisabled: boolean = true;
  company: ClaCompanyModel;

  constructor(
    public navParams: NavParams,
    public navCtrl: NavController,
    public viewCtrl: ViewController,
    public formBuilder: FormBuilder,
    private claService: ClaService,
    private location: PlatformLocation
  ) {
    this.form = formBuilder.group({
      search: ['', Validators.compose([Validators.required])]
    });
    this.location.onPopState(() => {
      this.viewCtrl.dismiss(false);
    });
  }

  ngOnInit() {
    this.company = this.navParams.get('company');
    this.getProjectsCcla();
  }

  getProjectsCcla() {
    const companyId = this.navParams.get('companyId');
    this.claService.getCompanyUnsignedProjects(companyId).subscribe((response) => {
      this.loading = false;
      this.projects = response
        .filter((a) => a != null && a.project_name != null && a.project_name.trim().length > 0)
        .sort((a, b) => {
          return ('' + a.project_name.trim()).localeCompare(b.project_name.trim());
        });
      this.form.value.search = '';
      this.projectsFiltered = this.projects;
    });
  }

  onSearch() {
    const searchTerm = this.form.value.search;
    if (searchTerm === '') {
      this.projectsFiltered = this.projects;
      this.submitDisabled = true;
    } else {
      if (this.projectsFiltered !== undefined) {
        this.projectsFiltered = this.projects.filter((a) => {
          return a.project_name.toLowerCase().includes(searchTerm.toLowerCase());
        });
      }
    }
  }

  selectProject(project) {
    this.submitDisabled = false
    this.form.controls['search'].setValue(project.project_name);
    this.projectId = project.project_id
  }

  submit() {
    if (!this.submitDisabled) {
      this.navCtrl.push('AuthorityYesnoPage', {
        projectId: this.projectId,
        company: this.company
      });
      this.dismiss();
    }
  }

  dismiss() {
    this.viewCtrl.dismiss();
  }
}
