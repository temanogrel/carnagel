// @flow
import { combineEpics } from 'redux-observable';
import { createSelector } from 'reselect';
import { Observable } from 'rxjs/Observable';

// ------------------------------------
// Constants
// ------------------------------------
export const MODAL_UNSHIFT = 'modals:unshift';
export const MODAL_PUSH = 'modals:push';
export const MODAL_POP = 'modals:pop';
export const MODAL_REPLACE = 'modals:replace';

export const MODAL_LOGIN = 'modals:login';
export const MODAL_ACCOUNT = 'modals:account';
export const MODAL_REGISTRATION = 'modals:registration';
export const MODAL_FORGOT_PASSWORD = 'modals:forgot-password';
export const MODAL_SET_NEW_PASSWORD = 'modals:set-new-password';
export const MODAL_PURCHASE_PAYMENT_PLAN = 'modals:purchase-transaction-plan';
export const MODAL_BANDWIDTH_EXPLAINED = 'modals:bandwidth:explained';
export const MODAL_SESSION_EXPIRED = 'modals:session-expired';
export const MODAL_OVERRIDE_EXISTING_SESSION = 'modals:override-existing-session';
export const MODAL_SESSION_BLACKLISTED_TODAY = 'modals:session-blacklisted-today';

// ------------------------------------
// Type definitions
// ------------------------------------
export type Modal = {
  type: string;
  payload: Object;
  onClose: Object;
};

// ------------------------------------
// Reducer
// ------------------------------------
const initialState = [];

export default function reducer(state: Modal[] = initialState, action) {
  switch (action.type) {
    case MODAL_UNSHIFT: {
      const copy = state.slice();
      copy.unshift(action.payload);

      return copy;
    }

    case MODAL_PUSH: {
      const copy = state.slice();
      copy.push(action.payload);

      return copy;
    }

    case MODAL_POP: {
      const copy = state.slice();
      copy.pop();

      return copy;
    }

    default:
      return state;
  }
}

/* ******************************************************************
    Selectors
****************************************************************** */
const getTopModal = ({ modals }) => modals.length > 0 ? modals[modals.length - 1] : null;

export const getCurrentModal = createSelector([getTopModal], (modal) => modal);

/* ******************************************************************
    Action creators
****************************************************************** */
export const replaceModal = (type, payload, onClose) => ({
  type: MODAL_REPLACE,
  payload: { type, payload, onClose },
});

export const unshiftModal = (type, payload, onClose) => ({
  type: MODAL_UNSHIFT,
  payload: { type, payload, onClose },
});

export const pushModal = (type, payload, onClose) => ({
  type: MODAL_PUSH,
  payload: { type, payload, onClose },
});

export const popModal = () => ({ type: MODAL_POP });

/** *****************************************************************
 Epics
 ****************************************************************** */
const replaceModalEpic = action$ =>
  action$
    .ofType(MODAL_REPLACE)
    .flatMap(({ payload }) =>
      Observable.concat(
        Observable.of(popModal()),
        Observable.of(pushModal(payload.type, payload.payload, payload.size, payload.onClose)),
      ),
    );

export const modalsEpics = combineEpics(
  replaceModalEpic,
);
