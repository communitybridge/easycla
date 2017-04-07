import { Component } from '@angular/core';

import { NavController, ModalController, NavParams } from 'ionic-angular';

import { ContactUpdate } from '../contact-update/contact-update';

@Component({
  selector: 'member',
  templateUrl: 'member.html'
})
export class MemberPage {
  selectedProject: any;
  selectedMember: any;
  contacts: Array<{
    firstname: string,
    lastname: string,
    photo: string,
    primary?: boolean,
    role: string,
    title: string,
    timezone?: string,
    email: string,
    phone?: string,
    bio?: string,
  }>;

  constructor(
    public navCtrl: NavController,
    public modalCtrl: ModalController,
    public navParams: NavParams
  ) {
    // If we navigated to this page, we will have an item available as a nav param
    this.selectedProject = navParams.get('project');
    this.selectedMember = navParams.get('member');

    // use selectedProject and selectedMember to get data from CINCO and populate this.contacts

    this.contacts = [
      {
        firstname: "John",
        lastname: "Mathis",
        photo: "https://api.adorable.io/avatars/136/Johanna",
        role: "Board",
        title: "VP Design",
        timezone: "EST (+4)",
        email: "board@google.com",
        phone: "123-456-7890",
        bio: "Something something some bio.",
      },
      {
        firstname: "Carl",
        lastname: "Carlson",
        photo: "https://api.adorable.io/avatars/136/Carl",
        primary: true,
        role: "Marketing",
        title: "Director, Product Development",
        timezone: "EST (+4)",
        email: "ccarlson@google.com",
        phone: "123-456-7890",
        bio: "Something something some bio.",
      },
      {
        firstname: "Susan",
        lastname: "Star",
        photo: "https://api.adorable.io/avatars/136/Carl",
        primary: true,
        role: "Marketing",
        title: "Director, Open Technology",
        timezone: "EST (+4)",
        email: "ccarlson@google.com",
        phone: "123-456-7890",
      },
    ]

  }

  contactSelected(event, contact) {
    let modal = this.modalCtrl.create(ContactUpdate, {
      project: this.selectedProject,
      member: this.selectedMember,
      contact: contact,
    });
    modal.present();
  }

}
