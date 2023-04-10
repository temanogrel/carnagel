import { PerformerEntity } from 'api/performer';
import type { CollectionResultMeta } from 'api/stdlib';
import { connect } from 'react-redux';
import React, { Component } from 'react';
import { Paginator, RecordingPreview } from 'components';
import { getCurrentPerformer } from 'store/ducks/performers';
import { AnimateUp } from 'utils/animations';
import { seoUpdate } from 'store/ducks/seo';
import { RecordingEntity } from 'api/recording';
import {
  getDescriptionForCollectionPage,
  getKeywordsOfPerformer,
} from 'utils/seo';
import { Spinner } from 'components/spinner/spinner';
import { PerformerRecordingCollectionForm } from 'views/performer/forms/recording-collection';

type Props = {
  performer: PerformerEntity;
  search: {
    recordings: RecordingEntity[];
    meta: CollectionResultMeta;
    loading: boolean;
  };

  seo: seoUpdate;
};

type State = {
  animateUp: AnimateUp;
};

const mapStateToProps = state => ({
  search: state.performers.search.recordings,
  performer: getCurrentPerformer(state),
});

const mapDispatchToProps = (dispatch) => ({
  seo: (...args: string) => dispatch(seoUpdate(...args)),
});

@connect(mapStateToProps, mapDispatchToProps)
export default class extends Component<Props, State> {
  constructor(props) {
    super(props);

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

    const page = (search.meta.offset / 30) + 1;

    this.props.seo(
      `Page ${page} - camtube.co`,
      getDescriptionForCollectionPage(page, Math.ceil(search.meta.total / 30), this.props.performer),
      getKeywordsOfPerformer(this.props.performer),
    );
  }

  render() {
    const { loading, meta, result } = this.props.search;

    return (
      <div className='container'>
        <Spinner loading={loading}/>
        <PerformerRecordingCollectionForm/>
        <h2 className='pt15 text-center thin'>{ this.props.performer.stageName }</h2>
        <h3 className='pb20 text-center black'>{ meta.total } Videos</h3>
        { result.map(r => <RecordingPreview key={ r.uuid } recording={ r }/>) }

        <Paginator { ...meta }/>
      </div>
    );
  }
}
