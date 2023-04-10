// @flow
import { SessionService } from 'api';
import { MODAL_SESSION_BLACKLISTED_TODAY, MODAL_SESSION_EXPIRED, pushModal, replaceModal } from 'store/ducks/modals';
import { Roles } from 'store/ducks/user';
import { SESSION_ENDED_BROADCAST } from 'store/ducks/websocket';
import { Observable } from 'rxjs/Observable';
import { combineEpics } from 'redux-observable';

// ------------------------------------
// Constants
// ------------------------------------
export const SESSION_EXPIRED = 'session:expired';
export const SET_SESSION = 'session:set';
export const LOGOUT = 'session:logout';

// ------------------------------------
// Reducer
// ------------------------------------

const INITIAL_STATE = null;

export default (state = INITIAL_STATE, action) => {
  switch (action.type) {
    case SET_SESSION:
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
export const logout = () => ({ type: LOGOUT });
export const setSession = payload => ({ type: SET_SESSION, payload  });

/** *****************************************************************
 Epics
 ****************************************************************** */
const checkIfAuthenticatedSessionIsBlackListedEpic = action$ =>
  action$
    .ofType(SET_SESSION)
    .switchMap(({ payload }) => {
      // Guest session cannot be blacklisted
      if (payload.role === Roles.GUEST) {
        return Observable.empty();
      }

      if (payload.blackListedToday) {
        return Observable.of(replaceModal(MODAL_SESSION_BLACKLISTED_TODAY));
      }

      return Observable.empty();
    });

const renewAuthenticatedSessionEpic = (action$, { getState }) =>
  action$
    .ofType(SET_SESSION)
    .switchMap(({ payload }) => {
      // Guest session cannot be blacklisted
      if (payload.role === Roles.GUEST) {
        return Observable.empty();
      }

      return Observable
        .timer(0, 1000 * 60)
        .switchMap(() => {
          const { session } = getState();

          // If session expires in less than five minutes, renew
          if (session.expiresIn < 5 * 60) {
            return Observable
              .fromPromise(SessionService.renewAuthenticatedSession())
              .flatMapTo(Observable.empty())
              .catch(() => Observable.of(logout()));
          }

          return Observable.empty();
        });
    });

const createGuestSessionOnLogoutEpic = action$ =>
  action$
    .ofType(LOGOUT)
    .switchMap(() =>
      Observable
        .fromPromise(SessionService.createGuestSession())
        .flatMapTo(Observable.empty()),
    );

const sessionEndedEpic = action$ =>
  action$
    .ofType(SESSION_ENDED_BROADCAST)
    .flatMapTo(Observable.concat(
      Observable.of(pushModal(MODAL_SESSION_EXPIRED)),
      Observable.of(logout()),
    ));

export const sessionEpics = combineEpics(
  renewAuthenticatedSessionEpic,
  checkIfAuthenticatedSessionIsBlackListedEpic,
  sessionEndedEpic,
  createGuestSessionOnLogoutEpic,
);
