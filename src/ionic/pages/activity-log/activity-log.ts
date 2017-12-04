import { Component } from '@angular/core';
import { NavController, IonicPage } from 'ionic-angular';
import { CincoService } from '../../services/cinco.service'
import { KeycloakService } from '../../services/keycloak/keycloak.service';

@IonicPage({
  segment: 'activity-log'
})
@Component({
  selector: 'activity-log',
  templateUrl: 'activity-log.html'
})
export class ActivityLogPage {
  events: any;

  constructor(
    public navCtrl: NavController,
    private cincoService: CincoService,
    private keycloak: KeycloakService
  ) {
    this.getDefaults();
  }

  getDefaults() {
    this.events = [];
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

  ngOnInit(){
    this.getEvents();
  }

  getEvents() {
    this.cincoService.getEventsForProject('a0941000002wByYAAU').subscribe(response => {
      if (response) {
        console.log(response);
      }
    });
  }

}
