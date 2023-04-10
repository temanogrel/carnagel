// @flow
import uuidv4 from 'uuid';
import { combineEpics } from 'redux-observable';
import { Observable } from 'rxjs/Observable';

// ------------------------------------
// Constants
// ------------------------------------
export const DISPLAY_TOASTER = 'toaster:display';
export const HIDE_TOASTER = 'toaster:hide';

// ------------------------------------
// Type definitions
// ------------------------------------
export type Toaster = {
  id: string;
  message: string;
};

// ------------------------------------
// Reducer
// ------------------------------------

const INITIAL_STATE = null;

export default (state: Toaster = INITIAL_STATE, action) => {
  switch (action.type) {
    case DISPLAY_TOASTER:
      return action.payload;

    case HIDE_TOASTER:
      return state?.id === action.payload ? null : state;

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
export const displayToaster = message => ({ type: DISPLAY_TOASTER, payload: { id: uuidv4(), message } });
export const hideToaster = id => ({ type: HIDE_TOASTER, payload: id });

/** *****************************************************************
 Epics
 ****************************************************************** */
const autoCloseEpic = action$ =>
  action$
    .ofType(DISPLAY_TOASTER)
    .delay(10000)
    .map(action => ({ type: HIDE_TOASTER, payload: action.payload.id }));

export const toasterEpics = combineEpics(
  autoCloseEpic,
);
