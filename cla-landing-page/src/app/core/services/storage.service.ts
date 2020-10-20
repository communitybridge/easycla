// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Injectable } from '@angular/core';

@Injectable()
export class StorageService {
  cookiesItems = [];
  constructor(
  ) { }

  getItem<T>(key: string): T {
    const result = localStorage.getItem(key);
    let resultJson = null;
    if (result != null) {
      resultJson = result;
    }
    return resultJson;
  }

  setItem<T>(key: string, value: T) {
    localStorage.setItem(key, JSON.stringify(value));
  }

  getItemFromCookies(key: string) {
    const name = key + '=';
    const ca = document.cookie.split(';');
    for (const i of ca) {
      let c = ca[i];
      while (c.charAt(0) === ' ') {
        c = c.substring(1);
      }
      if (c.indexOf(name) === 0) {
        return c.substring(name.length, c.length);
      }
    }
    return '';
  }

  removeItem<T>(key: string) {
    localStorage.removeItem(key);
  }

  removeAll<T>() {
    localStorage.clear();
  }
}
