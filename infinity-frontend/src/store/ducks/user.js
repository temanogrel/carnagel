// @flow
import { UserService } from 'api';
import { combineEpics } from 'redux-observable';
import { LOGOUT, SESSION_EXPIRED, SET_SESSION } from 'store/ducks/session';
import { Observable } from 'rxjs/Observable';

// ------------------------------------
// Constants
// ------------------------------------
export const RESOLVE_CURRENT_USER = 'user:resolve-current';
export const RESOLVE_CURRENT_USER_ERROR = 'user:resolve-current-error';
export const RESOLVE_CURRENT_USER_SUCCESS = 'user:resolve-current-success';

export const Roles = {
  GUEST: 'guest',
  USER: 'user',
  ADMIN: 'admin',
};

// ------------------------------------
// ------------------------------------
// Reducer
// ------------------------------------

const INITIAL_STATE = null;

export default (state: State = INITIAL_STATE, action) => {
  switch (action.type) {
    case SESSION_EXPIRED:
    case LOGOUT:
    case RESOLVE_CURRENT_USER_ERROR:
      return null;

    case RESOLVE_CURRENT_USER_SUCCESS:
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
export const resolveCurrentUser = () => ({ type: RESOLVE_CURRENT_USER });

/** *****************************************************************
 Epics
 ****************************************************************** */
const loadCurrentUserEpic = action$ =>
  action$
    .ofType(RESOLVE_CURRENT_USER)
    .debounceTime(50)
    .switchMap(() => Observable
      .fromPromise(UserService.getCurrentUser())
      .map(user => ({ type: RESOLVE_CURRENT_USER_SUCCESS, payload: user }))
      .catch(() => Observable.of({ type: RESOLVE_CURRENT_USER_ERROR })),
    );

const loadCurrentUserForNewSessionEpic = action$ =>
  action$
    .ofType(SET_SESSION)
    .flatMap(({ payload }) => {
      if (payload.role === 'guest') {
        return Observable.empty();
      }

      return Observable.of(resolveCurrentUser());
    });

export const userEpics = combineEpics(
  loadCurrentUserEpic,
  loadCurrentUserForNewSessionEpic,
);
