// @flow
import { PerformerService, RecordingService } from 'api';
import { combineEpics } from 'redux-observable';
import { Observable } from 'rxjs/Observable';
import { push } from 'react-router-redux';
import { FETCH_RECORDING_SUCCESS } from 'store/ducks/recordings';

// ------------------------------------
// Constants
// ------------------------------------
export const FETCH_USER_FAVORITES = 'user-favorites:fetch';
export const FETCH_USER_FAVORITES_SUCCESS = 'user-favorites:fetch-success';
export const FETCH_USER_FAVORITES_ERROR = 'user-favorites:fetch-error';

// ------------------------------------
// ------------------------------------
// Reducer
// ------------------------------------

const INITIAL_STATE = {
  loading: false,
  result: [],
  meta: {},
};

export default (state = INITIAL_STATE, action) => {
  switch (action.type) {
    case FETCH_USER_FAVORITES_ERROR:
      return { loading: false, result: [], meta: {} };

    case FETCH_USER_FAVORITES:
      return { loading: true, result: [], meta: {} };

    case FETCH_USER_FAVORITES_SUCCESS:
      return { loading: false, result: action.payload.items, meta: action.payload.meta };

    case FETCH_RECORDING_SUCCESS: {
      if (action.payload.isFavorite) {
        return state;
      }

      const result = state.result.filter(r => r.uuid !== action.payload.uuid);

      return { ...state, result };
    }

    default:
      return state;
  }
}

/** *****************************************************************
 Selectors
 ****************************************************************** */

/** *****************************************************************
 Action creators
 ****************************************************************** */
export const fetchUserFavorites = payload => ({ type: FETCH_USER_FAVORITES, payload });

/** *****************************************************************
 Epics
 ****************************************************************** */
const fetchUserFavoritesEpic = (action$, { getState }) =>
  action$
    .ofType(FETCH_USER_FAVORITES)
    .debounceTime(300)
    .switchMap(({ payload }) =>
      Observable
        .fromPromise(RecordingService.getUserFavorites(getState().user.uuid, payload))
        .map(payload => ({ type: FETCH_USER_FAVORITES_SUCCESS, payload }))
        .catch(() => Observable.of({ type: FETCH_USER_FAVORITES_ERROR })),
    );

export const userFavoritesEpics = combineEpics(
  fetchUserFavoritesEpic,
);
