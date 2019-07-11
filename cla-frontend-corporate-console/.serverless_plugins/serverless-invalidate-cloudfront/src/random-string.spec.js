// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

const randomString = require('./random-string');

describe('random-string', () => {
  let stacks;

  it('will generate a string of the correct length', () => {
    const result = randomString(16);
    expect(result.length).toBe(16);
  });

  it("won't generate the same string twice", () => {
    const result1 = randomString(16);
    const result2 = randomString(16);
    expect(result1).not.toBe(result2);
  });
});
