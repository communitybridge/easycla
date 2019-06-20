// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

exports.addHeaders = function(event, headerList) {
  const { response } = event.Records[0].cf;
  const { headers } = response;

  Object.keys(headerList).forEach(headerName => {
    const headerValue = headerList[headerName];
    headers[headerName.toLowerCase()] = [
      {
        key: headerName,
        value: headerValue
      }
    ];
  });

  return response;
};
