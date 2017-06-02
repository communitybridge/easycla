import { Component } from '@angular/core';
import { FormBuilder, Validators } from '@angular/forms';

import { NavController, IonicPage } from 'ionic-angular';

import { CincoService } from '../../app/services/cinco.service'

import { ProjectModel } from '../../models/project-model';

@IonicPage({
  segment: 'add-project'
})
@Component({
  selector: 'add-project',
  templateUrl: 'add-project.html'
})
export class AddProjectPage {

  project = new ProjectModel();
  public projectForm: any;
  constructor(public navCtrl: NavController, private cincoService: CincoService, public _form: FormBuilder,) {
    this.projectForm = this._form.group({
      "name": ["", Validators.required],
      "startDate": ["", Validators.required],
      "status": ["", Validators.required],
      "category": ["", Validators.required],
      "sector": ["", Validators.required],
      "url": ["", Validators.required],
      "thoroughfare": ["", Validators.required],
      "postalCode": ["", Validators.required],
      "localityName": ["", Validators.required],
      "administrativeArea": ["", Validators.required],
      "country": ["", Validators.required],
      "description": ["", Validators.required]
    });
    this.getDefaults();
  }

  submitNewProject() {
    this.project.startDate = new Date(this.project.startDate).toISOString();
    this.cincoService.postProject(this.project).subscribe(response => {
      this.navCtrl.push('ProjectPage', {
        projectId: response
      });
    });
  }

  changeLogo() {
    // TODO: WIP
    alert("Change Logo");
  }

  getDefaults() {
    this.project = {
      id: "",
      name: "",
      description: "",
      managers: "",
      members: "",
      status: "",
      category: "",
      sector: "",
      url: "",
      startDate: "",
      logoRef: "",
      agreementRef: "",
      mailingListType: "",
      emailAliasType: "",
      address: {
        address: {
          administrativeArea: "",
          country: "",
          localityName: "",
          postalCode: "",
          thoroughfare: ""
        },
        type: "BILLING"
      }
    };
  }

}
