import { PerformerEntity } from 'api/performer';
import { reduxForm } from 'redux-form';
import { withRouter } from 'react-router-dom';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import queryString from 'query-string';
import { FETCH_PERFORMER_RECORDINGS, getCurrentPerformer } from 'store/ducks/performers';

const onSubmit = (values, dispatch, { location, performer }) => {
  const params = queryString.parse(location.search);
  const page = parseInt(params.page) ? parseInt(params.page) : 1;

  dispatch({
    type: FETCH_PERFORMER_RECORDINGS, payload: {
      id: performer.uuid,
      params: {
        offset: (page - 1) * 30,
        limit: 30,
      },
    },
  });
};

type Props = {
  submit: Function;
  location: Object;
  performer: PerformerEntity;
};

class PerformerRecordingCollection extends Component<Props> {
  shouldComponentUpdate({ location }) {
    const { search, pathname } = this.props.location;

    return search !== location.search || pathname !== location.pathname;
  }

  componentDidMount() {
    this.props.submit();
  }

  componentDidUpdate(prevProps) {
    if (this.props.location.pathname !== prevProps.location.pathname) {
      this.props.submit();
    } else if (this.props.location.search !== prevProps.location.search) {
      // Url query is different, submit
      this.props.submit();
    }
  }

  render() {
    return null;
  }
}

const mapStateToProps = state => ({
  performer: getCurrentPerformer(state),
});

export const PerformerRecordingCollectionForm = withRouter(connect(mapStateToProps)(reduxForm({
  form: 'performer:recording-collection',
  onSubmit,
})(PerformerRecordingCollection)));
