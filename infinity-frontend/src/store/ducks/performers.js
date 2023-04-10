// @flow
import { PerformerService, RecordingService } from 'api';
import { combineEpics } from 'redux-observable';
import { Observable } from 'rxjs/Observable';
import { push } from 'react-router-redux';
import { FETCH_RECORDING_SUCCESS } from 'store/ducks/recordings';

// ------------------------------------
// Constants
// ------------------------------------
export const FETCH_PERFORMER = 'performers:fetch';
export const FETCH_PERFORMER_SUCCESS = 'performers:fetch-success';

export const FETCH_PERFORMERS = 'performers:fetch-collection';
export const FETCH_PERFORMERS_SUCCESS = 'performers:fetch-collection-success';
export const FETCH_PERFORMERS_ERROR = 'performers:fetch-collection-error';
export const FETCH_PERFORMERS_RESET = 'performers:fetch-collection-reset';

export const FETCH_PERFORMER_RECORDINGS = 'performers:recordings-fetch';
export const FETCH_PERFORMER_RECORDINGS_SUCCESS = 'performers:recordings-fetch-success';
export const FETCH_PERFORMER_RECORDINGS_ERROR = 'performers:recordings-fetch-error';

export const SET_CURRENT_PERFORMER = 'performers:set-current';

// ------------------------------------
// ------------------------------------
// Reducer
// ------------------------------------

const INITIAL_STATE = {
  collection: {},
  current: null,

  search: {
    recordings: {
      loading: false,
      meta: {},
      result: [],
    },
    performers: {
      header: {
        loading: false,
        meta: {},
        result: [],
      },
      embedded: {
        loading: false,
        meta: {},
        result: [],
      },
    },
  },
};

export default (state = INITIAL_STATE, action) => {
  switch (action.type) {
    case SET_CURRENT_PERFORMER:
      return { ...state, current: action.payload };

    case FETCH_RECORDING_SUCCESS: {
      const recordings = { ...state.search.recordings };
      const index = recordings.result.map(r => r.uuid).indexOf(action.payload.uuid);

      if (index > -1) {
        recordings.result[index] = action.payload;
      }

      return { ...state, search: { ...state.search, recordings } };
    }

    case FETCH_PERFORMER_SUCCESS: {
      const collection = { ...state.collection };
      collection[action.payload.uuid] = action.payload;
      collection[action.payload.slug] = action.payload;

      return { ...state, collection };
    }

    case FETCH_PERFORMERS: {
      const performers = { ...state.search.performers };
      performers[action.payload.id] = { loading: true, meta: {}, result: [] };

      return { ...state, search: { ...state.search, performers } };
    }

    case FETCH_PERFORMERS_SUCCESS: {
      const performers = { ...state.search.performers };
      performers[action.payload.id] = { loading: false, meta: action.payload.meta, result: action.payload.items };

      return { ...state, search: { ...state.search, performers } };
    }

    case FETCH_PERFORMERS_RESET:
    case FETCH_PERFORMERS_ERROR: {
      const performers = { ...state.search.performers };
      performers[action.payload.id] = { loading: false, meta: {}, result: [] };

      return { ...state, search: { ...state.search, performers } };
    }

    case FETCH_PERFORMER_RECORDINGS:
      return { ...state, search: { ...state.search, recordings: { loading: true, meta: {}, result: [] } } };

    case FETCH_PERFORMER_RECORDINGS_SUCCESS:
      return {
        ...state,
        search: {
          ...state.search,
          recordings: { loading: false, meta: action.payload.meta, result: action.payload.items },
        },
      };

    case FETCH_PERFORMER_RECORDINGS_ERROR:
      return { ...state, search: { ...state.search, recordings: { loading: false, meta: {}, result: [] } } };

    default:
      return state;
  }
}

/** *****************************************************************
 Selectors
 ****************************************************************** */
export const getCurrentPerformer = state => {
  if (!state.performers.current) {
    return null;
  }

  const performer = state.performers.collection[state.performers.current];

  // Return it like this so we can check values against null
  return performer ? performer : null;
};

/** *****************************************************************
 Action creators
 ****************************************************************** */
export const setCurrentPerformer = id => ({ type: SET_CURRENT_PERFORMER, payload: id });
export const fetchPerformers = (id, params) => ({ type: FETCH_PERFORMERS, payload: { id, params } });
export const fetchPerformersReset = id => ({ type: FETCH_PERFORMERS_RESET, payload: { id } });

/** *****************************************************************
 Epics
 ****************************************************************** */
const fetchPerformersEpic = (action$, { getState }) =>
  action$
    .ofType(FETCH_PERFORMERS)
    .debounceTime(300)
    .switchMap(({ payload }) =>
      Observable
        .fromPromise(PerformerService.search(payload.params))
        .switchMap(result => {
          const { performers } = getState();

          // Check if the form was reset after api call
          if (!performers.search.performers[payload.id].loading) {
            return Observable.empty();
          }

          return Observable.of({ type: FETCH_PERFORMERS_SUCCESS, payload: { ...result, id: payload.id } });
        })
        .catch(() => Observable.of({ type: FETCH_PERFORMERS_ERROR, payload: { id: payload.id } })),
    );

const fetchPerformerRecordingsEpic = action$ =>
  action$
    .ofType(FETCH_PERFORMER_RECORDINGS)
    .switchMap(({ payload }) =>
      Observable
        .fromPromise(RecordingService.getByPerformer(payload.id, payload.params))
        .map(payload => ({ type: FETCH_PERFORMER_RECORDINGS_SUCCESS, payload }))
        .catch(() => Observable.of({ type: FETCH_PERFORMER_RECORDINGS_ERROR })),
    );

const fetchPerformerEpic = (action$, { getState }) =>
  action$
    .ofType(SET_CURRENT_PERFORMER, FETCH_PERFORMER)
    .flatMap(({ payload }) => {
      const { performers } = getState();

      if (performers.collection[payload]) {
        return Observable.empty();
      }

      return Observable
        .fromPromise(PerformerService.getByUuidOrSlug(payload))
        .map(performer => ({ type: FETCH_PERFORMER_SUCCESS, payload: performer }))
        .catch(() => Observable.of(push('/404')));
    });

export const performersEpics = combineEpics(
  fetchPerformersEpic,
  fetchPerformerRecordingsEpic,
  fetchPerformerEpic,
);
