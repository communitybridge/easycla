// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component, ViewChild, Input } from '@angular/core';
import { IonicPage, Nav, NavController, NavParams, Events } from 'ionic-angular';
import { ClaService } from '../../../services/cla.service';
import { Restricted } from '../../../decorators/restricted';
import { FormGroup, FormControl, Validators } from '@angular/forms';

@Restricted({
  roles: ['isAuthenticated', 'isPmcUser']
})

@IonicPage({
  segment: 'project/:projectId/cla/template/:projectTemplateId'
})

@Component({
  selector: 'project-cla-template',
  templateUrl: 'project-cla-template.html'
})

export class ProjectClaTemplatePage {
  projectId: string;
  templates: any[] = [];
  selectedTemplate: any;
  pdfPath = {
    corporatePDFURL: '',
    individualPDFURL: ''
  };
  currentPDF = 'corporatePDFURL';
  step = 'selection';
  buttonGenerateEnabled = true;
  message = null;
  loading = {
    documents: false
  };
  submitAttempt = false;

  @Input() form: FormGroup;
  @ViewChild(Nav) nav: Nav;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public claService: ClaService,
    public events: Events
  ) {
    this.projectId = navParams.get('projectId');
    this.getTemplates();
  }

  getTemplates() {
    this.claService.getTemplates().subscribe((templates) => (this.templates = templates));
  }

  ngOnInit() {
    this.setLoadingSpinner(false);
  }
  
  getPdfPath() {
    // PDF renderer library - or we can simply use the default browser behavior
    return this.pdfPath[this.currentPDF];
  }

  showPDF(type) {
    this.currentPDF = type;
  }

  selectTemplate(template) {
    this.selectedTemplate = template;
    this.generateFormGroup();
    this.step = 'values';
  }

  generateFormGroup() {
    let group = {};
    const metaFields = this.selectedTemplate.metaFields;
    metaFields.forEach(template => {
      let field = this.getFieldName(template);
      group[field] = new FormControl('', Validators.required);
    })
    this.form = new FormGroup(group);
  }

  getFieldName(template) {
    return template.name.split(' ').join('');
  }

  reviewSelectedTemplate() {
    this.submitAttempt = true;
    if (this.form.valid) {
      this.setLoadingSpinner(true);
      this.buttonGenerateEnabled = false;
      this.message = 'Generating PDFs...';

      const metaFields = this.selectedTemplate.metaFields;
      metaFields.forEach((metaField) => {
        metaField.value = this.form.value[this.getFieldName(metaField)];
      });
      let data = {
        templateID: this.selectedTemplate.ID,
        metaFields: metaFields
      };

      this.claService.postClaGroupTemplate(this.projectId, data).subscribe(
        (response) => {
          this.setLoadingSpinner(false);
          this.buttonGenerateEnabled = true;
          this.message = null;
          this.pdfPath = response;
          this.goToStep('review');
        },
        (error) => {
          this.setLoadingSpinner(false);
          this.buttonGenerateEnabled = true;
          this.message = 'Error creating PDFs: ' + error;
        }
      );
    }
  }

  goToStep(step) {
    this.step = step;
  }

  backToProject() {
    if (this.navCtrl.getViews().length > 1) {
      if (this.pdfPath.corporatePDFURL || this.pdfPath.individualPDFURL) {
        this.events.publish('reloadProjectCla');
      }
      this.navCtrl.pop();
      return;
    }
    this.navCtrl.push('AllProjectsPage');
  }

  setLoadingSpinner(value: boolean) {
    this.loading = {
      documents: value
    };
  }
}


