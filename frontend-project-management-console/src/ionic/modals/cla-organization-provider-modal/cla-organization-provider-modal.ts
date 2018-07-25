import { Component } from '@angular/core';
import { NavController, NavParams, ViewController, IonicPage, } from 'ionic-angular';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { FormControl } from '@angular/forms';
import { ClaService } from 'cla-service';
import { Http } from '@angular/http';

@IonicPage({
  segment: 'cla-organization-provider-modal'
})
@Component({
  selector: 'cla-organization-provider-modal',
  templateUrl: 'cla-organization-provider-modal.html',
})
export class ClaOrganizationProviderModal {
  form: FormGroup;
  submitAttempt: boolean = false;
  currentlySubmitting: boolean = false;

  claProjectId: any;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public viewCtrl: ViewController,
    private formBuilder: FormBuilder,
    public http: Http,
    public claService: ClaService,
  ) {
    this.claProjectId = this.navParams.get('claProjectId');
    this.form = formBuilder.group({
      // provider: ['', Validators.required],
      orgName: ['', Validators.compose([Validators.required])/*, this.urlCheck.bind(this)*/],
    });
    // this.getDefaults();
  }

  // urlCheck(control: FormControl) {
  //   let url = 'https://api.github.com/users/' + control.value;
  //   return this.http.get(url).map(res => {
  //     let json = res.json();
  //     return json;
  //   });
  //
  //   CLA Service
  //   GET /github/check/namespace/<namespace>
  //    Will return true if namespace exists in GitHub, false otherwise
  //    Example: true
  //    Example: false
  //
  // GET /github/get/namespace/<namespace>
  //    Will return an object of data on the account requested
  //    Example: {"bio": "The Linux Foundation is a non-profit consortium dedicated to fostering the growth of Linux.", "company": null, "email": null, "created_at": <datetime object>, "location": null, "login": "linuxfoundation", "type": "Organization"}
  //    Example: {"errors": {"namespace": "Invalid GitHub account namespace"}}
  //
  //  getGithubOrganizations()
  //    will return if there is already a github org in the cla with this name
  //
  // }
  //
  getDefaults() {

  }

  ngOnInit() {

  }

  submit() {
    this.submitAttempt = true;
    this.currentlySubmitting = true;
    if (!this.form.valid) {
      this.currentlySubmitting = false;
      // prevent submit
      return;
    }
    this.postClaGithubOrganization();
  }

  postClaGithubOrganization() {
    let organization = {
      organization_project_id: this.claProjectId,
      organization_name: this.form.value.orgName,
    };
    this.claService.postGithubOrganization(organization).subscribe((response) => {
      this.dismiss()
    });
  }

  dismiss() {
    this.viewCtrl.dismiss();
  }

}
