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
    this.getCommitActivity();
    this.getcommitsDistribution();
    this.getIssuesStatus();
    this.getPrsPipeline();
    this.getIssuesActivity();
    this.getPrsActivity();
    this.getPageViews();
    this.getMaintainers();
    this.redrawCharts();
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

  getCommitActivity() {
    let index = 'hyperledger2';
    let metricType = 'code.commits';
    let groupBy = 'day';
    let tsFrom = '1510430520000';
    let tsTo =   '1514764800000';
    this.analyticsService.getMetrics(index, metricType, groupBy, tsFrom, tsTo).subscribe(metrics => {
      if (metrics) {
        Object.entries(metrics.value).forEach(
          ([key, value]) => {
            if(value) {
              this.commitsActivityChart.dataTable.push([key, value]);
              this.commitsActivityChart = Object.create(this.commitsActivityChart);
            }
          }
        );
      }
    });
  }

  getcommitsDistribution() {
    let index = 'hyperledger2';
    let metricType = 'maintainers';
    let groupBy = 'month,author';
    let tsFrom = '1512150915000';
    let tsTo =   '1514764800000';
    let maintainers;
    let maintainersCommitsTop10 = 0;
    let maintainersCommitsTotal = 0;
    let top10Percentage = 0;
    let restPercentage = 0;
    this.analyticsService.getMetrics(index, metricType, groupBy, tsFrom, tsTo).subscribe(metrics => {
      if (metrics) {
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
              if(i<10){
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
        this.commitsDistributionChart = Object.create(this.commitsDistributionChart);
      }
    });
  }

  getIssuesStatus() {
    let index = 'hyperledger2';
    let metricType = 'issues';
    let groupBy = 'month,issue_status';
    let tsFrom = '1512150915000';
    let tsTo =   '1514764800000';
    let issuesStatus;
    this.analyticsService.getMetrics(index, metricType, groupBy, tsFrom, tsTo).subscribe(metrics => {
      if (metrics) {
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
        this.issuesStatusChart = Object.create(this.issuesStatusChart);
      }
    });
  }

  getPrsPipeline() {
    let index = 'hyperledger2';
    let metricType = 'prs.submitted';
    let groupBy = 'day';
    let tsFrom = '1512490997000';
    let tsTo =   '1514764800000';
    let sumOpenPRs = 0;
    this.analyticsService.getMetrics(index, metricType, groupBy, tsFrom, tsTo).subscribe(metrics => {
      if (metrics) {
        Object.entries(metrics.value).forEach(
          ([key, value]) => {
            sumOpenPRs = sumOpenPRs + value;
          }
        );
        this.prsPipelineChart.dataTable.push(['PRs', sumOpenPRs]);
        this.prsPipelineChart = Object.create(this.prsPipelineChart);
      }
    });
  }

  getIssuesActivity() {
    let index = 'hyperledger2';
    let metricType = 'issues';
    let groupBy = 'day,issue_status';
    let tsFrom = '1511150915000';
    let tsTo =   '1514764800000';
    this.analyticsService.getMetrics(index, metricType, groupBy, tsFrom, tsTo).subscribe(metrics => {
      if (metrics) {
        Object.entries(metrics.value).forEach(
          ([key, value]) => {
            if(value.value['To Do'] && !value.value['Done']) this.issuesActivityChart.dataTable.push([key, value.value['To Do'], 0]);
            if(!value.value['To Do']  && value.value['Done']) this.issuesActivityChart.dataTable.push([key, 0, value.value['Done']]);
            if(value.value['To Do']  && value.value['Done']) this.issuesActivityChart.dataTable.push([key, value.value['To Do'], value.value['Done']]);
          }
        );
        this.issuesActivityChart = Object.create(this.issuesActivityChart);
      }
    });
  }

  getPrsActivity() {
    let index = 'hyperledger2';
    let metricType = 'prs';
    let groupBy = 'month,issue_status';
    let tsFrom = '1388534400000';
    let tsTo =   '1514764800000';
    this.analyticsService.getMetrics(index, metricType, groupBy, tsFrom, tsTo).subscribe(metrics => {
      if (metrics) {
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
        this.prsActivityChart = Object.create(this.prsActivityChart);
      }
    });
  }

  getPageViews() {
    let index = 'hyperledger2';
    let metricType = 'website.duration';
    let groupBy = 'day';
    let tsFrom = '1510430520000';
    let tsTo =   '1514764800000';
    this.analyticsService.getMetrics(index, metricType, groupBy, tsFrom, tsTo).subscribe(metrics => {
      if (metrics) {
        Object.entries(metrics.value).forEach(
          ([key, value]) => {
            if(value) {
              this.pageViewsChart.dataTable.push([key, value]);
              this.pageViewsChart = Object.create(this.pageViewsChart);
            }
          }
        );
      }
    });
  }

  getMaintainers() {
    let index = 'hyperledger2';
    let metricType = 'maintainers';
    let groupBy = 'month,author';
    let tsFrom = '1512150915000';
    let tsTo =   '1514764800000';
    let maintainers;
    this.analyticsService.getMetrics(index, metricType, groupBy, tsFrom, tsTo).subscribe(metrics => {
      if (metrics) {
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
        this.maintainersTable = Object.create(this.maintainersTable);
      }
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
        textStyle:{ color: '#ffffff'},
        gridlines: {
          color: "#FFFFFF"
        },
        baselineColor: '#FFFFFF'
      },
      vAxis: {title: '# of commits'},
      colors: ['#eeeeee'],
      backgroundColor: '#4e92df',
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
        textStyle:{ color: '#888888'},
        gridlines: {
          color: "#FFFFFF"
        },
        baselineColor: '#FFFFFF'
      },
      colors: ['#cccccc','#eeeeee'],
      backgroundColor: '#4e92df',
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
        textStyle:{ color: '#ffffff'},
        gridlines: {
          color: "#FFFFFF"
        },
        baselineColor: '#FFFFFF'
      },
      vAxis: {},
      colors: ['#eeeeee'],
      backgroundColor: '#4e92df',
      legend: 'none',
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
      greenColor: '#4e92df'
    }
  };

  public issuesActivityChart:any =  {
    chartType: 'AreaChart',
    dataTable: [
      ['Date', 'Issues Open', 'Issues Closed']
    ],
    options: {
      hAxis: {
        textStyle:{ color: '#ffffff'},
        gridlines: {
          color: "#FFFFFF"
        },
        baselineColor: '#FFFFFF'
      },
      vAxis: {title: '# of Issues'},
      colors: ['#9344dd', '#abab45'],
      backgroundColor: '#4e92df',
      legend: 'none'
    }
  };

  public prsActivityChart:any =  {
    chartType: 'AreaChart',
    dataTable: [
      ['Date', 'PRs Open', 'PRs Merged', 'PRs Closed']
    ],
    options: {
      hAxis: {
        textStyle:{ color: '#ffffff'},
        gridlines: {
          color: "#FFFFFF"
        },
        baselineColor: '#FFFFFF'
      },
      vAxis: {title: '# of PRs'},
      colors: ['#9344dd', '#abab45'],
      backgroundColor: '#4e92df',
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
        textStyle:{ color: '#ffffff'},
        gridlines: {
          color: "#FFFFFF"
        },
        baselineColor: '#FFFFFF'
      },
      vAxis: {title: 'Page Views (in thousands)', minValue: 0, maxValue: 15},
      colors: ['#95c2e2'],
      backgroundColor: '#4e92df',
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
