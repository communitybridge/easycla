// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

function getHeaders(env, isDevServer) {
  return {
    'X-Content-Type-Options': 'nosniff',
    'X-Frame-Options': 'DENY',
    'Strict-Transport-Security': 'max-age=31536000; includeSubDomains',
    'X-XSS-Protection': '1',
    'Referrer-Policy': 'no-referrer',
    'Content-Security-Policy': generateCSP(env, isDevServer),
    'Cache-Control': 's-maxage=31536000'
  };
}

function getSources(environmentSources, sourceType) {
  if (environmentSources[sourceType] === undefined) {
    return [];
  }
  return environmentSources[sourceType].filter(source => {
    return typeof source === 'string';
  });
}

function generateCSP(env, isDevServer) {
  const SELF = "'self'";
  const UNSAFE_INLINE = "'unsafe-inline'";
  const UNSAFE_EVAL = "'unsafe-eval'";
  const NONE = "'none'";

  let connectSources = [SELF,
    'https://linuxfoundation-dev.auth0.com/',
    'https://linuxfoundation-staging.auth0.com/',
    'https://sso.linuxfoundation.org/',
    'https://api.staging.lfcla.com/',
    'https://api.dev.lfcla.com/',
    'https://api.test.lfcla.com/',
    'https://api.lfcla.com/',
    'https://communitybridge.org'
  ];
  let scriptSources = [SELF, UNSAFE_EVAL, UNSAFE_INLINE];
  let styleSources = [SELF, UNSAFE_INLINE, 'https://fonts.googleapis.com', 'https://communitybridge.org'];

  if (isDevServer) {
    connectSources = [...connectSources, 'https://localhost:8100/sockjs-node/', 'wss://localhost:8100/sockjs-node/'];
    // The webpack dev server uses system js which violates the unsafe-eval exception. This doesn't happen in the
    // production AOT build.
    scriptSources = [...scriptSources, UNSAFE_EVAL];
    // The development build needs unsafe inline assets.
  }

  const CSP_SOURCES = env ? env.CSP_SOURCES : undefined;
  const environmentSources = JSON.parse(CSP_SOURCES || '{}');

  const sources = {
    'default-src': [NONE],
    'img-src': [SELF, 'data:', 'https://s3.amazonaws.com/'],
    'script-src': scriptSources,
    'style-src': styleSources, // Unfortunately using Angular basically requires inline styles.
    'font-src': [SELF, 'https://fonts.googleapis.com', 'https://communitybridge.org'],
    'connect-src': connectSources,
    'frame-ancestors': [NONE],
    'form-action': [NONE],
    'worker-src': [SELF],
    'base-uri': [SELF],
    'frame-src': [],
    'child-src': [],
    'media-src': [],
    'manifest-src': [SELF],
  };

  return Object.entries(sources)
    .map(keyValuePair => {
      const additionalSources = getSources(environmentSources, keyValuePair[0]);
      return [keyValuePair[0], [...keyValuePair[1], ...additionalSources]];
    })
    .filter(keyValuePair => keyValuePair[1].length !== 0)
    .map(keyValuePair => {
      const entry = keyValuePair[1].join(' ');
      return `${keyValuePair[0]} ${entry};`;
    })
    .join(' ');
}

module.exports = getHeaders;
