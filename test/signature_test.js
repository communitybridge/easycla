const moment = require('moment');
const should = require('should');

const sig = require('../lib/signature');

const mockKey = {
  keyId: 'test_key',
  secret: 'test_secret'
};

suite('Signature', function () {

  suite('Version One', function () {

    test('signing requests without a body', function () {
      let req = sig.signRequestVersionOne(mockKey, {
        method: 'GET',
        uri: 'http://localhost:5000/projects'
      });

      should(req).be.an.Object();

      should(req.uri).equal('http://localhost:5000/projects');
      should(req.method).equal('GET');

      should(req.headers).be.an.Object();
      should(req.headers['Date']).match(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z$/);
      should(req.headers['Authorization']).match(/^CINCO test_key:(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$/);
      should(req.headers['Date']).match(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z$/);
      should(req.headers['Content-Type']).equal('application/json; charset=UTF-8');
      should(req.headers['Signature-Version']).equal('1');
      should(req.headers).not.have.property('Content-MD5');
    });

    test('signing requests including a body', function () {
      let req = sig.signRequestVersionOne(mockKey, {
        method: 'GET',
        uri: 'http://localhost:5000/projects',
        body: JSON.stringify({
          question: 'What is love?',
          response: 'Baby, don\'t hurt me'
        })
      });

      should(req).be.an.Object();

      should(req.uri).equal('http://localhost:5000/projects');
      should(req.method).equal('GET');

      should(req.headers).be.an.Object();
      should(req.headers['Authorization']).match(/^CINCO test_key:(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$/);
      should(req.headers['Date']).match(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z$/);
      should(req.headers['Content-MD5']).equal('4ca47577b0a914d6f8c4649d4419bd43');
      should(req.headers['Content-Type']).equal('application/json; charset=UTF-8');
      should(req.headers['Signature-Version']).equal('1');
    });

  });

  suite('Version Two', function () {

    test('signing requests without a body', function () {
      let req = sig.signRequestVersionTwo(mockKey, {
        method: 'GET',
        uri: 'http://localhost:5000/projects'
      });

      should(req).be.an.Object();

      should(req.uri).equal('http://localhost:5000/projects');
      should(req.method).equal('GET');

      should(req.headers['Content-Type']).equal('application/json; charset=UTF-8');
      should(req.headers['Host']).equal('localhost:5000');
      should(req.headers['Signature-Version']).equal('2');
      should(req.headers['X-Amz-Date']).match(/^\d{8}T\d{6}Z$/);

      let headerDate = req.headers['X-Amz-Date'].slice(0, 8);

      should(req.headers['Authorization'])
        .startWith('AWS4-HMAC-SHA256 Credential=test_key/' + headerDate + '/internal/cinco/aws4_request, SignedHeaders=content-type;host;signature-version;x-amz-date, ');

      let authSigned = req.headers['Authorization'].slice(137);

      should(authSigned).match(/^Signature=[A-Fa-f0-9]{64}$/);
    });

    test('signing requests including a body', function () {
      let req = sig.signRequestVersionTwo(mockKey, {
        method: 'GET',
        uri: 'http://localhost:5000/projects',
        body: JSON.stringify({
          comments: 'You can dance if you want to'
        })
      });

      should(req).be.an.Object();

      should(req.uri).equal('http://localhost:5000/projects');
      should(req.method).equal('GET');

      should(req.headers['Content-Length']).equal(43);
      should(req.headers['Content-Type']).equal('application/json; charset=UTF-8');
      should(req.headers['Host']).equal('localhost:5000');
      should(req.headers['Signature-Version']).equal('2');
      should(req.headers['X-Amz-Date']).match(/^\d{8}T\d{6}Z$/);

      let headerDate = req.headers['X-Amz-Date'].slice(0, 8);

      should(req.headers['Authorization'])
        .startWith('AWS4-HMAC-SHA256 Credential=test_key/' + headerDate + '/internal/cinco/aws4_request, SignedHeaders=content-length;content-md5;content-type;host;signature-version;x-amz-date, ');

      let authSigned = req.headers['Authorization'].slice(164);

      should(authSigned).match(/^Signature=[A-Fa-f0-9]{64}$/);
    });

  });

});