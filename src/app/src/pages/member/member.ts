import { Component } from '@angular/core';
import { NavController, ModalController, NavParams, IonicPage } from 'ionic-angular';
import { CincoService } from '../../app/services/cinco.service'
import { ContactUpdate } from '../contact-update/contact-update';


@IonicPage({
  segment: 'project-page/:projectId/member-page/:memberId'
})
@Component({
  selector: 'page-member',
  templateUrl: 'member.html'
})
export class MemberPage {
  projectId: any;
  memberId: any;
  member: any;
  // contacts: Array<{
  //   firstname: string,
  //   lastname: string,
  //   profile_photo: string,
  //   primary?: boolean,
  //   role: string,
  //   title: string,
  //   timezone?: string,
  //   email: string,
  //   email_groups?: Array<{
  //     name: string,
  //   }>,
  //   phone?: string,
  //   bio?: string,
  //   photos?: Array<{}>,
  // }>;

  constructor(
    public navCtrl: NavController,
    public modalCtrl: ModalController,
    public navParams: NavParams,
    private cincoService: CincoService
  ) {
    // If we navigated to this page, we will have an item available as a nav param
    this.projectId = navParams.get('projectId');
    this.memberId = navParams.get('memberId');

    this.getDefaults();
  }

  ngOnInit() {
    // use selectedProject and selectedMember to get data from CINCO and populate this.contacts
    this.getMember(this.projectId, this.memberId);
  };

  getMember(projectId, memberId) {
    this.cincoService.getMember(projectId, memberId).subscribe(response => {
      if(response) {
        this.member = response;
      }
    });
  }

  contactSelected(event, contact) {
    let modal = this.modalCtrl.create(ContactUpdate, {
      projectId: this.projectId,
      member: this.member,
      contact: contact,
    });
    modal.present();
  }

  getDefaults() {
    this.member = {
      id: this.memberId,
      contacts: [
        {
          bio: "",
          email: "",
          familyName: "",
          givenName: "",
          headshotRef: "",
          id: "",
          phone: "",
          type: "",
        }
      ],
      invoices: [],
      org: {
        addresses: [],
        id: "",
        logoRef: "",
        name: "",
        phone: "",
        url: "",
      },
      projectId: "",
      renewalDate: "",
      startDate: "",
      tier: {
        qualifier: "",
        type: "",
      }
    };
    // console.log(this.member);
    // this.contacts = [
    //   {
    //     firstname: "John",
    //     lastname: "Mathis",
    //     profile_photo: "https://api.adorable.io/avatars/40/Johanna",
    //     role: "Board",
    //     title: "VP Design",
    //     timezone: "EST (+4)",
    //     email: "board@google.com",
    //     phone: "123-456-7890",
    //     bio: "Something something some bio.",
    //     photos: [
    //       {
    //         name: "alreadysaved.png",
    //       },
    //       {
    //         name: "existing.tif",
    //       }
    //     ],
    //   },
    //   {
    //     firstname: "Carl",
    //     lastname: "Carlson",
    //     profile_photo: "https://api.adorable.io/avatars/40/Carl",
    //     primary: true,
    //     role: "Marketing",
    //     title: "Director, Product Development",
    //     timezone: "EST (+4)",
    //     email: "ccarlson@google.com",
    //     email_groups: [
    //       {
    //         name: "marketing@zephyr.com",
    //       },
    //       {
    //         name: "info@zephyr.com",
    //       },
    //       {
    //         name: "leadership@zephyr.com",
    //       },
    //     ],
    //     phone: "123-456-7890",
    //   },
    //   {
    //     firstname: "Susan",
    //     lastname: "Star",
    //     profile_photo: "https://api.adorable.io/avatars/40/Susan",
    //     role: "Technical",
    //     title: "Director, Open Technology",
    //     email: "email@google.com",
    //     phone: "123-456-7890",
    //   },
    //   {
    //     firstname: "Name",
    //     lastname: "Name",
    //     profile_photo: "https://api.adorable.io/avatars/40/Name",
    //     role: "Marketing 2",
    //     title: "Title",
    //     email: "email@google.com",
    //     phone: "123-456-7890",
    //   },
    //   {
    //     firstname: "Name",
    //     lastname: "Name",
    //     profile_photo: "https://api.adorable.io/avatars/40/Name2",
    //     role: "Events",
    //     title: "Title",
    //     email: "email@google.com",
    //     phone: "123-456-7890",
    //     bio: "something something bio bio.",
    //     photos: [
    //       {
    //         name: "alreadysaved.png",
    //       },
    //       {
    //         name: "existing.tif",
    //       }
    //     ],
    //   },
    //   {
    //     firstname: "Mary",
    //     lastname: "Almond",
    //     profile_photo: "https://api.adorable.io/avatars/40/Mary",
    //     role: "Billing",
    //     title: "Billing Administrator",
    //     email: "email@google.com",
    //     phone: "123-456-7890",
    //   },
    //   {
    //     firstname: "Sasha",
    //     lastname: "Maxwell",
    //     profile_photo: "https://api.adorable.io/avatars/40/Sasha",
    //     role: "Signatory",
    //     title: "Legal Services",
    //     email: "email@google.com",
    //     phone: "123-456-7890",
    //   },
    // ];
  }

}
