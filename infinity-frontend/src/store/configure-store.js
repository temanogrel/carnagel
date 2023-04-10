import { routerMiddleware, routerReducer } from 'react-router-redux';
import { applyMiddleware, compose, createStore, combineReducers } from 'redux';
import { combineEpics, createEpicMiddleware } from 'redux-observable';
import { default as thunkMiddleware } from 'redux-thunk';
import { reducer as formReducer } from 'redux-form';
import 'rxjs/Rx';
import history from 'utils/history';
import session, { sessionEpics } from './ducks/session';
import payments, { paymentsEpics} from './ducks/payments';
import websocket, { websocketEpics } from './ducks/websocket';
import bandwidth, { bandwidthEpics } from './ducks/bandwidth';
import seo from './ducks/seo';
import user, { userEpics } from './ducks/user';
import recordings, { recordingsEpics } from './ducks/recordings';
import toaster, { toasterEpics } from './ducks/toaster';
import modals, { modalsEpics } from 'store/ducks/modals';
import performers, { performersEpics } from 'store/ducks/performers';
import userFavorites, { userFavoritesEpics } from 'store/ducks/user-favorites';

const rootEpic = combineEpics(
  sessionEpics,
  websocketEpics,
  bandwidthEpics,
  paymentsEpics,
  userEpics,
  recordingsEpics,
  toasterEpics,
  modalsEpics,
  performersEpics,
  userFavoritesEpics,
);

const rootReducer = combineReducers({
  router: routerReducer,
  form: formReducer,
  session,
  websocket,
  performers,
  bandwidth,
  seo,
  payments,
  user,
  recordings,
  toaster,
  modals,
  userFavorites,
});

const middlewares = [
  routerMiddleware(history),
  thunkMiddleware,
  createEpicMiddleware(rootEpic),
];

if (__DEV__) {
  const createLogger = require('redux-logger'); // eslint-disable-line
  const logger = createLogger({
    collapsed: true,
    // Suppress a few redux-form actions
    predicate: (getState, action) =>
      !(
        action.type === '@@redux-form/CHANGE' ||
        action.type === '@@redux-form/FOCUS' ||
        action.type === '@@redux-form/BLUR' ||
        action.type === '@@redux-form/REGISTER_FIELD' ||
        action.type === '@@redux-form/UNREGISTER_FIELD' ||
        action.type === '@@redux-form/DESTROY'
      ),
  });

  middlewares.push(logger);
}

const initialState = {};

export const store = createStore(rootReducer, initialState, compose(applyMiddleware(...middlewares)));
