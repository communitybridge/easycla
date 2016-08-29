var request = require('request');
var bodyParser = require('body-parser');
var crypto = require('crypto');
var async = require('async');

const integration_user = process.env['CONSOLE_INTEGRATION_USER'];
const integration_pass = process.env['CONSOLE_INTEGRATION_PASSWORD'];


var eventAsPostEntity = JSON.stringify(validEvent);
var md5 = checksum(eventAsPostEntity);
var signature = signRequestVersionOne(secretKey,'POST', path, now, md5);

var post_options = {
  body: eventAsPostEntity,
  headers : {
    'Content-MD5' : md5,
    'Content-Type' : 'application/json',
    'Date' : now,
    'Signature-Version' : 1,
    'Authorization' : 'SIDECAR ' + accessKey + ":" + signature
  }
};

var req = request.post(baseUrl + path, post_options, function(error, response, body) {
  assert.equal(response.statusCode,202,"Response Code for Posting an Event was not 200, Was: " + response.statusCode );
  callback();
});




module.exports = function (apiRootUrl) {
  return {
    apiRootUrl: apiRootUrl,

    getKeysForLfId: function (lfId, next) {
      request.get(apiRootUrl + "auth/trusted/cas/" + lfId, function (err, res, body) {
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

      var md5checksum = function() {
        md5sum = crypto.createHash('md5')
        for (var i = 0; i < arguments.length; i++) {
          md5sum.update(arguments[i]);
        }
        return md5sum.digest('hex');
      };

      var signRequestVersionOne = function(secretKey, method, path, dateString, entityMd5) {
        var version = "1";
        var stringToSign = method + '\n' + path + '\n' + dateString + '\n';
        if (entityMd5) {
          stringToSign += entityMd5 + '\n';
        }
        stringToSign += version

        var hmac = crypto.createHmac("sha1", secretKey);
        hmac.write(stringToSign);
        return hmac.digest('base64');
      };

      var makeSignedRequest = function(reqOpts, next) {
        if (reqOpts.body) {
          reqOpts.bodyChecksum = md5checksum(JSON.stringify(reqOpts.body));
        }
        var now = new Date().toISOString();
        var signature = signRequestVersionOne(reqOpts.apiKeys.secret, reqOpts.method, reqOpts.uriPath,
            now, reqOpts.bodyChecksum);

        reqOpts.uri = apiRootUrl + reqOpts.uriPath;
        request({})
      };



      return {
        createUser: function(lfId, next) {
          var body = {
            "lfId": lfId
          };
          var opts = {
            method: 'POST',
            uriPath: 'user/',
            body: body,
            apiKeys: apiKeys
          };
          makeSignedRequest(opts, function(err, res, body) {
            if (err) {
              next(err, false);
            } else if (res.statusCode == 202) {
              next(null,true)
            } else if (res.statusCode == 409) {
              next(null, false);
            } else {
              next(new Error('User with lfId of [' +lfId + '] not created.' +
                'Status: [' + res.code + '].  Response Body: ' + body))
            }
          });

        },
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
