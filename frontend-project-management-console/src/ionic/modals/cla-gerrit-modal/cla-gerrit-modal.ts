import { Component } from '@angular/core';
import { NavController, NavParams, ViewController, IonicPage, Events } from 'ionic-angular';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ExtraValidators } from '../../validators/requireSelfAnd';
import { ClaService } from '../../services/cla.service';
import { CincoService } from '../../services/cinco.service'
import { Http } from '@angular/http';

@IonicPage({
  segment: 'cla-gerrit-modal'
})
@Component({
  selector: 'cla-gerrit-modal',
  templateUrl: 'cla-gerrit-modal.html',
})
export class ClaGerritModal {
  form: FormGroup;
  submitAttempt: boolean = false;
  currentlySubmitting: boolean = false;
  projectId: string;
  user: any;

  constructor(
    public navParams: NavParams,
    public viewCtrl: ViewController,
    private cincoService: CincoService,
    private formBuilder: FormBuilder,
    public http: Http,
    public claService: ClaService,
    public events: Events
  ) {
    this.projectId = this.navParams.get('projectId'); 
    this.form = formBuilder.group({
      gerritName: ['', Validators.compose([Validators.required])],
      URL: ['', Validators.compose([Validators.required])],
      groupIdIcla: ['', (control)=>{ return ExtraValidators.requireSelfOr(control, "groupIdCcla")}],
      groupIdCcla: ['', (control)=>{ return ExtraValidators.requireSelfOr(control, "groupIdIcla")}],
    });

    events.subscribe('modal:close', () => {
      this.dismiss();
    });
  }

  ngOnInit() {
  }

  getDefaults() {
  }

  dismiss() {
    this.viewCtrl.dismiss();
  }


  submit() {
    this.submitAttempt = true;
    this.currentlySubmitting = true;
    if (!this.form.valid) {
      this.currentlySubmitting = false;
      // prevent submit
      return;
    }
    this.postGerritInstance();
  }

  postGerritInstance() {
    let gerrit = {
      project_id: this.projectId,
      gerrit_name: this.form.value.gerritName,
      gerrit_url: this.form.value.URL,
    };
    if (this.form.value.groupIdIcla && this.form.value.groupIdIcla.trim() != ''){
      gerrit["group_id_icla"] = this.form.value.groupIdIcla
    }
    if (this.form.value.groupIdCcla && this.form.value.groupIdCcla.trim() != ''){
      gerrit["group_id_ccla"] = this.form.value.groupIdCcla
    }
    this.claService.postGerritInstance(gerrit).subscribe((response) => {
      if(response.error_icla) {
        this.form.controls['groupIdIcla'].setErrors({groupNotExistentError: "The specified LDAP group for ICLA does not exist."})
      }
      else if(response.error_ccla) {
        this.form.controls['groupIdCcla'].setErrors({groupNotExistentError: "The specified LDAP group for CCLA does not exist."})
      }
      else {
        this.dismiss();
      }
    }, 
    error => {
      let errorObject = error.json()
        if(errorObject.errors) { 
          //TODO: Handle other types of backend errors. 
          this.form.controls['URL'].setErrors({invalidURL: "Invalid URL specified."})
        }
    },
    completion => {
    });
  }

  
}
