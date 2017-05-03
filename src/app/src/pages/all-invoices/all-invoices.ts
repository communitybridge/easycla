import { Component, ViewChild } from '@angular/core';

import { NavController, IonicPage } from 'ionic-angular';

import { CincoService } from '../../app/services/cinco.service'
import { Chart } from 'chart.js';

@IonicPage({
  segment: 'all-invoices'
})
@Component({
  selector: 'all-invoices',
  templateUrl: 'all-invoices.html'
})
export class AllInvoicesPage {
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
      },
      {
        id: 'm00000000004',
        project: 'TODO',
        name: 'ADP',
        level: 'Gold',
        status: 'Invoice Sent',
        annual_dues: '$30,000',
        renewal_date: '4/1/2017',
      },
      {
        id: 'm00000000005',
        project: 'Zephyr',
        name: 'BlackRock',
        level: 'Bronze',
        status: 'Renewal < 60 Days',
        annual_dues: '$30,000',
        renewal_date: '6/1/2017',
      },
      {
        id: 'm00000000006',
        project: 'OpenSwitch',
        name: 'Fox',
        level: 'Bronze',
        status: 'Renewal < 60 Days',
        annual_dues: '$30,000',
        renewal_date: '10/1/2017',
      },
      {
        id: 'm00000000007',
        project: 'TODO',
        name: 'Google',
        level: 'Gold',
        status: 'Renewal < 60 Days',
        annual_dues: '$30,000',
        renewal_date: '10/1/2017',
      },
      {
        id: 'm00000000008',
        project: 'TODO',
        name: 'Joyent',
        level: 'Gold',
        status: 'Renewal < 60 Days',
        annual_dues: '$30,000',
        renewal_date: '10/1/2017',
      },
      {
        id: 'm00000000009',
        project: 'TODO',
        name: 'KrVolk',
        level: 'Gold',
        status: 'Renewal < 60 Days',
        annual_dues: '$30,000',
        renewal_date: '10/1/2017',
      },
      {
        id: 'm00000000010',
        project: 'TODO',
        name: 'Netflix',
        level: 'Gold',
        status: 'Renewal < 60 Days',
        annual_dues: '$30,000',
        renewal_date: '10/1/2017',
      },
      {
        id: 'm00000000011',
        project: 'TODO',
        name: 'Company Name',
        level: 'Silver',
        status: 'Renewal < 60 Days',
        annual_dues: '$30,000',
        renewal_date: '10/1/2017',
      },
    ];
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

}
