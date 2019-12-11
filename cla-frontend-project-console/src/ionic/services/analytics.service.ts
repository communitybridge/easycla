// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Injectable } from '@angular/core';
import { Http } from '@angular/http';

import 'rxjs/Rx';
import { EnvConfig } from './cla.env.utils';

@Injectable()
export class AnalyticsService {
  analyticsApiUrl: string;
  apiVersion: string;

  constructor(public http: Http) {
    this.analyticsApiUrl = EnvConfig['analytics-api-url'];
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
   * prs
   * prs.submitted
   * prs.open
   * prs.merged
   * prs.closed
   *
   * Further details https://github.linuxfoundation.org/Engineering/analytics-service/blob/develop/docs/analytics.api.md
   **/

  getMetrics(index, metricType, groupBy, tsFrom, tsTo) {
    return this.http
      .get(
        this.analyticsApiUrl +
          this.apiVersion +
          '/projects/*/metrics/' +
          metricType +
          '?groupby=' +
          groupBy +
          '&tsFrom=' +
          tsFrom +
          '&tsTo=' +
          tsTo +
          '&index=/' +
          index +
          '/docs/'
      )
      .map(res => res.json());
  }

  //////////////////////////////////////////////////////////////////////////////
}
