import { Component } from '@angular/core';
import { NavController, NavParams, ViewController, IonicPage, } from 'ionic-angular';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { FormControl } from '@angular/forms';
import { CincoService } from '../../services/cinco.service';
import { Http } from '@angular/http';

@IonicPage({
  segment: 'cla-organization-provider-modal'
})
@Component({
  selector: 'cla-organization-provider-modal',
  templateUrl: 'cla-organization-provider-modal.html',
  providers: [CincoService]
})
export class ClaOrganizationProviderModal {
  form: FormGroup;
  submitAttempt: boolean = false;
  currentlySubmitting: boolean = false;

  uploadInfo: any;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public viewCtrl: ViewController,
    private formBuilder: FormBuilder,
    public http: Http,
  ) {
    this.getDefaults();
    // this.uploadInfo = this.navParams.get('uploadInfo');
    this.form = formBuilder.group({
      provider: ['', Validators.required],
      organization: ['', Validators.compose([Validators.required])/*, this.urlCheck.bind(this)*/],
    });
  }

  // urlCheck(control: FormControl) {
  //   let url = 'https://api.github.com/users/' + control.value;
  //   return this.http.get(url).map(res => {
  //     let json = res.json();
  //     return json;
  //   });
  // }

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

    // After submitting for the organization
  }

  dismiss() {
    this.viewCtrl.dismiss();
  }

}
