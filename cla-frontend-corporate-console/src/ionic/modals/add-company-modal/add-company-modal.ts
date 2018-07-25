import { Component, ChangeDetectorRef } from '@angular/core';
import { NavController, NavParams, ModalController, ViewController, AlertController, IonicPage } from 'ionic-angular';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ClaService } from 'cla-service';
import { ClaCompanyModel } from '../../models/cla-company';

@IonicPage({
  segment: 'add-company-modal'
})
@Component({
  selector: 'add-company-modal',
  templateUrl: 'add-company-modal.html',
})
export class AddCompanyModal {

  form: FormGroup;
  submitAttempt: boolean = false;
  currentlySubmitting: boolean = false;

  mode: string;
  company: ClaCompanyModel;

  constructor(
    public navParams: NavParams,
    public viewCtrl: ViewController,
    public formBuilder: FormBuilder,
    private claService: ClaService,
  ) {
    this.getDefaults();
  }

  getDefaults() {
    this.company = this.navParams.get('company');
    this.mode = this.company
      ? 'edit'
      : 'add';
    let companyName = this.company
      ? this.company.company_name
      : '';
    this.form = this.formBuilder.group({
      name:[companyName, Validators.compose([Validators.required])],
    });
  }

  ngOnInit() {

  }

  submit() {
    this.submitAttempt = true;
    this.currentlySubmitting = true;
    if (!this.form.valid){
      this.currentlySubmitting = false;
      // prevent submit
      return;
    }
    if (this.mode === 'add') {
      this.addCompany();
    } else {
      this.updateCompany();
    }
  }

  addCompany() {
    let company = {
      company_name: this.form.value.name,
      company_whitelist: [],
      company_whitelist_patterns: [],
    };
    this.claService.postCompany(company).subscribe(
      (response) => {
        this.currentlySubmitting = false;
        this.dismiss();
      },
      (error) => {
        this.currentlySubmitting = false;
      }
    );
  }

  updateCompany() {
    let company = {
      company_id: this.company.company_id,
      company_name: this.form.value.name,
      company_whitelist: this.company.company_whitelist,
      company_whitelist_patterns: this.company.company_whitelist_patterns,
    };
    this.claService.putCompany(company).subscribe(
      (response) => {
        this.currentlySubmitting = false;
        this.dismiss();
      },
      (error) => {
        this.currentlySubmitting = false;
      }
    );
  }

  dismiss() {
    this.viewCtrl.dismiss();
  }

}
