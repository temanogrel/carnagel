import type { CollectionResultMeta } from 'api/stdlib';
import { connect } from 'react-redux';
import { withRouter } from 'react-router-dom';
import React, { Component } from 'react';
import { RecordingEntity } from 'api';
import { Paginator, RecordingPreview } from 'components';
import { seoUpdate } from 'store/ducks/seo';
import { AnimateUp } from 'utils/animations';
import { getDescriptionOfSortedRecordings, getKeywordsOfRecordings, getTitleForPathname } from 'utils/seo';
import { Spinner } from 'components/spinner/spinner';
import { RecordingCollectionForm } from 'views/recording/forms/collection';

type Props = {
  search: {
    loading: boolean;
    result: RecordingEntity[];
    meta: CollectionResultMeta;
  };

  sorting: string;
  seo: Function;
};

type State = {
  animateUp: Object;
};

const mapStateToProps = (state) => ({
  sorting: state.recordings.sorting,
  search: state.recordings.search,
});

const mapDispatchToProps = (dispatch) => ({
  seo: (...args) => dispatch(seoUpdate(...args)),
});

@withRouter
@connect(mapStateToProps, mapDispatchToProps)
export default class extends Component<Props, State> {
  constructor(props) {
    super(props);

    this.state = { animateUp: null };
  }

  componentDidUpdate(prevProps) {
    const { search, sorting } = this.props;
    const { meta, result } = search;

    if (!prevProps.search.loading || search.loading) {
      return;
    }

    setTimeout(() => {
      this.setState({ animateUp: new AnimateUp() });
    }, 180);

    this.props.seo(
      getTitleForPathname(this.props.location.pathname),
      getDescriptionOfSortedRecordings(Math.ceil(meta.offset / 90) + 1, (Math.ceil(meta.total / 90)), sorting),
      getKeywordsOfRecordings(result),
    );
  }

  componentWillUnmount() {
    if (this.state.animateUp) {
      this.state.animateUp.destroy();
    }
  }

  render() {
    const { search } = this.props;

    return (
      <div className='container'>
        <RecordingCollectionForm/>
        <Spinner loading={ search.loading }/>
        { search.result.map(r => <RecordingPreview key={ r.uuid } recording={ r }/>) }
        <Paginator { ...search.meta }/>
      </div>
    );
  }
}
