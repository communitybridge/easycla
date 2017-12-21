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

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    private cincoService: CincoService,
    private keycloak: KeycloakService,
    private domSanitizer : DomSanitizer,
    public modalCtrl: ModalController,
    public rolesService: RolesService,
  ) {
    this.projectId = navParams.get('projectId');
  }

  ngOnInit() {
    this.getProjectConfig(this.projectId);
    this.getDefaults();
  }

  getDefaults() {
    this.keysGetter = Object.keys;
    this.getProjectPrivacy();
  }

  getProjectConfig(projectId) {
    this.cincoService.getProjectConfig(projectId).subscribe(response => {
      if (response) {
        console.log(response)
      }
    });
  }

  getProjectPrivacy() {
    this.projectPrivacy = ['membersVisible'];
    // TODO Implement CINCO side
    // this.cincoService.getProjectPrivacy().subscribe(response => {
    //   this.projectPrivacy = response;
    // });
  }

}
