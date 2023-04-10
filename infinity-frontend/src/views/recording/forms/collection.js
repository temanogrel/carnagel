import { reduxForm } from 'redux-form';
import { withRouter } from 'react-router-dom';
import { connect } from 'react-redux';
import React, { Component } from 'react';
import queryString from 'query-string';
import { FETCH_RECORDINGS } from 'store/ducks/recordings';

const onSubmit = (values, dispatch, { location, sorting }) => {
  const params = queryString.parse(location.search);
  const page = parseInt(params.page) ? parseInt(params.page) : 1;

  dispatch({
    type: FETCH_RECORDINGS, payload: {
      interval: 7,
      offset: (page - 1) * 90,
      limit: 90,
      sortMode: sorting ? sorting : undefined,
    },
  });
};

type Props = {
  submit: Function;
  location: Object;
};

class RecordingCollection extends Component<Props> {
  componentDidMount() {
    this.props.submit();
  }

  componentDidUpdate(prevProps) {
    const { key, search, pathname } = this.props.location;

    if (key !== prevProps.location.key && pathname === '/' && search === '') {
      // Pressing the logo should reload the newest videos
      this.props.submit();
    } else if (search !== prevProps.location.search) {
      // Url query is different, submit
      this.props.submit();
    }
  }

  render() {
    return null;
  }
}

const mapStateToProps = state => ({
  sorting: state.recordings.sorting,
});

export const RecordingCollectionForm = withRouter(connect(mapStateToProps, null)(reduxForm({
  form: 'recording:collection',
  onSubmit,
})(RecordingCollection)));
