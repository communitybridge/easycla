import { Component } from "@angular/core";
import { IonicPage, ModalController, NavController } from "ionic-angular";
import { ClaService } from "../../services/cla.service";
import { RolesService } from "../../services/roles.service";
import { Restricted } from "../../decorators/restricted";
import { ColumnMode, SelectionType, SortType } from "@swimlane/ngx-datatable";
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { EmailValidator } from "../../validators/email";


@IonicPage()
@Component({
  selector: 'page-cla-manager-onboarding',
  templateUrl: 'cla-manager-onboarding.html',
})

export class ClaManagerOnboardingPage {
  loading: any;
  companies: any;

  userId: string;
  userEmail: string;
  userName: string;

  manager: string;
  formErrors: any[]
  form: FormGroup;
  submitAttempt: boolean = false;
  currentlySubmitting: boolean = false;
  formSuccessfullySubmitted: boolean = false;
  claManagerApproved: boolean = false;


  constructor(
    public navCtrl: NavController,
    private claService: ClaService,
    private formBuilder: FormBuilder,
    private rolesService: RolesService // for @Restricted

  ) {
    this.form = formBuilder.group({
      project_name: [''],
      compnay_name: ['', Validators.compose([Validators.required])],
      full_name: ['', Validators.compose([Validators.required])],
      lfid: [''],
      email_address: ['', Validators.compose([Validators.required, EmailValidator.isValid])],
    });
    this.formErrors = [];
  }

  ngOnInit() {
    this.getDefaults();
  }

  getDefaults() {
    this.loading = {
      companies: true
    };
    this.userId = localStorage.getItem("userid");
    this.userEmail = localStorage.getItem("user_email");
    this.userName = localStorage.getItem("user_name");
    this.setUserDetails();
  }

  submit() {
    // Reset our status and error messages
    this.submitAttempt = true;
    this.currentlySubmitting = true;

    this.claService.createCLAManagerRequest(
      this.form.value.lfid,
      this.form.value.project_name,
      this.form.value.compnay_name,
      this.form.value.full_name,
      this.form.value.email_address
    ).subscribe((response) => {
      console.log(response, 'this is response')
    })
  }

  setUserDetails() {
    this.form.controls['lfid'].setValue(this.userId);
    this.form.controls['email_address'].setValue(this.userEmail)
    this.form.controls['full_name'].setValue(this.userName);
  }
}
