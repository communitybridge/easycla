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
  segment: 'project/:projectId/groups'
})
@Component({
  selector: 'project-groups',
  templateUrl: 'project-groups.html',
  providers: [CincoService]
})

export class ProjectGroupsPage {

  projectId: string;
  keysGetter;
  projectPrivacy;

  groupName: string;
  groupDescription: string;
  groupPrivacy;
  subgroupPermissions;

  form: FormGroup;
  submitAttempt: boolean = false;
  currentlySubmitting: boolean = false;

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

    this.form = formBuilder.group({
      groupName:[this.groupName, Validators.compose([Validators.required])],
      groupPrivacy:[this.groupPrivacy],
      groupDescription:[this.groupDescription],
      subgroupPermissions:[this.subgroupPermissions],
    });

  }

  ngOnInit() {
    this.getProjectConfig(this.projectId);
    this.getDefaults();
  }

  getDefaults() {
    this.keysGetter = Object.keys;
    this.getGroupPrivacy();
    this.getSubgroupPermissions();
  }

  getProjectConfig(projectId) {
    this.cincoService.getProjectConfig(projectId).subscribe(response => {
      if (response) {
        console.log(response)
      }
    });
  }

  getGroupPrivacy() {
    this.groupPrivacy = ['membersVisible'];
    // TODO Implement CINCO side
    // this.cincoService.getProjectGroupPrivacy().subscribe(response => {
    //   this.groupPrivacy = response;
    // });
  }

  getSubgroupPermissions() {
    this.subgroupPermissions = ['ownersOnly'];
    // TODO Implement CINCO side
  }

}
