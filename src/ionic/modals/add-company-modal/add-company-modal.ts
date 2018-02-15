import { Component, ChangeDetectorRef } from '@angular/core';
import { NavController, NavParams, ModalController, ViewController, AlertController, IonicPage } from 'ionic-angular';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ClaService } from 'cla-service';

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

  constructor(
    public viewCtrl: ViewController,
    public formBuilder: FormBuilder,
    private claService: ClaService,
  ) {
    this.getDefaults();
  }

  getDefaults() {
    this.form = this.formBuilder.group({
      name:['', Validators.compose([Validators.required])],
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
    let company = {
      company_name: this.form.value.name,
      company_whitelist: [],
      company_whitelist_patterns: [],
    };
    this.claService.postCompany(company).subscribe(
      (response) => {
        console.log('response');
        console.log(response);
        this.currentlySubmitting = false;
        this.dismiss();
      },
      (error) => {
        console.log('error');
        console.log(error);
        this.currentlySubmitting = false;
      }
    );

  }

  dismiss() {
    this.viewCtrl.dismiss();
  }

}
