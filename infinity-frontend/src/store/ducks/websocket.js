// @flow
import { websocketService } from 'api';
import { Observable } from 'rxjs/Observable';
import { combineEpics } from 'redux-observable';
import { SET_SESSION } from 'store/ducks/session';

// ------------------------------------
// Constants
// ------------------------------------
export const WEBSOCKET_SESSION_SET = 'websocket:session-set';
export const WEBSOCKET_CONNECTION_CHANGED = 'websocket:connection-changed';

// Broadcasts
export const SESSION_ENDED_BROADCAST = 'session:ended';
export const PURCHASE_UPDATE_BROADCAST = 'purchase:update';

// ------------------------------------
// Type definitions
type State = {
  connected: false,
};

// ------------------------------------
// Reducer
// ------------------------------------

const INITIAL_STATE = {
  connected: false,
};

export default (state: State = INITIAL_STATE, action) => {
  switch (action.type) {
    case WEBSOCKET_CONNECTION_CHANGED:
      return { connected: action.payload };

    default:
      return state;
  }
};

/** *****************************************************************
 Selectors
 ****************************************************************** */

/** *****************************************************************
 Action creators
 ****************************************************************** */

export const websocketConnectionChanged = connected => ({ type: WEBSOCKET_CONNECTION_CHANGED, payload: connected });
export const websocketBroadcast = (broadcast, payload) => ({ type: broadcast, payload });

/** *****************************************************************
 Epics
 ****************************************************************** */
const setSessionOnReconnectEpic = (action$, { getState }) =>
  action$
    .ofType(WEBSOCKET_CONNECTION_CHANGED)
    .switchMap(({ payload }) => {
      if (payload && getState().session) {
        return Observable
          .fromPromise(websocketService.sendRpc('session:set', { token: getState().session.rawToken }))
          .mapTo(({ type: WEBSOCKET_SESSION_SET }))
          .catch(() => Observable.empty());
      }

      return Observable.empty();
    });

const websocketSetSessionEpic = action$ =>
  action$
    .ofType(SET_SESSION)
    .switchMap(({ payload }) =>
      Observable
        .fromPromise(websocketService.sendRpc('session:set', { token: payload.rawToken }))
        .mapTo(({ type: WEBSOCKET_SESSION_SET }))
        .catch(() => Observable.empty()),
    );

export const websocketEpics = combineEpics(
  setSessionOnReconnectEpic,
  websocketSetSessionEpic,
);
