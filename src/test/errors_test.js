const should = require('should');
const errors = require('../lib/errors');

suite('Errors', function () {

  test('creating errors from a response', function () {
    let bod1 = JSON.stringify({
      message: 'All right ramblers,',
      details: 'let\'s get rambling'
    });

    let err1 = errors.fromResponse({ statusCode: 401 });
    let err2 = errors.fromResponse({ statusCode: 401 }, 'Everyone and their mums is packin\' round here!');
    let err3 = errors.fromResponse({ statusCode: 401, body: '{}' });
    let err4 = errors.fromResponse({ statusCode: 401, body: bod1 });

    should(err1).be.an.Error();
    should(err1.statusCode).equal(401);
    should(err1.userMessage).equal('');
    should(err1.message).equal('Status: [401].');

    should(err2).be.an.Error();
    should(err2.statusCode).equal(401);
    should(err2.userMessage).equal('');
    should(err2.message).equal('Everyone and their mums is packin\' round here! Status: [401].');

    should(err3).be.an.Error();
    should(err3.statusCode).equal(401);
    should(err3.userMessage).equal('');
    should(err3.message).equal('Status: [401]. Response body: {}');

    should(err4).be.an.Error();
    should(err4.statusCode).equal(401);
    should(err4.userMessage).equal('All right ramblers, let\'s get rambling');
    should(err4.message).equal('Status: [401]. Response body: ' + bod1);
  });

});