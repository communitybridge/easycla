import { Component, ViewChild } from '@angular/core';

import { NavController, IonicPage } from 'ionic-angular';

import { CincoService } from '../../app/services/cinco.service'
import { Chart } from 'chart.js';

import { ProjectModel } from '../../models/project-model';

@IonicPage({
  segment: 'all-invoices'
})
@Component({
  selector: 'all-invoices',
  templateUrl: 'all-invoices.html'
})
export class AllInvoicesPage {

  allProjects: any;

  project = new ProjectModel();

  membersCount: any;

  members: Array<{
    id: string,
    project: string,
    name: string,
    level: string,
    status: string,
    annual_dues: string,
    renewal_date: string,
  }>;
  upcomingInvoices: {
    base: Array<number>,
    additional: Array<number>,
    alert: Array<number>,
  };
  sentInvoices: {
    base: Array<number>,
    additional: Array<number>,
    alert: Array<number>,
  };
  chartColors: {
    base: {
      background: string,
      hoverBackground: string,
    },
    additional: {
      background: string,
      hoverBackground: string,
    },
    alert: {
      background: string,
      hoverBackground: string,
    },
  };

  @ViewChild('upcomingInvoicesCanvas') upcomingInvoicesCanvas;
  upcomingInvoicesChart: any;

  @ViewChild('sentInvoicesCanvas') sentInvoicesCanvas;
  sentInvoicesChart: any;

  constructor(public navCtrl: NavController, private cincoService: CincoService) {
    this.getDefaults();
  }

  ngOnInit(){
    this.getAllProjects();
  }

  getAllProjects(){
    this.cincoService.getAllProjects().subscribe(response => {
      console.log(response);
      this.allProjects = response;
      // this.getProject(this.projectId);
    });
  }

  getProject(projectId) {
    let getMembers = true;
    this.cincoService.getProject(projectId, getMembers).subscribe(response => {
      if(response) {
        this.project.id = response.id;
        this.project.name = response.name;
        this.project.description = response.description;
        this.project.managers = response.managers;
        this.project.status = response.status;
        this.project.category = response.category;
        this.project.sector = response.sector;
        this.project.url = response.url;
        this.project.startDate = response.startDate;
        this.project.logoRef = response.logoRef;
        this.project.agreementRef = response.agreementRef;
        this.project.mailingListType = response.mailingListType;
        this.project.emailAliasType = response.emailAliasType;
        this.project.address = response.address;
        this.project.members = response.members;
        this.membersCount = this.project.members.length;
      }
    });
  }

  ionViewDidLoad() {
    var barOptions_stacked = {
      tooltips: {
          enabled: true
      },
      hover :{
          animationDuration:0
      },
      scales: {
          xAxes: [{
              ticks: {
                  beginAtZero:true,
                  fontSize:11
              },
              scaleLabel:{
                  display:false
              },
              gridLines: {
              },
              stacked: true
          }],
          yAxes: [{
              gridLines: {
                  display:false,
                  color: "#fff",
                  zeroLineColor: "#fff",
                  zeroLineWidth: 0
              },
              ticks: {
                  fontSize:11
              },
              stacked: true
          }]
      },
      legend:{
          display: false
      },
    };

    this.upcomingInvoicesChart = new Chart(this.upcomingInvoicesCanvas.nativeElement, {
      type: 'bar',
      data: {
        labels: ["< 60 Days", "< 30 Days"],
        // Base
        datasets: [{
          data: this.upcomingInvoices.base,
          backgroundColor: this.chartColors.base.background,
          hoverBackgroundColor: this.chartColors.base.hoverBackground
        },
        // Additional
        {
          data: this.upcomingInvoices.additional,
          backgroundColor: this.chartColors.additional.hoverBackground,
          hoverBackgroundColor: this.chartColors.additional.background
        },
        // Alert
        {
          data: this.upcomingInvoices.alert,
          backgroundColor: this.chartColors.alert.background,
          hoverBackgroundColor: this.chartColors.alert.hoverBackground
        },
      ]
      },
      options: barOptions_stacked,
    });

    this.sentInvoicesChart = new Chart(this.sentInvoicesCanvas.nativeElement, {
      type: 'bar',
      data: {
        labels: ["SENT", "LATE", "PAID (Last 30 Days)"],
        // Base
        datasets: [{
          label: 'Base',
          data: this.sentInvoices.base,
          backgroundColor: this.chartColors.base.background,
          hoverBackgroundColor: this.chartColors.base.hoverBackground
        },
        // Additional
        {
          label: 'Additional',
          data: this.sentInvoices.additional,
          backgroundColor: this.chartColors.additional.hoverBackground,
          hoverBackgroundColor: this.chartColors.additional.background
        },
        // Alert
        {
          label: 'Alert',
          data: this.sentInvoices.alert,
          backgroundColor: this.chartColors.alert.background,
          hoverBackgroundColor: this.chartColors.alert.hoverBackground
        },
      ]
      },

      options: barOptions_stacked,
    });
  }

  getDefaults() {
    this.upcomingInvoices = {
      base: [4, 4],
      additional: [0, 0],
      alert: [0, 0],
    };
    this.sentInvoices = {
      base: [2, 0, 0],
      additional: [0, 0, 1],
      alert: [0, 1, 0],
    };

    this.chartColors = {
      base: {
        background: "rgba(136,186,22,1)",
        hoverBackground: "rgba(50,90,100,1)",
      },
      additional: {
        background: "rgba(64,86,45,1)",
        hoverBackground: "rgba(64,116,5,1)",
      },
      alert: {
        background: "rgb(245,166,35)",
        hoverBackground: "rgb(213,107,24)",
      },
    };

    this.members = [
      {
        id: 'm00000000001',
        project: 'Zephyr',
        name: 'Abbie',
        level: 'Gold',
        status: 'Invoice Paid',
        annual_dues: '$30,000',
        renewal_date: '3/1/2017',
      },
      {
        id: 'm00000000002',
        project: 'Zephyr',
        name: 'Acrombie',
        level: 'Gold',
        status: 'Invoice Sent (Late)',
        annual_dues: '$30,000',
        renewal_date: '3/2/2017',
      },
      {
        id: 'm00000000003',
        project: 'Zephyr',
        name: 'Adobe',
        level: 'Gold',
        status: 'Contract: Pending',
        annual_dues: '$30,000',
        renewal_date: '4/1/2017',
      }
    ];
  }

}
