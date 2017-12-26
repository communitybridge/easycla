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
  groupPrivacy = [];
  subgroupPermissions = [];

  form: FormGroup;
  submitAttempt: boolean = false;
  currentlySubmitting: boolean = false;

  group: any;

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
    this.groupPrivacy = [];
    this.cincoService.getGroupPrivacy(this.projectId).subscribe(response => {
      this.groupPrivacy = response;
    });
  }

  getSubgroupPermissions() {
    this.subgroupPermissions = ['To_implement_CINCO_side'];
    // TODO Implement CINCO side
    // this.cincoService.getSubgroupPermissions(this.projectId).subscribe(response => {
    //   this.subgroupPermissions = response;
    // });
  }

  submitGroup() {
    this.submitAttempt = true;
    this.currentlySubmitting = true;
    if (!this.form.valid) {
      this.currentlySubmitting = false;
      // prevent submit
      return;
    }
    this.group = {
      projectId: this.projectId,
      name: this.form.value.groupName,
      description: this.form.value.groupDescription,
      privacy: this.groupPrivacy[this.form.value.groupPrivacy],
    };
    console.log("this.group");
    console.log(this.group);
    this.cincoService.createProjectGroup(this.projectId, this.group).subscribe(response => {
      this.currentlySubmitting = false;
      this.navCtrl.setRoot('ProjectPage', {
        projectId: this.projectId
      });
    });
  }

}
