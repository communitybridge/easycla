import { Component } from '@angular/core';

import { NavController, NavParams } from 'ionic-angular';

@Component({
  selector: 'page-project',
  templateUrl: 'project.html'
})
export class ProjectPage {
  selectedProject: any;
  project: {
    icon: string,
    title: string,
    description: string,
    datas: Array<{
      label: string,
      value: string,
    }>,
  };

  members: Array<{
    alert?: string,
    name: string,
    level: string,
    status: string,
    annual_dues: string,
    renewal_date: string,
  }>;

  constructor(public navCtrl: NavController, public navParams: NavParams) {
    // If we navigated to this page, we will have an item available as a nav param
    this.selectedProject = navParams.get('project');

    // use selectedProject to get data from CINCO and populate this.project

    this.project = {
      icon: "https://dummyimage.com/600x250/ffffff/000.png&text=zephyr+logo",
      title: "Zephyr",
      description: "Zephyr Project is a small, scalable, real-time operating system for use on resource-constraned systems supporting multiple architectures...",
      datas: [
        {
          label: "Budget",
          value: "$3,000,000",
        },
        {
          label: "Categories",
          value: "Embedded & IoT",
        },
        {
          label: "Launched",
          value: "5/1/2016",
        },
        {
          label: "Current",
          value: "$2,000,000 ($1,000,000)",
        },
        {
          label: "Members",
          value: "41",
        },
      ],
    };

    this.members = [
      {
        alert: '',
        name: 'Abbie',
        level: 'Gold',
        status: 'Invoice Paid',
        annual_dues: '$30,000',
        renewal_date: '3/1/2017',
      },
      {
        alert: 'alert',
        name: 'Acrombie',
        level: 'Gold',
        status: 'Invoice Sent (Late)',
        annual_dues: '$30,000',
        renewal_date: '3/2/2017',
      },
      {
        alert: 'notice',
        name: 'Adobe',
        level: 'Gold',
        status: 'Contract: Pending',
        annual_dues: '$30,000',
        renewal_date: '4/1/2017',
      },
      {
        alert: '',
        name: 'ADP',
        level: 'Gold',
        status: 'Invoice Sent',
        annual_dues: '$30,000',
        renewal_date: '4/1/2017',
      },
      {
        alert: '',
        name: 'BlackRock',
        level: 'Bronze',
        status: 'Renewal < 60',
        annual_dues: '$30,000',
        renewal_date: '6/1/2017',
      },
      {
        alert: '',
        name: 'Fox',
        level: 'Bronze',
        status: 'Renewal < 60',
        annual_dues: '$30,000',
        renewal_date: '10/1/2017',
      },
      {
        alert: '',
        name: 'Google',
        level: 'Gold',
        status: 'Renewal < 60',
        annual_dues: '$30,000',
        renewal_date: '10/1/2017',
      },
      {
        alert: '',
        name: 'Joyent',
        level: 'Gold',
        status: 'Renewal < 60',
        annual_dues: '$30,000',
        renewal_date: '10/1/2017',
      },
      {
        alert: '',
        name: 'KrVolk',
        level: 'Gold',
        status: 'Renewal < 60',
        annual_dues: '$30,000',
        renewal_date: '10/1/2017',
      },
      {
        alert: '',
        name: 'Netflix',
        level: 'Gold',
        status: 'Renewal < 60',
        annual_dues: '$30,000',
        renewal_date: '10/1/2017',
      },
      {
        alert: '',
        name: 'Company Name',
        level: 'Silver',
        status: 'Renewal < 60',
        annual_dues: '$30,000',
        renewal_date: '10/1/2017',
      },
    ];

  }

  memberSelected(event, project, member) {
    alert("check the console!");
    console.log({project, member});
    // this.navCtrl.push(MemberPage, {
    //   project: project
    // });
  }
}
