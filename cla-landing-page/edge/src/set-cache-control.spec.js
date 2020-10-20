/*
 * Copyright The Linux Foundation and each contributor to CommunityBridge.
 * SPDX-License-Identifier: MIT
 */

const handler = require('./set-cache-control');

describe('setCacheControl', () => {
  it("turns on Cache-Control with max-age when resource isn't in filesNotToCache", () => {
    const headers = {
      'Some-Header-One': '1',
      'Some-Header-Two': '2'
    };
    const filesNotToCache = ['index.html'];
    const timeToLive = 60 * 60 * 24 * 365;

    const result = handler.setCacheControl(headers, 'some-image.png', filesNotToCache, timeToLive);
    expect(result).toEqual({
      ...headers,
      'Cache-Control': `max-age=${timeToLive}`
    });
  });

  it('turns off Cache-Control when resource is in filesNotToCache', () => {
    const headers = {
      'Some-Header-One': '1',
      'Some-Header-Two': '2'
    };
    const filesNotToCache = ['index.html'];
    const timeToLive = 60 * 60 * 24 * 365;

    const result = handler.setCacheControl(headers, 'index.html', filesNotToCache, timeToLive);
    expect(result).toEqual({
      ...headers,
      'Cache-Control': `no-cache, no-store, must-revalidate`
    });
  });

  it("doesn't change the s-maxage property in an existing Cache-Control header", () => {
    const headers = {
      'Some-Header-One': '1',
      'Some-Header-Two': '2',
      'Cache-Control': 's-maxage=100'
    };
    const filesNotToCache = ['index.html'];
    const timeToLive = 60 * 60 * 24 * 365;

    const result = handler.setCacheControl(headers, 'index.html', filesNotToCache, timeToLive);
    expect(result).toEqual({
      ...headers,
      'Cache-Control': `s-maxage=100, no-cache, no-store, must-revalidate`
    });
  });

  it('overrides max-age property in an existing Cache-Control header', () => {
    const headers = {
      'Some-Header-One': '1',
      'Some-Header-Two': '2',
      'Cache-Control': 'max-age=100'
    };
    const filesNotToCache = ['index.html'];
    const timeToLive = 60 * 60 * 24 * 365;

    const result = handler.setCacheControl(headers, 'some-file.html', filesNotToCache, timeToLive);
    expect(result).toEqual({
      ...headers,
      'Cache-Control': `max-age=${timeToLive}`
    });
  });
});
