var request = require('request');
var bodyParser = require('body-parser');
var crypto = require('crypto');
var async = require('async');

const integration_user = process.env['CONSOLE_INTEGRATION_USER'];
const integration_pass = process.env['CONSOLE_INTEGRATION_PASSWORD'];

module.exports = function (apiRootUrl) {
  return {
    apiRootUrl: apiRootUrl,

    getKeysForLfId: function (lfId, next) {
      request.get(apiRootUrl + "/auth/trusted/cas/" + lfId, function (err, res, body) {
            if (err) {
              next(err)
            } else if (res.statusCode != 200) {
              next(new Error('Unable to get keys for LfId of [' + lfId + '].  ' +
                  'Status: [' + res.statusCode + '].  Response Body: ' + body));
            } else {
              body = JSON.parse(body);
              next(null, {lfId: lfId, keyId: body.keyId, secret: body.secret})
            }
          }
      ).auth(integration_user, integration_pass, true);
    },


    client: function (apiKeys) {
      return {

      }
    },

    provisionUser: function (payload, appKeyId, appSecret, callback) {
      payload = JSON.stringify(payload);
      var httpMethod = 'POST';
      var uriPath = '/rest/v1/provision/application/user';
      var currentTime = new Date().toISOString();
      var md5 = crypto.createHash('md5').update(payload).digest('hex');
      var signatureVersion = '1';
      var toSign = httpMethod + '\n' + uriPath + '\n' + currentTime + '\n' + md5 + '\n' + signatureVersion;
      var signature = crypto.createHmac('sha1', appSecret).update(toSign).digest('base64')
      request({
        method: httpMethod,
        url: hostURL + uriPath,
        headers: {
          'Content-Type': 'application/json',
          'Date': currentTime,
          'Signature-Version': '1',
          'Content-MD5': md5,
          'Authorization': 'CINCO ' + appKeyId + ': ' + signature
        },
        body: payload
      }, function (error, response) {
        if (!error) {
          callback(error, response);
        }
        else {
          callback(error, response);
          console.log(error);
        }
      });
    },

    getUserKeys: function (payload, appKeyId, appSecret, callback) {
      payload = JSON.stringify(payload);
      var httpMethod = 'POST';
      var uriPath = '/rest/v1/provision/application/auth';
      var currentTime = new Date().toISOString();
      var md5 = crypto.createHash('md5').update(payload).digest('hex');
      var signatureVersion = '1';
      var toSign = httpMethod + '\n' + uriPath + '\n' + currentTime + '\n' + md5 + '\n' + signatureVersion;
      var signature = crypto.createHmac('sha1', appSecret).update(toSign).digest('base64')
      request({
        method: httpMethod,
        url: hostURL + uriPath,
        headers: {
          'Content-Type': 'application/json',
          'Date': currentTime,
          'Signature-Version': '1',
          'Content-MD5': md5,
          'Authorization': 'CINCO ' + appKeyId + ': ' + signature
        },
        body: payload
      }, function (error, response) {
        if (!error) {
          callback(error, response);
        }
        else {
          callback(error, response);
          console.log(error);
        }
      });
    },

    getDeviceId: function () {
      var deviceId = configuration.developerKeys.deviceId;
      return deviceId;
    },

    apiStatus: function (callback) {
      var statusURI = '/rest/status';
      request({
        method: 'GET',
        url: hostURL + statusURI,
      }, function (error, response, body) {
        if (response.statusCode == 200 && body != '[]' && body) {
          var response = JSON.parse(body);
          callback(error, response);
        }
        else {
          console.log('apiStatus: ' + body);
          console.log("apiStatus response: " + response.statusCode);
        }
      });
    }
  }
};






