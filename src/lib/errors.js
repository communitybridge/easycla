/* Public */

/**
 * Creates an error with optional metadata fields.
 *
 * Metadata fields are copied into the error payload
 * after creation, in order to further the details
 * available in the error object.
 *
 * @param msg the base error message.
 * @param meta the metadata to copy over.
 * @returns {Error}
 */
function create(msg, meta = {}) {
  let error = new Error(msg);
  Object.assign(error, meta);
  return error;
}

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

  let userMessage = '';

  if (res.body) {
    let errorResponse = JSON
      .parse(res.body);

    if (errorResponse.message) {
      userMessage += errorResponse.message;
    }

    if (errorResponse.details) {
      userMessage += ' ' + errorResponse.details;
    }
  }

  return create(message, {
    statusCode: res.statusCode,
    userMessage: userMessage
  });
}

/* Exports */

module.exports.create = create;
module.exports.fromResponse = fromResponse;