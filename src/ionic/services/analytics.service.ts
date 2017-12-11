import { Injectable } from '@angular/core';
import { Http } from '@angular/http';

import 'rxjs/Rx';

import { ANALYTICS_API_URL } from './constants'; // TODO: Make sure ANALYTICS_API_URL corresponds to ANALYTICS URL

@Injectable()
export class AnalyticsService {

  analyticsApiUrl: string;
  apiVersion: string;

  constructor(public http: Http) {
    this.analyticsApiUrl = ANALYTICS_API_URL;
    this.apiVersion = '/v1';
  }

  //////////////////////////////////////////////////////////////////////////////

  /**
  * This service should ONLY contain methods calling Analytics API
  **/

  //////////////////////////////////////////////////////////////////////////////

  /**
  * Get Analytics Metrics
  *
  * As Dec-11-2017 supports the following Metric Type:
  *
  * issues.firststatus.avg
  * issues.firstresponse.avg
  * issues.toclose.avg
  * issues
  * issues.open
  * website.duration
  * code.commits
  *
  * Further details https://github.linuxfoundation.org/Engineering/analytics-service/blob/develop/docs/analytics.api.md
  **/

  getMetrics(metricType, groupBy, tsFrom, tsTo)
  {
    return this.http.get(
      this.analyticsApiUrl +
      this.apiVersion +
      '/projects/*/metrics/' + metricType +
      '?groupby=' + groupBy +
      '&tsFrom=' + tsFrom +
      '&tsTo='+ tsTo)
      .map(res => res.json());
  }

  //////////////////////////////////////////////////////////////////////////////

}
