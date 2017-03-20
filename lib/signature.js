const aws4 = require('aws4');
const crypto = require('crypto');
const moment = require('moment');
const url = require('url');

module.exports = {

  // default to version one until version two is ready
  signRequest: signRequestVersionFour,

  // allow calling either directly (for testing)
  signRequestVersionOne: signRequestVersionOne,
  signRequestVersionFour: signRequestVersionFour

};

/**
 * Signs a request using the version 1 request signing described at
 * https://confluence.linuxfoundation.org/pages/viewpage.action?pageId=7143516.
 *
 * The intent is to phase this out and begin using the version 2
 * authentication, but remains here in the meantime for use and will
 * likely live on afterwards for reference.
 *
 * @param apiKey a pair of keyId and secret.
 * @param requestOpts the request options.
 * @returns {*} a modified option set.
 */
function signRequestVersionOne(apiKey, requestOpts) {
  let currentDate = moment.utc().toISOString();

  let stringToSign = (requestOpts.method || 'GET') + '\n' +
    url.parse(requestOpts.uri).path + '\n' +
    currentDate + '\n';

  let checksum;

  if (requestOpts.body) {
    checksum = _hashToHex('md5', requestOpts.body);
    stringToSign += checksum + '\n';
  }

  stringToSign += '1';

  let hmac = crypto.createHmac('sha256', apiKey.secret);

  hmac.write(stringToSign);

  let signature = hmac.digest('base64');

  requestOpts.headers = {
    'Content-Type': 'application/json; charset=UTF-8',
    'Date': currentDate,
    'Signature-Version': '1',
    'Authorization': 'CINCO ' + apiKey.keyId + ':' + signature
  };

  if (checksum) {
    requestOpts.headers['Content-MD5'] = checksum;
  }

  return requestOpts;
}

/**
 * Signs a request using the version 2 authentication which is designed
 * to have parity with AWS v4 signatures, described in the AWS documentation
 * at http://docs.aws.amazon.com/general/latest/gr/sigv4-create-canonical-request.html.
 *
 * Most of this is delegated to the aws4 module to handle the signing,
 * rather than doing it manually and probably breaking everything.
 *
 * @param apiKey a pair of keyId and secret.
 * @param requestOpts the request options.
 * @returns {*} a modified option set.
 */
function signRequestVersionFour(apiKey, requestOpts) {
  let parsedUri = url.parse(requestOpts.uri);
  let headers = requestOpts.headers || {};

  if (!headers['Content-Type']) {
    headers['Content-Type'] = 'application/json; charset=UTF-8';
  }

  headers['Signature-Version'] = '4';

  if (requestOpts.body) {
    headers['Content-MD5'] = _hashToHex('md5', requestOpts.body);
  }

  let pathname = parsedUri.pathname;

  if (pathname !== '/' && pathname.endsWith('/')) {
    pathname = pathname.slice(0, -1);
  }

  if (parsedUri.query) {
    pathname += '?' + parsedUri.query;
  }

  let signingOpts = {
    host: parsedUri.host,
    path: pathname,
    method: requestOpts.method,
    headers: headers,
    service: 'cinco',
    region: 'internal',
    body: requestOpts.body,
    doNotEncodePath: true
  };

  let sig = aws4.sign(signingOpts, {
    accessKeyId: apiKey.keyId,
    secretAccessKey: apiKey.secret
  });

  requestOpts.headers = sig.headers;

  return requestOpts;
}

/**
 * Small wrapper to Node.js Crypto to hash variable arguments
 * using a provided algorithm. The hash is digested as hex.
 *
 * @param algorithm the hashing algorithm to use.
 * @param args variable arguments to add to the hash.
 * @returns {String} a hex hash.
 * @private
 */
function _hashToHex(algorithm, ...args) {
  let hash = crypto.createHash(algorithm);
  for (let i = 0; i < args.length; i++) {
    hash.update(args[i]);
  }
  return hash.digest('hex');
}
