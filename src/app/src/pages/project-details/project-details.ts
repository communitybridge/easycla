import { Component, ViewChild } from '@angular/core';
import { NavController, NavParams, IonicPage, Content } from 'ionic-angular';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { CincoService } from '../../services/cinco.service'
import { ProjectModel } from '../../models/project-model';

@IonicPage({
  segment: 'project-details/:projectId'
})
@Component({
  selector: 'project-details',
  templateUrl: 'project-details.html'
})
export class ProjectDetailsPage {

  projectId: string;

  project = new ProjectModel();

  membershipsCount: number;

  editProject: any;
  loading: any;

  _form: FormGroup;
  submitAttempt: boolean = false;
  currentlySubmitting: boolean = false;
  @ViewChild(Content) content: Content;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    private cincoService: CincoService,
    private formBuilder: FormBuilder,
  ) {
    this.editProject = {};
    this.projectId = navParams.get('projectId');
    this.getDefaults();
    this._form = formBuilder.group({
      name:[this.project.name],
      startDate:[this.project.startDate],
      status:[this.project.status],
      category:[this.project.category],
      sector:[this.project.sector],
      url:[this.project.url],
      addressThoroughfare:[this.project.address.address.thoroughfare],
      addressPostalCode:[this.project.address.address.postalCode],
      addressLocalityName:[this.project.address.address.localityName],
      addressAdministrativeArea:[this.project.address.address.administrativeArea],
      addressCountry:[this.project.address.address.country],
      description:[this.project.description],
    });
  }

  ngOnInit() {
    this.getProject(this.projectId);
  }

  getProject(projectId) {
    let getMembers = true;
    this.cincoService.getProject(projectId, getMembers).subscribe(response => {
      if (response) {
        // this.project.id = response.id;
        // this.project.name = response.name;
        // this.project.description = response.description;
        // this.project.managers = response.managers;
        // this.project.members = response.members
        // this.project.status = response.status;
        // this.project.category = response.category;
        // this.project.sector = response.sector;
        // this.project.url = response.url;
        // this.project.startDate = response.startDate;
        // this.project.logoRef = response.logoRef;
        // this.project.agreementRef = response.agreementRef;
        // this.project.mailingListType = response.mailingListType;
        // this.project.emailAliasType = response.emailAliasType;
        // this.project.address = response.address;
        this.project = response;
        this.loading.project = false;

        this._form.patchValue({
          name:this.project.name,
          startDate:this.project.startDate,
          status:this.project.status,
          category:this.project.category,
          sector:this.project.sector,
          url:this.project.url,
          addressThoroughfare:this.project.address.address.thoroughfare,
          addressPostalCode:this.project.address.address.postalCode,
          addressLocalityName:this.project.address.address.localityName,
          addressAdministrativeArea:this.project.address.address.administrativeArea,
          addressCountry:this.project.address.address.country,
          description:this.project.description,
        });
      }
    });
  }

  submitEditProject() {
    this.submitAttempt = true;
    this.currentlySubmitting = true;
    if (!this._form.valid) {
      this.content.scrollToTop();
      this.currentlySubmitting = false;
      // prevent submit
      return;
    }
    this.editProject = {
      project_name: this.project.name,
      project_description: this.project.description,
      project_url: this.project.url,
      project_sector: this.project.sector,
      project_address: this.project.address,
      project_status: this.project.status,
      project_category: this.project.category,
      project_start_date: this.project.startDate
    };
    this.cincoService.editProject(this.projectId, this.editProject).subscribe(response => {
      this.currentlySubmitting = false;
      this.navCtrl.setRoot('ProjectPage', {
        projectId: this.projectId
      });
    });
  }

  cancelEditProject() {
    this.navCtrl.setRoot('ProjectPage', {
      projectId: this.projectId
    });
  }

  changeLogo() {
    // TODO: WIP
    alert("Change Logo");
  }

  getDefaults() {
    this.loading = {
      project: true,
    };
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
        type: ""
      }
    };
  }

}
