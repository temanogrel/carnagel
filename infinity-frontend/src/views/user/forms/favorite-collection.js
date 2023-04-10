import { reduxForm, Field } from 'redux-form';
import { withRouter } from 'react-router-dom';
import React, { Component } from 'react';
import queryString from 'query-string';
import { fetchUserFavorites } from 'store/ducks/user-favorites';

const onSubmit = (values, dispatch, { location }) => {
  const params = queryString.parse(location.search);
  const page = parseInt(params.page) ? parseInt(params.page) : 1;

  dispatch(fetchUserFavorites({
    offset: (page - 1) * 60,
    limit: 60,
    stageName: values.filter.length > 0 ? values.filter.replace(' ', '_') : undefined,
  }));
};

type Props = {
  submit: Function;
  location: Object;
  filter: Object;
};

const renderFilter = ({ input }) => (
  <input
    { ...input }
    placeholder="Filter"
    type="text"
    style={ { maxWidth: '100%' } }
    className="search button red large transparent text-left margin-auto"/>
);

const onChange = (values, dispatch, { submit }) => submit();

class UserFavoriteCollection extends Component<Props> {
  componentDidMount() {
    this.props.submit();
  }

  componentDidUpdate(prevProps) {
    const { search } = this.props.location;

    if (search !== prevProps.location.search) {
      // Url query is different, submit
      this.props.submit();
    } else if (this.props.filter !== prevProps.filter) {
      this.props.submit();
    }
  }

  render() {
    return <Field name="filter" component={ renderFilter }/>;
  }
}

export const UserFavoriteCollectionForm = withRouter(reduxForm({
  form: 'user-favorites:collection',
  fields: ['filter'],
  onChange,
  onSubmit,
})(UserFavoriteCollection));
