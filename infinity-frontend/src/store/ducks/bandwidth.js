// @flow
import { BandwidthService } from 'api';
import { combineEpics } from 'redux-observable';
import { Observable } from 'rxjs/Observable';
import { WEBSOCKET_SESSION_SET } from 'store/ducks/websocket';

// ------------------------------------
// Constants
// ------------------------------------
export const FETCH_BANDWIDTH = 'bandwidth:remaining';
export const FETCH_BANDWIDTH_SUCCESS = 'bandwidth:remaining-success';
export const SET_BANDWIDTH_POLL_INTERVAL = 'bandwidth:set-poll-interval';

export const BandwidthPollInterval = {
  DEFAULT: 5000,
  RECORDING_VIEW: 1000,
};

// ------------------------------------
// ------------------------------------
// Reducer
// ------------------------------------

const INITIAL_STATE = { total: 0, remaining: 0 };

export default (state: State = INITIAL_STATE, action) => {
  switch (action.type) {
    case FETCH_BANDWIDTH_SUCCESS:
      return action.payload;

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
export const setBandwidthPollInterval = interval => ({ type: SET_BANDWIDTH_POLL_INTERVAL, payload: interval });

/** *****************************************************************
 Epics
 ****************************************************************** */
const fetchBandwidthEpic = action$ =>
  action$
    .ofType(WEBSOCKET_SESSION_SET, FETCH_BANDWIDTH)
    .switchMap(() =>
      Observable
        .fromPromise(BandwidthService.getRemainingBandwidth())
        .map(payload => ({ type: FETCH_BANDWIDTH_SUCCESS, payload }))
        .catch(() => Observable.empty()),
    );

const pollBandwidthEpic = (action$, { getState }) =>
  action$
    .ofType(SET_BANDWIDTH_POLL_INTERVAL)
    .switchMap(({ payload }) =>
      Observable
        .timer(0, payload)
        .flatMap(() => {
          if (!getState().session) {
            return Observable.empty();
          }

          return Observable.of({ type: FETCH_BANDWIDTH });
        }),
    );

export const bandwidthEpics = combineEpics(
  pollBandwidthEpic,
  fetchBandwidthEpic,
);
