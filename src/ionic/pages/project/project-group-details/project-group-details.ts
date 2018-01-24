import { Component } from '@angular/core';
import { NavController, ModalController, NavParams, IonicPage } from 'ionic-angular';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { CincoService } from '../../../services/cinco.service';
import { KeycloakService } from '../../../services/keycloak/keycloak.service';
import { DomSanitizer} from '@angular/platform-browser';
import { RolesService } from '../../../services/roles.service';
import { Restricted } from '../../../decorators/restricted';

@Restricted({
  roles: ['isAuthenticated', 'isPmcUser'],
})
@IonicPage({
  segment: 'project/:projectId/group-details/:groupName'
})
@Component({
  selector: 'project-group-details',
  templateUrl: 'project-group-details.html',
  providers: [CincoService]
})

export class ProjectGroupDetailsPage {

  projectId: string;
  keysGetter;
  projectPrivacy;

  groupName: string;
  mailingGroup: any;
  groupDescription: string;
  groupPrivacy = [];
  subgroupPermissions = [];

  form: FormGroup;
  submitAttempt: boolean = false;
  currentlySubmitting: boolean = false;

  group: any;
  projectGroups: any;
  groupParticipants: any;

  participantName: any;
  participantEmail: any;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    private cincoService: CincoService,
    private keycloak: KeycloakService,
    private domSanitizer : DomSanitizer,
    public modalCtrl: ModalController,
    public rolesService: RolesService,
    private formBuilder: FormBuilder,
  ) {
    this.projectId = navParams.get('projectId');
    this.groupName = navParams.get('groupName');

    this.form = formBuilder.group({
      participantName:[this.participantName, Validators.compose([Validators.required])],
      participantEmail:[this.participantEmail],
    });

  }

  ngOnInit() {
    this.getProjectConfig(this.projectId);
    this.getDefaults();
    this.getAllGroupParticipants();
  }

  getDefaults() {
    this.keysGetter = Object.keys;
    //TODO: Get participants via CINCO
    this.groupParticipants = [
      {
        address: ""
      }
    ]

  }

  getProjectConfig(projectId) {
    this.cincoService.getProjectConfig(projectId).subscribe(response => {
      if (response) {
        console.log(response);
        if (response.mailingGroup) this.mailingGroup = response.mailingGroup;
        else console.log("no domain");
      }
    });
  }

  getProjectGroup() {
    this.cincoService.getProjectGroup(this.projectId, this.groupName).subscribe(response => {
      console.log(response)
    });
  }

  getAllGroupParticipants(){
    this.cincoService.getAllGroupParticipants(this.projectId, this.groupName).subscribe(response => {
      console.log("getAllGroupParticipants");
      console.log(response);
      this.groupParticipants = response;
    });
  }

  submitParticipant() {
    this.submitAttempt = true;
    this.currentlySubmitting = true;
    if (!this.form.valid) {
      this.currentlySubmitting = false;
      // prevent submit
      return;
    }
    let participant = [{
        address: this.form.value.participantEmail,
        name: this.form.value.participantName
    }];
    console.log(participant);
    this.cincoService.addGroupParticipant(this.projectId, this.groupName, participant).subscribe(response => {
      console.log("addGroupParticipant")
      console.log(response)
      this.currentlySubmitting = false;
      this.cincoService.getAllGroupParticipants(this.projectId, this.groupName).subscribe(response => {
        console.log("getAllGroupParticipants");
        console.log(response);
        this.groupParticipants = response;
      });
    });
  }


}
