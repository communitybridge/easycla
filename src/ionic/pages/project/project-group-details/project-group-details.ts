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
  segment: 'project/:projectId/group-details'
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
  groupDescription: string;
  groupPrivacy = [];
  subgroupPermissions = [];

  form: FormGroup;
  submitAttempt: boolean = false;
  currentlySubmitting: boolean = false;

  group: any;
  projectGroups: any;
  groupParticipants: any;

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
      groupName:[this.groupName, Validators.compose([Validators.required])],
      groupDescription:[this.groupDescription],
      groupPrivacy:[this.groupPrivacy],
      subgroupPermissions:[this.subgroupPermissions],
    });

  }

  ngOnInit() {
    this.getProjectConfig(this.projectId);
    this.getDefaults();
  }

  getDefaults() {
    this.keysGetter = Object.keys;
    //TODO: Get participants via CINCO
    this.groupParticipants = [
      {
        address: "todo@todo.com"
      }
    ]

  }

  getProjectConfig(projectId) {
    this.cincoService.getProjectConfig(projectId).subscribe(response => {
      if (response) {
        console.log(response);
        if (!response.domain) console.log("no domain");
      }
    });
  }

  getProjectGroup() {
    this.cincoService.getProjectGroup(this.projectId, this.groupName).subscribe(response => {
      console.log(response)
    });
  }

}
