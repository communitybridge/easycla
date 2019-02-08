export const AUTH_ROUTE = '#/auth';

const HTTP_PORT = '80';
const HTTPS_PORT = '443';

export function getURL(route: string, location: Location) {
  const port = location.port;
  const hasStandardPort = port === HTTPS_PORT || port === HTTP_PORT || port === '';
  const redirectPort = hasStandardPort ? '' : `:${port}`;
  return `${location.protocol}//${location.hostname}${redirectPort}/${route}`;
}

export function getAuthURL(location: Location) {
  return getURL(AUTH_ROUTE, location);
}

// The URL to which Auth0 will redirect the browser after authorization has been granted for the user.
export function getAuthURLFromWindow() {
  return getAuthURL(window.location);
}
