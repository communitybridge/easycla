import { Component } from '@angular/core';
import { DatePipe } from '@angular/common';
import { NavController, NavParams, ViewController, IonicPage, } from 'ionic-angular';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ClaService } from '../../services/cla.service';
import { Events } from 'ionic-angular';

@IonicPage({
  segment: 'cla-contract-upload-modal'
})
@Component({
  selector: 'cla-contract-upload-modal',
  templateUrl: 'cla-contract-upload-modal.html',
})
export class ClaContractUploadModal {
  loading: any;
  form: FormGroup;
  submitAttempt: boolean = false;
  currentlySubmitting: boolean = false;
  keysGetter: any;

  claProjectId: any;
  documentType: string; // individual | corporate
  templateOptions: any;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public viewCtrl: ViewController,
    private formBuilder: FormBuilder,
    public claService: ClaService,
    private datePipe: DatePipe,
    public events: Events
  ) {
    this.documentType = this.navParams.get('documentType');
    this.claProjectId = this.navParams.get('claProjectId');
    this.form = formBuilder.group({
      templateName: ['', Validators.compose([Validators.required])],
      legalEntityName:['', Validators.compose([Validators.required])],
      preamble:['', Validators.compose([Validators.required])],
      newSignature:[false],
    });
    this.getDefaults();
    this.keysGetter = Object.keys;

    events.subscribe('modal:close', () => {
      this.dismiss();
    });
  }

  getDefaults() {
    this.loading = {
      document: true,
    };
    this.templateOptions = {
      CNCFTemplate: 'CNCF Template',
      OpenBMCTemplate: 'OpenBMC Template',
      TungstenFabricTemplate: 'Tungsten Fabric Template',
      OpenColorIOTemplate: 'OpenColorIO Template',
      OpenVDBTemplate: 'OpenVDB Template'
    };
  }

  ngOnInit() {
    this.getProjectDocument();
  }

  getProjectDocument() {
    this.claService.getProjectDocument(this.claProjectId, this.documentType).subscribe((document) => {
      this.form.patchValue({
        legalEntityName: document.document_legal_entity_name,
        preamble: document.document_preamble,
      });
      this.loading.document = false;
    });
  }

  autofillFields() {
    let templateKey = this.form.value.templateName;
    let templateName = this.templateOptions[templateKey];
    this.form.patchValue({
      legalEntityName: templateName,
      preamble: templateName,
    });
  }

  generateDocumentName(text) {
    let simplifiedEntityName = text.replace(/[ ]/g, '-');
    let docType = 'cla';
    if (this.documentType == 'individual') {
      docType = 'icla'
    } else if (this.documentType == 'corporate') {
      docType = 'ccla';
    }
    let date = this.datePipe.transform(Date(), 'yyyy-MM-dd');
    return simplifiedEntityName + '_' + docType + '_' + date;
  }

  submit() {
    this.submitAttempt = true;
    this.currentlySubmitting = true;
    let template = this.form.value.templateName;
    if(template && template !== 'custom') {
      // autofill if using template
      this.autofillFields();
    }
    if (!this.form.valid) {
      this.currentlySubmitting = false;
      // prevent submit
      return;
    }
    this.postProjectDocumentTemplate();
  }

  postProjectDocumentTemplate() {
    let documentName = this.generateDocumentName(this.form.value.legalEntityName);

    let document = {
      document_name: documentName,
      document_preamble: this.form.value.preamble,
      document_legal_entity_name: this.form.value.legalEntityName,
      template_name: this.form.value.templateName,
      new_major_version: this.form.value.newSignature,
    };

     this.claService.postProjectDocumentTemplate(this.claProjectId, this.documentType, document).subscribe(
       (response) => {
         this.dismiss();
         this.currentlySubmitting = false;
       },
       (error) => {
         this.currentlySubmitting = false; // don't lock up form on failure
       }
     );
  }

  dismiss() {
    this.viewCtrl.dismiss();
  }

}
