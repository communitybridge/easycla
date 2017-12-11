import { Component } from '@angular/core';
import { DatePipe } from '@angular/common';
import { NavController, NavParams, ViewController, IonicPage, } from 'ionic-angular';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ClaService } from 'cla-service';

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

  claProjectId: any;
  documentType: string; // individual | corporate

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public viewCtrl: ViewController,
    private formBuilder: FormBuilder,
    public claService: ClaService,
    private datePipe: DatePipe,
  ) {
    this.documentType = this.navParams.get('documentType');
    this.claProjectId = this.navParams.get('claProjectId');
    this.form = formBuilder.group({
      legalEntityName:['', Validators.compose([Validators.required])],
      preamble:['', Validators.compose([Validators.required])],
      newSignature:[false],
    });
    this.getDefaults();
  }

  getDefaults() {
    this.loading = {
      document: true,
    };
  }

  ngOnInit() {
    this.getProjectDocument();
  }

  getProjectDocument() {
    this.claService.getProjectDocument(this.claProjectId, this.documentType).subscribe((document) => {
      console.log("document");
      console.log(document);
      this.form.patchValue({
        legalEntityName: document.document_legal_entity_name,
        preamble: document.document_preamble,
      });
      this.loading.document = false;
    });
  }

  submit() {
    this.submitAttempt = true;
    this.currentlySubmitting = true;
    if (!this.form.valid) {
      this.currentlySubmitting = false;
      // prevent submit
      return;
    }
    this.postProjectDocumentTemplate();
  }

  postProjectDocumentTemplate() {
    let simplifiedEntityName = this.form.value.legalEntityName.replace(/[ ]/g, '-');
    let docType = 'cla';
    if (this.documentType == 'individual') {
      docType = 'icla'
    } else if (this.documentType == 'corporate') {
      docType = 'ccla';
    }
    let date = this.datePipe.transform(Date(), 'yyyy-MM-dd');
    let documentName = simplifiedEntityName + '_' + docType + '_' + date;

    let document = {
      document_name: documentName,
      document_preamble: this.form.value.preamble,
      document_legal_entity_name: this.form.value.legalEntityName,
      template_name: 'CNCFTemplate', // only template supported for now
      new_major_version: this.form.value.newSignature,
    };
    console.log(this.documentType);
    console.log(document);

     this.claService.postProjectDocumentTemplate(this.claProjectId, this.documentType, document).subscribe((response) => {
       this.dismiss();
     });
  }

  dismiss() {
    this.viewCtrl.dismiss();
  }

}
