// @flow
import { RecordingService } from 'api';
import { combineEpics } from 'redux-observable';
import { Observable } from 'rxjs/Observable';
import { cloneObject } from 'utils/object';
import { push } from 'react-router-redux';

// ------------------------------------
// Constants
// ------------------------------------
export const FETCH_RECORDINGS = 'recordings:fetch-collection';
export const FETCH_RECORDINGS_SUCCESS = 'recordings:fetch-collection-success';
export const FETCH_RECORDINGS_ERROR = 'recordings:fetch-collection-error';

export const FETCH_RECORDING = 'recordings:fetch';
export const FETCH_RECORDING_SUCCESS = 'recordings:fetch-success';

export const TOGGLE_RECORDING_LIKE = 'recordings:toggle-like';
export const TOGGLE_RECORDING_FAVORITE = 'recordings:toggle-favorite';
export const ADD_RECORDING_VIEW = 'recordings:view';

export const SET_CURRENT_RECORDING = 'recordings:set-current';

export const SET_RECORDINGS_SORTING = 'recordings:set-sorting';

export const RecordingSortMode = {
  POPULARITY: 'popularity',
  VIEWS: 'views',
  LATEST: undefined,
};

// ------------------------------------
// ------------------------------------
// Reducer
// ------------------------------------

const INITIAL_STATE = {
  collection: {},
  current: null,

  sorting: undefined,

  search: {
    result: [],
    meta: {},
    loading: false,
  },
};

export default (state = INITIAL_STATE, action) => {
  switch (action.type) {
    case SET_CURRENT_RECORDING:
      return { ...state, current: action.payload };

    case FETCH_RECORDING_SUCCESS: {
      const collection = { ...state.collection };
      collection[action.payload.uuid] = action.payload;
      collection[action.payload.slug] = action.payload;

      const search = { ...state.search };
      const index = search.result.map(r => r.uuid).indexOf(action.payload.uuid);

      if (index > -1) {
        search.result[index] = action.payload;
      }

      return { ...state, search, collection };
    }

    case SET_RECORDINGS_SORTING:
      return { ...state, sorting: action.payload };

    case FETCH_RECORDINGS:
      return { ...state, search: { loading: true, result: [], meta: {} } };

    case FETCH_RECORDINGS_ERROR:
      return { ...state, search: { loading: false, result: [], meta: {} } };

    case FETCH_RECORDINGS_SUCCESS:
      return { ...state, search: { loading: false, result: action.payload.items, meta: action.payload.meta } };

    default:
      return state;
  }
}

/** *****************************************************************
 Selectors
 ****************************************************************** */
export const getCurrentRecording = state => {
  if (!state.recordings.current) {
    return null;
  }

  const recording = state.recordings.collection[state.recordings.current];

  // Return it like this so we can check values against null
  return recording ? recording : null;
};

/** *****************************************************************
 Action creators
 ****************************************************************** */
export const toggleRecordingLike = recording => ({ type: TOGGLE_RECORDING_LIKE, payload: recording });
export const toggleRecordingFavorite = recording => ({ type: TOGGLE_RECORDING_FAVORITE, payload: recording });
export const addRecordingView = recording => ({ type: ADD_RECORDING_VIEW, payload: recording });
export const setCurrentRecording = id => ({ type: SET_CURRENT_RECORDING, payload: id });
export const setRecordingsSorting = sorting => ({ type: SET_RECORDINGS_SORTING, payload: sorting });

/** *****************************************************************
 Epics
 ****************************************************************** */
const fetchRecordingsEpic = action$ =>
  action$
    .ofType(FETCH_RECORDINGS)
    .switchMap(({ payload }) =>
      Observable
        .fromPromise(RecordingService.getAll(payload))
        .map(payload => ({ type: FETCH_RECORDINGS_SUCCESS, payload }))
        .catch(() => Observable.of({ type: FETCH_RECORDINGS_ERROR })),
    );

const toggleRecordingLikeEpic = action$ =>
  action$
    .ofType(TOGGLE_RECORDING_LIKE)
    .flatMap(({ payload }) =>
      Observable
        .fromPromise(RecordingService.toggleLike(payload.uuid))
        .map(isLiked => {
          const recording = cloneObject(payload);
          recording.toggleLike(isLiked);

          return ({ type: FETCH_RECORDING_SUCCESS, payload: recording });
        })
        .catch(() => Observable.empty()),
    );

const toggleRecordingFavoriteEpic = action$ =>
  action$
    .ofType(TOGGLE_RECORDING_FAVORITE)
    .flatMap(({ payload }) =>
      Observable
        .fromPromise(RecordingService.toggleFavorite(payload.uuid))
        .map(isFavorite => {
          const recording = cloneObject(payload);
          recording.toggleFavorite(isFavorite);

          return ({ type: FETCH_RECORDING_SUCCESS, payload: recording });
        })
        .catch(() => Observable.empty()),
    );

const addRecordingViewEpic = action$ =>
  action$
    .ofType(ADD_RECORDING_VIEW)
    .flatMap(({ payload }) =>
      Observable
        .fromPromise(RecordingService.addView(payload.uuid))
        .flatMap(added => {
          if (!added) {
            return Observable.empty();
          }

          const recording = cloneObject(payload);
          recording.addView();

          return Observable.of({ type: FETCH_RECORDING_SUCCESS, payload: recording });
        })
        .catch(() => Observable.empty()),
    );

const fetchRecordingEpic = (action$, { getState }) =>
  action$
    .ofType(SET_CURRENT_RECORDING, FETCH_RECORDING)
    .flatMap(({ payload }) => {
      const { recordings } = getState();

      if (recordings.collection[payload]) {
        return Observable.empty();
      }

      return Observable
        .fromPromise(RecordingService.getByUuidOrSlug(payload))
        .map(recording => ({ type: FETCH_RECORDING_SUCCESS, payload: recording }))
        .catch(() => Observable.of(push('/404')));
    });

export const recordingsEpics = combineEpics(
  fetchRecordingsEpic,
  fetchRecordingEpic,
  toggleRecordingLikeEpic,
  toggleRecordingFavoriteEpic,
  addRecordingViewEpic,
);
