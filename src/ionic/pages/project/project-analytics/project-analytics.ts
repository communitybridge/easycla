import { Component } from '@angular/core';
import { NavController, NavParams, IonicPage } from 'ionic-angular';
import { CincoService } from '../../../services/cinco.service';
import { KeycloakService } from '../../../services/keycloak/keycloak.service';
import { DomSanitizer} from '@angular/platform-browser';

@IonicPage({
  segment: 'project/:projectId/analytics'
})
@Component({
  selector: 'project-analytics',
  templateUrl: 'project-analytics.html',
  providers: [CincoService]
})
export class ProjectAnalyticsPage {

  projectId: string;

  hasBitergia: boolean;
  analyticsUrl: any;
  sanitizedBitergiaUrl: any;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    private cincoService: CincoService,
    private keycloak: KeycloakService,
    private domSanitizer : DomSanitizer
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
    this.getProjectConfig(this.projectId);
  }

  getDefaults() {
    this.hasBitergia = true;
    this.sanitizedBitergiaUrl = this.domSanitizer.bypassSecurityTrustResourceUrl(this.analyticsUrl);
  }

  addAnalyticsUrl() {
    this.cincoService.getProjectConfig(this.projectId).subscribe(response => {
      if (response) {
        let updatedConfig = response;
        updatedConfig.analyticsUrl = this.analyticsUrl;
        this.cincoService.editProjectConfig(this.projectId, updatedConfig).subscribe(response => {
          if (response) {
            console.log(response);
          }
        });
      }
    });
  }

  getProjectConfig(projectId) {
    this.cincoService.getProjectConfig(projectId).subscribe(response => {
      if (response) {
        let projectConfig = response;
        if(projectConfig.analyticsUrl) {
          this.analyticsUrl = projectConfig.analyticsUrl;
          this.hasBitergia = true;
        }
        else{
          this.hasBitergia = false;
        }
      }
    });
  }

  openAnaylticsConfigModal(projectId) {

  }

}
