import { Component, ViewChild } from '@angular/core';
import { NavController, NavParams, IonicPage, ModalController, Content } from 'ionic-angular';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { UrlValidator } from  '../../validators/url';
import { CincoService } from '../../services/cinco.service'
import { KeycloakService } from '../../services/keycloak/keycloak.service';
import { ProjectModel } from '../../models/project-model';

@IonicPage({
  segment: 'project-details/:projectId'
})
@Component({
  selector: 'project-details',
  templateUrl: 'project-details.html',
})
export class ProjectDetailsPage {

  projectId: string;

  project = new ProjectModel();

  membershipsCount: number;

  editProject: any;
  loading: any;
  form: FormGroup;
  submitAttempt: boolean = false;
  currentlySubmitting: boolean = false;
  @ViewChild(Content) content: Content;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    private cincoService: CincoService,
    public modalCtrl: ModalController,
    private formBuilder: FormBuilder,
    private keycloak: KeycloakService
  ) {

    this.editProject = {};
    this.projectId = navParams.get('projectId');
    this.getDefaults();
    this.form = formBuilder.group({
      name:[this.project.name, Validators.compose([Validators.required])],
      startDate:[this.project.startDate],
      status:[this.project.status],
      category:[this.project.category],
      sector:[this.project.sector],
      url:[this.project.url, Validators.compose([UrlValidator.isValid])],
      addressThoroughfare:[this.project.address.address.thoroughfare],
      addressPostalCode:[this.project.address.address.postalCode],
      addressLocalityName:[this.project.address.address.localityName],
      addressAdministrativeArea:[this.project.address.address.administrativeArea],
      addressCountry:[this.project.address.address.country],
      description:[this.project.description],
    });
  }

  ionViewCanEnter() {
    if(!this.keycloak.authenticated())
    {
      this.navCtrl.setRoot('LoginPage');
      this.navCtrl.popToRoot();
    }
    return this.keycloak.authenticated();
  }

  ionViewWillEnter() {
    if(!this.keycloak.authenticated())
    {
      this.navCtrl.push('LoginPage');
    }
  }

  ngOnInit() {
    this.getProject(this.projectId);
  }

  getProject(projectId) {
    let getMembers = true;
    this.cincoService.getProject(projectId, getMembers).subscribe(response => {
      if (response) {
        this.project = response;
        this.getProjectLogos(projectId);
        this.loading.project = false;

        this.form.patchValue({
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

  getProjectLogos(projectId){
    this.cincoService.getProjectLogos(projectId).subscribe(response => {
      if(response) {
        if(response[0]) {
          this.project.logoRef = response[0].publicUrl;
        }
      }
    });
  }

  submitEditProject() {
    this.submitAttempt = true;
    this.currentlySubmitting = true;
    if (!this.form.valid) {
      this.content.scrollToTop();
      this.currentlySubmitting = false;
      // prevent submit
      return;
    }
    let address = {
      address: {
        thoroughfare: this.form.value.addressThoroughfare,
        postalCode: this.form.value.addressPostalCode,
        localityName: this.form.value.addressLocalityName,
        administrativeArea: this.form.value.addressAdministrativeArea,
        country: this.form.value.addressCountry,
      },
      type: "BILLING",
    }
    this.editProject = {
      name: this.form.value.name,
      description: this.form.value.description,
      url: this.form.value.url,
      sector: this.form.value.sector,
      address: address,
      status: this.form.value.status,
      category: this.form.value.category,
      startDate: this.form.value.startDate,
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


  openAssetManagementModal() {
    let modal = this.modalCtrl.create('AssetManagementModal', {
      projectId: this.projectId,
    });
    modal.present();
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
