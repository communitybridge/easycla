// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component, ViewChild, Input } from '@angular/core';
import { IonicPage, Nav, NavController, NavParams } from 'ionic-angular';
import { ClaService } from '../../../services/cla.service';
import { Restricted } from '../../../decorators/restricted';
import { DomSanitizer } from '@angular/platform-browser';
import { RolesService } from '../../../services/roles.service';
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
  sfdcProjectId: string;
  projectId: string;
  templates: any[] = [];
  selectedTemplate: any;
  templateValues = {};
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
  @Input() form: FormGroup;
  @ViewChild(Nav) nav: Nav;
  submitAttempt = false;
  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public claService: ClaService,
    public sanitizer: DomSanitizer,
    public rolesService: RolesService
  ) {
    this.sfdcProjectId = navParams.get('sfdcProjectId');
    this.projectId = navParams.get('projectId');
    this.getDefaults();
  }

  getDefaults() {
    this.getTemplates();
  }

  getTemplates() {
    this.claService.getTemplates().subscribe((templates) => (this.templates = templates));
  }

  ngOnInit() {
    this.setLoadingSpinner(false);
  }

  /**
   * Get the PDF path based on which CLA the user has selected on the UI.
   */
  getPdfPath() {
    // https://stackoverflow.com/questions/291813/recommended-way-to-embed-pdf-in-html#291823
    //return this.sanitizer.bypassSecurityTrustResourceUrl('https://drive.google.com/viewerng/viewer?embedded=true&url=' + this.pdfPath[this.currentPDF]);
    // Note: Google drive may not be accessible in certain countries - use default for now until we incorporate a
    // PDF renderer library - or we can simply use the default browser behavior
    console.log('generated', this.pdfPath[this.currentPDF]);
    return this.pdfPath[this.currentPDF];
    // return this.sanitizer.bypassSecurityTrustResourceUrl(this.pdfPath[this.currentPDF]);
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

  validateEmail(template: string) {
    if (template.toLowerCase().indexOf('email') >= 0) {
      if (/^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,3})+$/.test(template)) {
        return true
      } else {
        return false;
      }
    }
    return null;
  }

  reviewSelectedTemplate() {
    this.submitAttempt = true;
    if (this.form.valid) {
      this.setLoadingSpinner(true);
      this.buttonGenerateEnabled = false;
      this.message = 'Generating PDFs...';

      const metaFields = this.selectedTemplate.metaFields;
      metaFields.forEach((metaField) => {
        if (this.templateValues.hasOwnProperty(metaField.templateVariable)) {
          metaField.value = this.templateValues[metaField.templateVariable];
        }
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
    this.navCtrl.pop();
  }

  setLoadingSpinner(value: boolean) {
    this.loading = {
      documents: value
    };
  }
}


