import { Component } from '@angular/core';
import { NavController, IonicPage } from 'ionic-angular';
import { KeycloakService } from '../../services/keycloak/keycloak.service';

@IonicPage({
  name: 'LoginPage',
  segment: 'login'
})
@Component({
  selector: 'login',
  templateUrl: 'login.html'
})
export class LoginPage {

  constructor(public navCtrl: NavController, private keycloak: KeycloakService) {

  }

  ionViewWillEnter() {
    console.log('login will enter');
    if(this.keycloak.authenticated())
    {
      this.navCtrl.setRoot('AllProjectsPage');
      this.navCtrl.popToRoot();
    }
  }

  ionViewCanLeave() {
    console.log('login can leave');
    return (this.keycloak.authenticated());
  }

  login() {
    console.log('login function called');
    if (this.keycloak.authenticated()) {
      this.navCtrl.setRoot('AllProjectsPage');
    }
    else{
      this.keycloak.login();
    }
  }

}
