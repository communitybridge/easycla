/*
 * Copyright The Linux Foundation and each contributor to CommunityBridge.
 * SPDX-License-Identifier: MIT
 */

const handler = require('./add-headers');

describe('addHeaders', () => {
  const EVENT = {
    Records: [
      {
        cf: {
          response: {
            headers: {
              host: [
                {
                  value: 'd123.cf.net',
                  key: 'Host'
                }
              ]
            },
            clientIp: '2001:cdba::3257:9652',
            uri: '/index.html',
            method: 'GET'
          },
          config: {
            distributionId: 'EXAMPLE'
          }
        }
      }
    ]
  };

  it('returns a response object with added headers', () => {
    const headers = {
      'Some-Header-One': '1',
      'Some-Header-Two': '2'
    };
    const output = handler.addHeaders(EVENT, headers);
    expect(output).toEqual({
      headers: {
        host: [
          {
            value: 'd123.cf.net',
            key: 'Host'
          }
        ],
        'some-header-one': [
          {
            value: '1',
            key: 'Some-Header-One'
          }
        ],
        'some-header-two': [
          {
            value: '2',
            key: 'Some-Header-Two'
          }
        ]
      },
      clientIp: '2001:cdba::3257:9652',
      uri: '/index.html',
      method: 'GET'
    });
  });
});
