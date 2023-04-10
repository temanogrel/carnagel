import { RecordingEntity } from 'api/recording';
import type { CollectionResultMeta } from 'api/stdlib';
import { connect } from 'react-redux';
import React, { Component } from 'react';
import { withRouter } from 'react-router-dom';
import { Paginator, RecordingPreview } from 'components';
import { AnimateUp } from 'utils/animations';
import { Spinner } from 'components/spinner/spinner';
import { UserFavoriteCollectionForm } from 'views/user/forms/favorite-collection';

type Props = {
  location: Object;
  search: {
    result: RecordingEntity[];
    meta: CollectionResultMeta;
    loading: boolean;
  }
};

type State = {
  animateUp: Object;
};

const mapStateToProps = state => ({
  search: state.userFavorites,
});

@withRouter
@connect(mapStateToProps)
export default class extends Component<Props, State> {
  constructor(props, context) {
    super(props, context);

    this.state = { animateUp: null };
  }

  componentDidUpdate(prevProps) {
    const { search } = this.props;

    if (!prevProps.search.loading || search.loading) {
      return;
    }

    setTimeout(() => {
      this.setState({ animateUp: new AnimateUp() });
    }, 180);
  }

  render() {
    const { result, loading, meta } = this.props.search;

    return (
      <div className="container">
        <Spinner loading={ loading }/>
        <h2 className="pt15 text-center thin">Your collection</h2>
        { !loading && <h3 className="text-center black mb20">{ meta.total } videos</h3> }
        { loading && <h3 className="text-center black mb20">Loading videos</h3> }
        <UserFavoriteCollectionForm initialValues={ { filter: '' } }/>
        { result.length === 0 && !loading ? <div className="text-center">No results found</div> : null }
        { result.map(r => <RecordingPreview key={ r.uuid } recording={ r }/>) }

        <Paginator { ...meta }/>
      </div>
    );
  }
}
