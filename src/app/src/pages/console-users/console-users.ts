import { Component } from '@angular/core';

import { NavController, IonicPage, ModalController } from 'ionic-angular';

import { CincoService } from '../../app/services/cinco.service'

@IonicPage({
  segment: 'console-users'
})
@Component({
  selector: 'console-users',
  templateUrl: 'console-users.html'
})
export class ConsoleUsersPage {
  users: any;

  constructor(
    public navCtrl: NavController,
    private cincoService: CincoService,
    public modalCtrl: ModalController,
  ) {
    this.getDefaults();
  }

  getDefaults() {
    this.users = [];
  }

  ngOnInit(){
    this.getAllUsers();
  }

  getAllUsers() {
    this.cincoService.getAllUsers().subscribe(response => {
      if(response) {
        this.users = response;
      }
    });
  }

  userSelected(user) {
    let modal = this.modalCtrl.create('ConsoleUserUpdateModal', {
      user: user,
    });
    modal.present();
    // this.navCtrl.push('MemberPage', {
    //   projectId: member.projectId,
    //   memberId: member.id,
    // });
  }

}
