// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import { Component, ViewChild } from '@angular/core';

import { NavController, IonicPage } from 'ionic-angular';

import { CincoService } from '../../services/cinco.service'
import { KeycloakService } from '../../services/keycloak/keycloak.service';
import { Chart } from 'chart.js';

import { ProjectModel } from '../../models/project-model';
import { RolesService } from '../../services/roles.service';
import { Restricted } from '../../decorators/restricted';

@Restricted({
  roles: ['isAuthenticated', 'isPmcUser'],
})
// @IonicPage({
//   segment: 'all-invoices'
// })
@Component({
  selector: 'all-invoices',
  templateUrl: 'all-invoices.html'
})
export class AllInvoicesPage {

  allProjects: any;

  project = new ProjectModel();

  allProjectsWithMembers: any[] = [];

  membersCount: any;

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

  constructor(
    public navCtrl: NavController,
    private cincoService: CincoService,
    private keycloak: KeycloakService,
    public rolesService: RolesService,
  ) {
    this.getDefaults();
  }

  async ngOnInit(){
    this.getAllProjects();
  }

  getAllProjects() {
    this.cincoService.getAllProjects().subscribe(response => {
      if(response) {
        this.allProjects = response;
        for(let eachProject of this.allProjects) {
          this.getProject(eachProject.id);
        }
      }
    });
  }

  getProject(projectId) {
    let getMembers = true;
    this.cincoService.getProject(projectId, getMembers).subscribe(response => {
      if(response) {
        this.allProjectsWithMembers.push(response);
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
  }

}
