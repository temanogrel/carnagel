import { PerformerEntity } from 'api/performer';
import { reduxForm } from 'redux-form';
import { withRouter } from 'react-router-dom';
import React, { Component } from 'react';
import queryString from 'query-string';
import { fetchPerformers } from 'store/ducks/performers';

const onSubmit = (values, dispatch, { location }) => {
  const params = queryString.parse(location.search);
  const page = parseInt(params.page) ? parseInt(params.page) : 1;

  dispatch(fetchPerformers('embedded', {
    query: params.query.replace(' ', '_'),
    limit: 30,
    offset: (page - 1) * 30,
    includeLatestRecording: 1,
  }));
};

type Props = {
  submit: Function;
  location: Object;
  performer: PerformerEntity;
};

class PerformerCollection extends Component<Props> {
  shouldComponentUpdate({ location }) {
    const { search } = this.props.location;

    return search !== location.search;
  }

  componentDidMount() {
    this.props.submit();
  }

  componentDidUpdate(prevProps) {
    if (this.props.location.search !== prevProps.location.search) {
      // Url query is different, submit
      this.props.submit();
    }
  }

  render() {
    return null;
  }
}

export const PerformerCollectionForm = withRouter(reduxForm({
  form: 'performer:collection',
  onSubmit,
})(PerformerCollection));
