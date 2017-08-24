import { Component } from '@angular/core';
import { NavController, NavParams, ViewController, IonicPage, } from 'ionic-angular';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { CincoService } from '../../services/cinco.service'

@IonicPage({
  segment: 'cla-contract-config-modal'
})
@Component({
  selector: 'cla-contract-config-modal',
  templateUrl: 'cla-contract-config-modal.html',
  providers: [CincoService]
})
export class ClaContractConfigModal {
  form: FormGroup;
  submitAttempt: boolean = false;
  currentlySubmitting: boolean = false;

  contract: any;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public viewCtrl: ViewController,
    private formBuilder: FormBuilder,
  ) {
    this.getDefaults();
    this.contract = this.navParams.get('contract');
    this.form = formBuilder.group({
      name:[this.contract.name, Validators.compose([Validators.required])],
      ccla:[this.contract.ccla, Validators.compose([Validators.required])],
      cclaAndIcla:[this.contract.cclaAndIcla, Validators.compose([Validators.required])],
      icla:[this.contract.icla, Validators.compose([Validators.required])],
    });
  }

  ngOnInit() {

  }

  getDefaults() {

  }

  submit() {
    this.submitAttempt = true;
    this.currentlySubmitting = true;
    if (!this.form.valid) {
      this.currentlySubmitting = false;
      // prevent submit
      return;
    }
    // let data = this.form.value.field;
    // do any pre-processing of data
    // this.dataService.sendData(data).subscribe(response => {
    //   this.currentlySubmitting = false;
    //   // call any success messaging
    //   // navigate to previous page, root, or destination
    // });
  }

  dismiss() {
    this.viewCtrl.dismiss();
  }

}
