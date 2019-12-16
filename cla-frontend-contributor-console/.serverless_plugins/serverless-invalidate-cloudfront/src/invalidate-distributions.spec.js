const invalidateDistributions = require('./invalidate-distributions');

describe('invalidateDistributions', () => {
  let aws, cli, cloudfront, promise, distributions;
  beforeEach(() => {
    const stacks = {
      StackResources: [
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
      ]
    };
    const stackName = 'stack';
    aws = {
      naming: jasmine.createSpyObj('naming', {
        getStackName: jasmine.createSpy('getStackName').and.returnValue(stackName)
      }),
      request: jasmine.createSpy('request').and.returnValue(Promise.resolve(stacks))
    };

    cloudfront = {
      createInvalidation: jasmine.createSpy('createInvalidation').and.returnValue({ promise: () => Promise.resolve() })
    };

    cli = jasmine.createSpyObj('cli', {
      log: jasmine.createSpy('log')
    });
    distributions = {
      dist1: {
        paths: ['a', 'b', 'c']
      }
    };
    promise = invalidateDistributions(aws, distributions, cloudfront, cli);
  });
  it('calls createInvalidation with the correct distribution id', (done) => {
    promise.then(() => {
      expect(cloudfront.createInvalidation).toHaveBeenCalledWith(
        jasmine.objectContaining({
          DistributionId: 'physical1'
        }),
        jasmine.any(Function)
      );
      done();
    });
  });
  it('calls createInvalidation with the correct paths', (done) => {
    promise.then(() => {
      expect(cloudfront.createInvalidation).toHaveBeenCalledWith(
        jasmine.objectContaining({
          InvalidationBatch: jasmine.objectContaining({
            Paths: {
              Quantity: 3,
              Items: distributions.dist1.paths
            }
          })
        }),
        jasmine.any(Function)
      );
      done();
    });
  });
  it('rejects distributions without paths', (done) => {
    distributions = {
      dist1: {}
    };
    promise = invalidateDistributions(aws, distributions, cloudfront, cli);
    promise.catch(() => {
      done();
    });
  });
});
