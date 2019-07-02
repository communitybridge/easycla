'use strict';

const randomCharacter = () => {
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
  const index = Math.floor(Math.random() * chars.length) % chars.length;
  return chars.charAt(index);
};

/**
 * Generates a random alphanumeric ASCII string of the given length.
 * @param {number} length
 */
function randomString(length) {
  let text = '';
  for (let i = 0; i < length; ++i) {
    text += randomCharacter();
  }
  return text;
}

module.exports = randomString;
