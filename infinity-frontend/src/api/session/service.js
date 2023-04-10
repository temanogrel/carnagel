import { httpClient, SessionEntity } from 'api';
import { store } from 'store/configure-store';
import environment from 'api/environment';
import cookies from 'browser-cookies';
import { setSession } from 'store/ducks/session';

const COOKIE_SESSION_KEY = 'session';

export class SessionService {
  static initSession(): Promise<SessionEntity> {
    let inStorage = cookies.get(COOKIE_SESSION_KEY);
    if (inStorage === null) {
      return SessionService.createGuestSession();
    }

    let session;

    try {
      session = SessionEntity.createFromJwtToken(inStorage);
    } catch (e) {
      return SessionService.createGuestSession();
    }

    if (session.isExpired) {
      return SessionService.createGuestSession();
    }

    if (session.role !== 'guest') {
      return SessionService.renewAuthenticatedSession();
    }

    return SessionService.createGuestSession();
  }

  static createGuestSession(): Promise<SessionEntity> {
    return httpClient
      .get('backend://rpc/session/new')
      .then(({ data }) => {
        cookies.set(COOKIE_SESSION_KEY, data.token, {domain: environment.cookieDomain});

        const session = SessionEntity.createFromJwtToken(data.token);
        store.dispatch(setSession(session));

        return session;
      })
  }

  static createAuthenticatedSession(credentials): Promise<SessionEntity> {
    return httpClient
      .post('backend://rpc/session/authenticate', credentials)
      .then(({ data }) => {
        cookies.set(COOKIE_SESSION_KEY, data.token, {domain: environment.cookieDomain});

        const session = SessionEntity.createFromJwtToken(data.token);
        store.dispatch(setSession(session));

        return session;
      })
  }

  static renewAuthenticatedSession(): Promise<SessionEntity> {
    return httpClient
      .get('backend://rpc/session/renew')
      .then(({ data }) => {
        cookies.set(COOKIE_SESSION_KEY, data.token, {domain: environment.cookieDomain});

        const session = SessionEntity.createFromJwtToken(data.token);
        store.dispatch(setSession(session));

        return session;
      })
      .catch(() => this.createGuestSession())
  }
}
