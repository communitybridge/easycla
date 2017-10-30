import { Component } from '@angular/core';
import { NavController, ModalController, NavParams, IonicPage } from 'ionic-angular';
import { CincoService } from '../../../services/cinco.service';
import { KeycloakService } from '../../../services/keycloak/keycloak.service';
import { SortService } from '../../../services/sort.service';

@IonicPage({
  segment: 'project/:projectId/files'
})
@Component({
  selector: 'project-files',
  templateUrl: 'project-files.html',
  providers: [CincoService]
})
export class ProjectFilesPage {
  projectId: string;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    private cincoService: CincoService,
    private sortService: SortService,
    public modalCtrl: ModalController,
    private keycloak: KeycloakService
  ) {
    this.projectId = navParams.get('projectId');
    this.getDefaults();
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

  }

  getDefaults() {

  }


}
