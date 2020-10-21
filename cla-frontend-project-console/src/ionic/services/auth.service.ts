import {Injectable} from '@angular/core';
import createAuth0Client from '@auth0/auth0-spa-js';
import Auth0Client from '@auth0/auth0-spa-js/dist/typings/Auth0Client';
import {BehaviorSubject, Observable, Subject} from 'rxjs';
import {catchError, concatMap, shareReplay, tap} from 'rxjs/operators';
import {from} from 'rxjs/observable/from';
import {of} from 'rxjs/observable/of';
import {combineLatest} from 'rxjs/observable/combineLatest';
import {reject} from 'lodash';
import {AUTH_ROUTE} from './auth.utils';
import {EnvConfig} from './cla.env.utils';

@Injectable()
export class AuthService {
  loading$ = new BehaviorSubject<any>(true);
  // Create subject and public observable of user profile data
  private userProfileSubject$ = new BehaviorSubject<any>(null);
  userProfile$ = this.userProfileSubject$.asObservable();
  // Create a local property for login status
  loggedIn = false;
  redirectRoot: Subject<any> = new Subject<any>();
  checkSession: Subject<any> = new Subject<any>();
  currentHref = window.location.href;

  auth0Options = {
    clientId: EnvConfig['auth0-clientId'],
    domain: EnvConfig['auth0-domain'],
    // redirectUri: `${window.location.origin}` //+ AUTH_ROUTE, // *info from allowed_logout_urls
  };

  // Create an observable of Auth0 instance of client
  auth0Client$ = (from(
    createAuth0Client({
      domain: this.auth0Options.domain,
      client_id: this.auth0Options.clientId,
      // Luis said turn off for now
      // redirect_uri: this.auth0Options.redirectUri,
      //cacheLocation: 'memory',
      //useRefreshTokens: true,
    })
  ) as Observable<Auth0Client>).pipe(
    shareReplay(1), // Every subscription receives the same shared value
    catchError((err) => {
      this.loading$.next(false);
      return reject(err);
    })
  );

  // Define observables for SDK methods that return promises by default
  // For each Auth0 SDK method, first ensure the client instance is ready
  // concatMap: Using the client instance, call SDK method; SDK returns a promise
  // from: Convert that resulting promise into an observable
  isAuthenticated$ = this.auth0Client$.pipe(
    concatMap((client: Auth0Client) => from(client.isAuthenticated())),
    tap((res: any) => {
      // *info: once isAuthenticated$ responses , SSO session is loaded
      this.loading$.next(false);
      this.loggedIn = res;
    })
  );

  handleRedirectCallback$ = this.auth0Client$.pipe(
    concatMap((client: Auth0Client) =>
      from(client.handleRedirectCallback(this.currentHref))
    ));


  constructor() {
    // On initial load, check authentication state with authorization server
    // Set up local auth streams if user is already authenticated
    this.localAuthSetup();
    // Handle redirect from Auth0 login
    this.handleAuthCallback();
  }

  // When calling, options can be passed if desired
  // https://auth0.github.io/auth0-spa-js/classes/auth0client.html#getuser
  getUser$(options?): Observable<any> {
    return this.auth0Client$.pipe(
      concatMap((client: Auth0Client) => from(client.getUser(options))),
      tap((user) => {
        this.setSession(user);
        this.userProfileSubject$.next(user);
      })
    );
  }

  private localAuthSetup() {
    // This should only be called on app initialization
    // Set up local authentication streams
    const checkAuth$ = this.isAuthenticated$.pipe(
      concatMap((loggedIn: boolean) => {
        console.log('localAuthSetup - loggedIn: ' + loggedIn);
        this.loggedIn = loggedIn;
        this.checkSession.next(loggedIn);
        if (loggedIn) {
          // If authenticated, get user and set in app
          // NOTE: you could pass options here if needed
          return this.getUser$();
        }
        // If not authenticated, return stream that emits 'false'
        return of(loggedIn);
      })
    );
    checkAuth$.subscribe();
  }

  login(redirectPath: string = '/') {
    // A desired redirect path can be passed to login method
    // (e.g., from a route guard)
    // Ensure Auth0 client instance exists
    this.auth0Client$.subscribe((client: Auth0Client) => {
      // Call method to log in
      client.loginWithRedirect({
        redirect_uri: `${window.location.origin}`,
        appState: { target: redirectPath },
      });
    });
  }

  private handleAuthCallback() {
    // Call when app reloads after user logs in with Auth0
    const params = window.location.search;
    if (params.includes('code=') && params.includes('state=')) {
      let targetRoute: string; // Path to redirect to after login processsed
      const authComplete$ = this.handleRedirectCallback$.pipe(
        // Have client, now call method to handle auth callback redirect
        tap((cbRes: any) => {
          const element: any = document.getElementById('lfx-header');
          if (element) {
            element.authuser = this.getUser$();
          }
          // Get and set target redirect route from callback results
          targetRoute = cbRes.appState && cbRes.appState.target ? cbRes.appState.target : '/';
        }),
        concatMap(() => {
          // Redirect callback complete; get user and login status
          return combineLatest([this.getUser$(), this.isAuthenticated$]);
        })
      );
      // Subscribe to authentication completion observable
      // Response will be an array of user and login status
      authComplete$.subscribe(() => {
        console.log('redirecting to: ' + targetRoute);
        // Redirect to target route after callback processing
        this.redirectRoot.next(targetRoute);
      });
    }
  }

  private setSession(authResult): void {
    localStorage.setItem('userid', authResult.nickname);
    localStorage.setItem('user_email', authResult.email);
    localStorage.setItem('user_name', authResult.name);
  }

  logout() {
    // Ensure Auth0 client instance exists
    this.auth0Client$.subscribe((client: Auth0Client) => {
      // Call method to log out
      let redirectUri = window.location.origin; // this.auth0Options.redirectUri;
      if (EnvConfig['lfx-header-enabled'] === "true") {
        redirectUri = EnvConfig['landing-page'];
      }
      client.logout({
        client_id: this.auth0Options.clientId,
        returnTo: redirectUri,
      });
    });
  }

  getTokenSilently$(options?): Observable<any> {
    return this.auth0Client$.pipe(concatMap((client: Auth0Client) => from(client.getTokenSilently(options))));
  }

  public getIdToken(): Promise<string> {
    return new Promise<string>((resolve, reject) => {
      const token = this.getIdToken$({ ignoreCache: true }).toPromise();
      resolve(token);
    });
  }

  getIdToken$(options?): Observable<any> {
    return this.auth0Client$.pipe(
      concatMap((client: Auth0Client) => from(client.getIdTokenClaims(options))),
      concatMap((claims: any) => of(claims.__raw))
    );
  }
}
