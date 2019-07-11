// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

const matchDistributions = require('./match-distributions');

describe('matchDistributions', () => {
  let stacks;

  beforeEach(() => {
    stacks = [
      {
        ResourceType: 'AWS::CloudFront::Distribution',
        LogicalResourceId: 'dist1',
        PhysicalResourceId: 'physical1'
      },
      {
        ResourceType: 'AWS::CloudFront::Distribution',
        LogicalResourceId: 'dist2',
        PhysicalResourceId: 'physical2'
      },
      {
        ResourceType: 'AWS::Lambda',
        LogicalResourceId: 'nondistribution',
        PhysicalResourceId: 'physical3'
      }
    ];
  });

  it('can match a single distribution to a stack resource', () => {
    const distributions = {
      dist1: {
        paths: ['a', 'b', 'c']
      }
    };
    const result = matchDistributions(distributions, stacks, 'staging');
    expect(result).toEqual([{ distribution: distributions.dist1, distributionId: 'physical1', name: 'dist1' }]);
  });
  it('can match multiple distributions to a stack resource', () => {
    const distributions = {
      dist1: {
        paths: ['a', 'b', 'c']
      },
      dist2: {
        paths: ['e', 'f', 'g']
      }
    };
    const result = matchDistributions(distributions, stacks, 'staging');
    expect(result).toEqual([
      { distribution: distributions.dist1, distributionId: 'physical1', name: 'dist1' },
      { distribution: distributions.dist2, distributionId: 'physical2', name: 'dist2' }
    ]);
  });

  it("throws an error when a distribution doesn't have a matching stack resource", () => {
    const distributions = {
      dist3: {
        paths: ['a', 'b', 'c']
      }
    };
    expect(() => matchDistributions(distributions, stacks, 'staging')).toThrow();
  });

  it("throws an error when a distribution matches a stack resource which isn't a cloudfront distribution", () => {
    const distributions = {
      nondistribution: {
        paths: ['a', 'b', 'c']
      }
    };
    expect(() => matchDistributions(distributions, stacks, 'staging')).toThrow();
  });
});
