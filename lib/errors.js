/* Public */

/**
 * Generates an error from a response.
 *
 * The error message contains the status code, response
 * body and any custom user messages and/or details.
 *
 * @param res a response object.
 * @param message a starting error message.
 * @returns {Error}
 */
function fromResponse(res, message) {
  if (message) {
    message += ' ';
  } else {
    message = '';
  }

  message += 'Status: [' + res.statusCode + '].';

  if (res.body) {
    message += ' Response body: ' + res.body;
  }

  let error = new Error(message);

  error.statusCode = res.statusCode;
  error.userMessage = "";

  if (res.body) {
    let errorResponse = JSON
      .parse(res.body);

    if (errorResponse.message) {
      error.userMessage += errorResponse.message;
    }

    if (errorResponse.details) {
      error.userMessage += " " + errorResponse.details;
    }
  }

  return error;
}

/* Exports */

module.exports.fromResponse = fromResponse;