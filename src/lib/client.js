const reqLib = require('request');
const signature = require('./signature');

/**
 * Creates a new request client using the provided API
 * and polling interval.
 *
 * @param apiKey the apiKey for requests.
 * @param waitFor the time to wait when polling.
 * @constructor
 */
function RequestClient(apiKey, waitFor) {
  this._apiKey  = apiKey;
  this._waitFor = waitFor || 500;
}

/**
 * Polls a location for the result of an async operation.
 *
 * If an error occurs or the task is not found, we return
 * and error. Otherwise we poll until we receive a non-204
 * status code.
 *
 * @param uri the uri location to poll.
 * @param callback the callback to pass a result to.
 */
RequestClient.prototype.poll = function poll(uri, callback) {
  // carry out request using internals
  this.request(uri, (err, res, body) => {
    if (err) {
      return callback(err);
    }

    // task not found
    if (res.statusCode === 404) {
      return callback(new Error('Unable to locate async task result'));
    }

    // parse the body
    let job = JSON.parse(body);

    // task not ready yet
    if (job.status === 'NOT_STARTED' || job.status === 'RUNNING') {
      return this._delayPoll(uri, callback);
    }

    // error in request
    if (job.status === 'ERROR') {
      return callback(new Error(job.error));
    }

    // result received
    callback(err, res, job.result);
  });
};

/**
 * Executes a request using the provided options, after
 * signing the options using the default signature version.
 *
 * If we receive a 202 result and a location header, we'll
 * go into polling mode until we receive an async result.
 *
 * @param opts the request options.
 * @param callback the callback to pass a result to.
 */
RequestClient.prototype.request = function request(opts, callback) {
  // normalize options into an object format
  typeof opts === 'string' && (opts = { uri: opts });

  // ensure method is set to a default
  !opts.method && (opts.method = 'GET');

  // ensure requests are signed
  let signedOpts = signature
    .signRequest(this._apiKey, opts);

  // carry out request with signed options
  reqLib(signedOpts, (err, res, body) => {
    if (err) {
      return callback(err);
    }

    // if we're async, start polling
    if (res.statusCode === 202 && res.headers.location) {
      return this._delayPoll(res.headers.location, callback);
    }

    // result received
    callback(err, res, body);
  });
};

/**
 * Super basic delay handler for executions requiring
 * a polling delay for async requests.
 *
 * @param exec the execution scope.
 * @private
 */
RequestClient.prototype._delay = function _delay(exec) {
  setTimeout(exec, this._waitFor);
};

/**
 * Triggers a delayed poll call to the given URI after
 * waiting the polling delay in advance.
 *
 * @param uri the URI to poll for a result.
 * @param callback the callback to pass a result to.
 * @private
 */
RequestClient.prototype._delayPoll = function _delayPoll(uri, callback) {
  this._delay(() => this.poll(uri, callback));
};

/**
 * Export the Client as the module.
 *
 * @type {RequestClient}
 */
module.exports = RequestClient;
