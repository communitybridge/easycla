// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Injectable } from '@angular/core';
declare var require: any;
var TimSort = require('timsort');

@Injectable()
export class SortService {

  toggleSort(config, prop, dataArray) {
    let current_sort = config[prop].sort;
    this.resetSort(config);
    if (current_sort == 'asc') {
      config[prop].sort = 'desc';
    } else {
      config[prop].sort = 'asc';
    }
    let sort = config[prop].sort;
    TimSort.sort(dataArray, this.sort(config[prop].arrayProp, config[prop].sortType, sort));
  }

  resetSort(config) {
    for (var key in config) {
      if (config.hasOwnProperty(key)) {
        config[key].sort = null;
      }
    }
  }

  sort(prop, type, dir = 'asc') {
    let sort = 1; // standard
    if (dir == 'desc') {
      sort = -1; // inverse
    }
    return function(a, b) {
      prop = prop.replace(/\[(\w+)\]/g, '.$1'); // convert indexes to properties
      prop = prop.replace(/^\./, '');           // strip a leading dot
      var props = prop.split('.');
      var props_len = props.length;
      for (var i = 0; i < props_len; ++i) {
        var k = props[i];
        if (k in a) {
            a = a[k];
        } else {
          return;
        }
        if (k in b) {
            b = b[k];
        } else {
          return;
        }
      }
      if (type == 'text') {
        if (a < b) {
          return -1 * sort;
        }
        if (a > b) {
          return 1 * sort;
        }
        return 0;
      }
      if (type == 'number') {
        return (a - b) * sort;
      }
      if (type == 'date') {
        a = new Date(a);
        b = new Date(b);
        return (a - b) * sort;
      }
      /*
        Semantic Version sorting adapted from https://stackoverflow.com/a/6832721
       */
      if (type == 'semver') {
        let lexicographical = false;
        let zeroExtend = true;
        let v1parts = a.split('.');
        let v2parts = b.split('.');

        let isValidPart = function(x) {
          return (lexicographical ? /^\d+[A-Za-z]*$/ : /^\d+$/).test(x);
        }

        if (!v1parts.every(isValidPart) || !v2parts.every(isValidPart)) {
          return NaN;
        }

        if (zeroExtend) {
          while (v1parts.length < v2parts.length) v1parts.push("0");
          while (v2parts.length < v1parts.length) v2parts.push("0");
        }

        if (!lexicographical) {
          v1parts = v1parts.map(Number);
          v2parts = v2parts.map(Number);
        }

        for (let i = 0; i < v1parts.length; ++i) {
          if (v2parts.length == i) {
            return 1 * sort;
          }

          if (v1parts[i] == v2parts[i]) {
            continue;
          }
          else if (v1parts[i] > v2parts[i]) {
            return 1 * sort;
          }
          else {
            return -1 * sort;
          }
        }

        if (v1parts.length != v2parts.length) {
          return -1 * sort;
        }

        return 0;
      }
    }
  }
}
