import { httpClient, SessionService } from 'api/index';
import { SESSION_EXPIRED } from 'store/ducks/session';
import environment from './environment';
import { store } from 'store/configure-store';

export const apiUriInterceptor = (request) => {
  if (request.url.substr(0, 10) === 'backend://') {
    request.url = environment.apiUri + request.url.substr(9);
  }

  return request;
};

export const handleUnauthorizedResponses = (error) => {
  if (!error.response) {
    return Promise.reject(error);
  }

  if (error.response.status === 401) {
    store.dispatch({ type: SESSION_EXPIRED });

    return SessionService
      .createGuestSession()
      .then(() => httpClient(error.response.config))
      .catch(() => Promise.reject(error));
  }

  return Promise.reject(error);
};

export type CollectionResultMeta = {
  offset: number;
  total: number;
  limit: number;
}

export type CollectionResult<T> = {
  items: T[];
  meta: CollectionResultMeta;
};
