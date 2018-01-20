import { Component } from '@angular/core';
import { NavController, ModalController, NavParams, IonicPage } from 'ionic-angular';
import { CincoService } from '../../../services/cinco.service';
import { KeycloakService } from '../../../services/keycloak/keycloak.service';
import { AnalyticsService } from '../../../services/analytics.service';
import { DomSanitizer} from '@angular/platform-browser';
import { RolesService } from '../../../services/roles.service';
import { Restricted } from '../../../decorators/restricted';
import { HostListener } from '@angular/core'

@Restricted({
  roles: ['isAuthenticated', 'isPmcUser'],
})
@IonicPage({
  segment: 'project/:projectId/analytics'
})
@Component({
  selector: 'project-analytics',
  templateUrl: 'project-analytics.html',
  providers: [CincoService, AnalyticsService]
})

export class ProjectAnalyticsPage {

  projectId: string;

  hasAnalyticsUrl: boolean;
  analyticsUrl: any;
  sanitizedAnalyticsUrl: any;

  index:any;
  timeNow:any;
  span:any;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    private cincoService: CincoService,
    private keycloak: KeycloakService,
    private analyticsService: AnalyticsService,
    private domSanitizer : DomSanitizer,
    public modalCtrl: ModalController,
    public rolesService: RolesService,
  ) {
    this.projectId = navParams.get('projectId');
  }

  ngOnInit() {
    this.getProjectConfig(this.projectId);
    this.getDefaults();
  }

  getDefaults() {
    this.setTimeNow();
    this.span = 'month';
    this.index = 'hyperledger';
    this.getCommitActivity(this.span);
    this.getcommitsDistribution(this.span);
    this.getIssuesStatus(this.span);
    this.getIssuesActivity(this.span);
    this.getPrsPipeline(this.span);
    this.getPrsActivity(this.span);
    this.getPageViews(this.span);
    this.getMaintainers(this.span);
    this.redrawCharts();
  }

  setTimeNow() {
    this.timeNow = new Date().getTime();
  }

  calculateTsFrom(span) {
    let rest;
    if(span == 'year') { rest = 365; }
    else if(span == 'quarter') { rest = 90; }
    else if(span == 'month') { rest = 30; }
    else if(span == 'week') { rest = 7; }
    else if(span == 'day') { rest = 1; }
    else { rest = 30; } // otherwise query to a month
    let date = new Date();
    let previousDate = date.getDate() - rest;
    date.setDate(previousDate);
    let tsFrom = date.getTime();
    return tsFrom;
  }

  getProjectConfig(projectId) {
    this.cincoService.getProjectConfig(projectId).subscribe(response => {
      if (response) {
        let projectConfig = response;
        if(projectConfig.analyticsUrl) {
          this.analyticsUrl = projectConfig.analyticsUrl;
          this.sanitizedAnalyticsUrl = this.domSanitizer.bypassSecurityTrustResourceUrl(this.analyticsUrl);
          this.hasAnalyticsUrl = true;
        }
        else{
          this.hasAnalyticsUrl = true;
        }
      }
    });
  }

  openAnaylticsConfigModal(projectId) {
    let modal = this.modalCtrl.create('AnalyticsConfigModal', {
      projectId: projectId,
    });
    modal.onDidDismiss(analyticsUrl => {
      if(analyticsUrl){
        this.analyticsUrl = analyticsUrl;
        this.hasAnalyticsUrl = true;
        this.sanitizedAnalyticsUrl = this.domSanitizer.bypassSecurityTrustResourceUrl(this.analyticsUrl);
      }
    });
    modal.present();
  }

  getCommitActivity(span) {
    let index = this.index;
    let metricType = 'code.commits';
    let groupBy = 'day';
    let tsFrom = this.calculateTsFrom(span);
    let tsTo = this.timeNow;
    this.analyticsService.getMetrics(index, metricType, groupBy, tsFrom, tsTo).subscribe(metrics => {
      this.commitsActivityChart.dataTable = [
        ['Date', 'Commits'] // Clean Array
      ];
      if(Object.keys(metrics.value).length) {
        Object.entries(metrics.value).forEach(
          ([key, value]) => {
            this.commitsActivityChart.dataTable.push([key, value]);
          }
        );
      }
      else {
        this.commitsActivityChart.dataTable.push(['No commits for a ' + span + ' now', 0]);
      }
      this.commitsActivityChart = Object.create(this.commitsActivityChart);
    });
  }

  getcommitsDistribution(span) {
    let index = this.index;
    let metricType = 'maintainers';
    let groupBy = 'year,author';
    let tsFrom = this.calculateTsFrom(span);
    let tsTo = this.timeNow;
    let maintainers;
    let maintainersCommitsTop10 = 0;
    let maintainersCommitsTotal = 0;
    let top10Percentage = 0;
    let restPercentage = 0;
    this.analyticsService.getMetrics(index, metricType, groupBy, tsFrom, tsTo).subscribe(metrics => {
      this.commitsDistributionChart.dataTable = [
        ['Date', 'Commits'] // Clean Array
      ];
      if (Object.keys(metrics.value).length) {
        Object.entries(metrics.value).forEach(
          ([key, value]) => {
            if(value) {
              maintainers = value;
            }
          }
        );
        let i = 0;
        Object.entries(maintainers.value).forEach(
          ([key, value]) => {
            if(value) {
              if(i < 10) {
                maintainersCommitsTop10 = maintainersCommitsTop10 + value;
              }
              maintainersCommitsTotal = maintainersCommitsTotal + value;
              i++
            }
          }
        );
        top10Percentage = Math.round( maintainersCommitsTop10 * 100 / maintainersCommitsTotal );
        restPercentage = 100 - top10Percentage;
        this.commitsDistributionChart.dataTable.push(['Top 10', top10Percentage])
        this.commitsDistributionChart.dataTable.push(['Rest', restPercentage]);
      }
      else {
        this.commitsDistributionChart.dataTable.push(['No commits for a ' + span + ' now', 100]);
      }
      this.commitsDistributionChart = Object.create(this.commitsDistributionChart);
    });
  }

  getIssuesStatus(span) {
    let index = this.index;
    let metricType = 'issues';
    let groupBy = 'year,issue_status';
    let tsFrom = this.calculateTsFrom(span);
    let tsTo = this.timeNow;
    let issuesStatus;
    this.analyticsService.getMetrics(index, metricType, groupBy, tsFrom, tsTo).subscribe(metrics => {
      this.issuesStatusChart.dataTable = [
        ['Status', 'Issues'] // Clean Array
      ];
      if (Object.keys(metrics.value).length) {
        Object.entries(metrics.value).forEach(
          ([key, value]) => {
            issuesStatus = value;
          }
        );
        Object.entries(issuesStatus.value).forEach(
          ([key, value]) => {
            this.issuesStatusChart.dataTable.push([key, value]);
          }
        );
      }
      else {
        this.issuesStatusChart.dataTable.push(['No issues for a ' + span + ' now', 0]);
      }
      this.issuesStatusChart = Object.create(this.issuesStatusChart);
    });
  }

  getIssuesActivity(span) {
    let index = this.index;
    let metricType = 'issues';
    let groupBy = 'day,issue_status';
    let tsFrom = this.calculateTsFrom(span);
    let tsTo = this.timeNow;
    this.analyticsService.getMetrics(index, metricType, groupBy, tsFrom, tsTo).subscribe(metrics => {
      this.issuesActivityChart.dataTable = [
        ['Date', 'Issues Open', 'Issues Closed'] // Clean Array
      ];
      if (Object.keys(metrics.value).length) {
        Object.entries(metrics.value).forEach(
          ([key, value]) => {
            if(value.value['open'] && !value.value['closed']) this.issuesActivityChart.dataTable.push([key, value.value['open'], 0]);
            if(!value.value['open']  && value.value['closed']) this.issuesActivityChart.dataTable.push([key, 0, value.value['closed']]);
            if(value.value['open']  && value.value['closed']) this.issuesActivityChart.dataTable.push([key, value.value['open'], value.value['closed']]);
          }
        );
      }
      else {
        this.issuesActivityChart.dataTable.push(['No issues for a ' + span + ' now', 0, 0]);
      }
      this.issuesActivityChart = Object.create(this.issuesActivityChart);
    });
  }

  getPrsPipeline(span) {
    span = 'year';
    let index = this.index;
    let metricType = 'prs.open';
    let groupBy = 'year';
    let tsFrom = this.calculateTsFrom(span);
    let tsTo = this.timeNow;
    let sumOpenPRs = 0;
    this.analyticsService.getMetrics(index, metricType, groupBy, tsFrom, tsTo).subscribe(metrics => {
      this.prsPipelineChart.dataTable = [
        ['Label', 'Value'] // Clean Array
      ];
      if (Object.keys(metrics.value).length) {
        Object.entries(metrics.value).forEach(
          ([key, value]) => {
            sumOpenPRs = sumOpenPRs + value;
          }
        );
        this.prsPipelineChart.dataTable.push(['PRs', sumOpenPRs]);
      }
      else {
        this.prsPipelineChart.dataTable.push(['No PRs for a ' + span + ' now', 0]);
      }
      this.prsPipelineChart = Object.create(this.prsPipelineChart);
    });
  }

  getPrsActivity(span) {
    let index = this.index;
    let metricType = 'prs';
    let groupBy = 'day,issue_status';
    let tsFrom = this.calculateTsFrom(span);
    let tsTo = this.timeNow;
    this.analyticsService.getMetrics(index, metricType, groupBy, tsFrom, tsTo).subscribe(metrics => {
      this.prsActivityChart.dataTable = [
        ['Date', 'PRs Open', 'PRs Merged', 'PRs Closed'] // Clean Array
      ];
      if (Object.keys(metrics.value).length) {
        Object.entries(metrics.value).forEach(
          ([key, value]) => {
            if(value.value.open && !value.value.merged && !value.value.closed) this.prsActivityChart.dataTable.push([key, value.value.open, 0, 0]);
            if(value.value.merged  && !value.value.open && !value.value.closed) this.prsActivityChart.dataTable.push([key, 0, value.value.merged, 0]);
            if(value.value.closed && !value.value.open && !value.value.merged) this.prsActivityChart.dataTable.push([key, 0, 0, value.value.closed]);
            if(value.value.open && value.value.merged && !value.value.closed) this.prsActivityChart.dataTable.push([key, value.value.open, value.value.merged, 0]);
            if(value.value.merged  && !value.value.open && value.value.closed) this.prsActivityChart.dataTable.push([key, 0, value.value.merged, value.value.closed]);
            if(value.value.closed && value.value.open && !value.value.merged) this.prsActivityChart.dataTable.push([key, value.value.open, 0, value.value.closed]);
            if(value.value.open && value.value.merged && value.value.closed) this.prsActivityChart.dataTable.push([key, value.value.open, value.value.merged, value.value.closed]);
          }
        );
      }
      else {
        this.prsActivityChart.dataTable.push(['No PRs for a ' + span + ' now', 0, 0, 0]);
      }
      this.prsActivityChart = Object.create(this.prsActivityChart);
    });
  }

  getPageViews(span) {
    // let index = this.index;
    //
    // let metricType = 'code.commits';
    // let groupBy = 'day';
    // let tsFrom = this.calculateTsFrom(span);
    // let tsTo = this.timeNow;
    // this.analyticsService.getMetrics(index, metricType, groupBy, tsFrom, tsTo).subscribe(metrics => {
    //   if (metrics) {
    //     Object.entries(metrics.value).forEach(
    //       ([key, value]) => {
    //         if(value) {
    //           this.pageViewsChart.dataTable.push([key, value]);
    //           this.pageViewsChart = Object.create(this.pageViewsChart);
    //         }
    //       }
    //     );
    //   }
    // });


    let index = this.index;
    //TODO: To query actual Page Views. EP doens't exist yet.
    let metricType = 'code.commits';
    let groupBy = 'day';
    let tsFrom = this.calculateTsFrom(span);
    let tsTo = this.timeNow;
    this.analyticsService.getMetrics(index, metricType, groupBy, tsFrom, tsTo).subscribe(metrics => {
      this.pageViewsChart.dataTable = [
        ['Date', 'Views'] // Clean Array
      ];
      if(Object.keys(metrics.value).length) {
        Object.entries(metrics.value).forEach(
          ([key, value]) => {
            this.pageViewsChart.dataTable.push([key, value]);
          }
        );
      }
      else {
        this.pageViewsChart.dataTable.push(['No page views for a ' + span + ' now', 0]);
      }
      this.pageViewsChart = Object.create(this.pageViewsChart);
    });

  }

  getMaintainers(span) {
    let index = this.index;
    let metricType = 'maintainers';
    let groupBy = 'year,author';
    let tsFrom = this.calculateTsFrom(span);
    let tsTo = this.timeNow;
    let maintainers;
    this.analyticsService.getMetrics(index, metricType, groupBy, tsFrom, tsTo).subscribe(metrics => {
      this.maintainersTable.dataTable = [
        ['Name <Email Address>', 'Commits'] // Clean Array
      ];
      if (Object.keys(metrics.value).length) {
        Object.entries(metrics.value).forEach(
          ([key, value]) => {
            if(value) {
              maintainers = value;
            }
          }
        );
        Object.entries(maintainers.value).forEach(
          ([key, value]) => {
            if(value) {
              this.maintainersTable.dataTable.push([key, value]);
            }
          }
        );
      }
      else {
        this.maintainersTable.dataTable.push(['No maintainers for a ' + span + ' now', 0]);
      }
      this.maintainersTable = Object.create(this.maintainersTable);
    });
  }

  redrawCharts() {
    this.commitsActivityChart = Object.create(this.commitsActivityChart);
    this.commitsDistributionChart = Object.create(this.commitsDistributionChart);
    this.issuesStatusChart = Object.create(this.issuesStatusChart);
    this.prsPipelineChart = Object.create(this.prsPipelineChart);
    this.issuesActivityChart = Object.create(this.issuesActivityChart);
    this.prsActivityChart = Object.create(this.prsActivityChart);
    this.pageViewsChart = Object.create(this.pageViewsChart);
  }

  @HostListener('window:resize', ['$event'])
  onResize(event) {
    event.target.innerWidth;
    this.redrawCharts();
  }

  public commitsActivityChart:any =  {
    chartType: 'ColumnChart',
    dataTable: [
      ['Date', 'Commits']
    ],
    options: {
      hAxis: {
        textStyle:{ color: '#0b4e73'},
        gridlines: {
          color: "#0b4e73"
        },
        baselineColor: '#0b4e73'
      },
      vAxis: {title: '# of commits'},
      colors: ['#2bb3e2'],
      backgroundColor: '#ffffff',
      legend: 'none',
    }
  };

  public commitsDistributionChart:any =  {
    chartType: 'PieChart',
    dataTable: [
      ['Date', 'Commits']
    ],
    options: {
      hAxis: {
        textStyle:{ color: '#0b4e73'},
        gridlines: {
          color: "#0b4e73"
        },
        baselineColor: '#0b4e73'
      },
      colors: ['#0b4e73','#2bb3e2'],
      backgroundColor: '#ffffff',
      legend: 'none'
    }
  };

  public issuesStatusChart:any =  {
    chartType: 'BarChart',
    dataTable: [
      ['Status', 'Issues']
    ],
    options: {
      hAxis: {
        textStyle:{ color: '#0b4e73'},
        gridlines: {
          color: "#0b4e73"
        },
        baselineColor: '#0b4e73'
      },
      vAxis: {},
      colors: ['#2bb3e2'],
      backgroundColor: '#ffffff',
      legend: 'none',
    }
  };

  public issuesActivityChart:any =  {
    chartType: 'AreaChart',
    dataTable: [
      ['Date', 'Issues Open', 'Issues Closed']
    ],
    options: {
      hAxis: {
        textStyle:{ color: '#0b4e73'},
        gridlines: {
          color: "#0b4e73"
        },
        baselineColor: '#FFFFFF'
      },
      vAxis: {title: '# of Issues'},
      colors: ['#2bb3e2', '#0b4e73'],
      backgroundColor: '#ffffff',
      legend: 'none'
    }
  };

  public prsPipelineChart:any =  {
    chartType: 'Gauge',
    dataTable: [
      ['Label', 'Value']
    ],
    options: {
      animation: {easing: 'out'},
      greenFrom: 0, greenTo: 32,
      minorTicks: 1,
      min: 0, max: 100,
      majorTicks: ['0', '20', '40', '60', '80', '100'],
      greenColor: '#ffffff'
    }
  };

  public prsActivityChart:any =  {
    chartType: 'AreaChart',
    dataTable: [
      ['Date', 'PRs Open', 'PRs Merged', 'PRs Closed']
    ],
    options: {
      hAxis: {
        textStyle:{ color: '#0b4e73'},
        gridlines: {
          color: "#0b4e73"
        },
        baselineColor: '#0b4e73'
      },
      vAxis: {title: '# of PRs'},
      colors: ['#2bb3e2', '#0b4e73'],
      backgroundColor: '#ffffff',
      legend: 'none'
    }
  };

  public pageViewsChart:any =  {
    chartType: 'AreaChart',
    dataTable: [
      ['Date', 'Page Views']
    ],
    options: {
      hAxis: {
        textStyle:{ color: '#0b4e73'},
        gridlines: {
          color: "#0b4e73"
        },
        baselineColor: '#0b4e73'
      },
      vAxis: {title: 'Page Views (in thousands)', minValue: 0, maxValue: 15},
      colors: ['#2bb3e2'],
      backgroundColor: '#ffffff',
      legend: 'none'
    }
  };

  public cssClassNames:any = {
    'headerRow': 'cssHeaderRow',
    'tableRow': 'cssTableRow',
    'oddTableRow': 'cssOddTableRow',
    'selectedTableRow': 'cssSelectedTableRow',
    'hoverTableRow': 'cssHoverTableRow',
    'headerCell': 'cssHeaderCell',
    'tableCell': 'cssTableCell',
    'rowNumberCell': 'cssRowNumberCell'
  };

  public maintainersTable:any =  {
    chartType: 'Table',
    dataTable: [
      ['Name <Email Address>', 'Commits'],
    ],
    options: {
      title: 'Maintainers',
      allowHtml: true,
      alternatingRowStyle: false,
      width: '100%',
      cssClassNames: this.cssClassNames,
      sortColumn: 1,
      sortAscending: false,
    }
  };

}
