// @flow
import axios from 'axios';
import { apiUriInterceptor, handleUnauthorizedResponses } from './stdlib';

export const httpClient = axios.create({
  withCredentials: true,
});

export type ResultCollection = {
  items: Object[];
  meta: {
    limit: number;
    offset: number;
    total: number;
  };
}

// configure the interceptors
httpClient.interceptors.request.use(apiUriInterceptor);
httpClient.interceptors.response.use(null, response => handleUnauthorizedResponses(response));

export * from './websocket';
export * from './payment';
export * from './performer';
export * from './recording';
export * from './user';
export * from './session';
export * from './bandwidth';
export * from './proxy';
