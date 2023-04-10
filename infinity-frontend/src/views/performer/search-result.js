import type { CollectionResultMeta } from 'api/stdlib';
import queryString from 'query-string';
import React, { Component } from 'react';
import LazyLoad from 'react-lazyload';
import { PerformerEntity, ProxyService } from 'api';
import { AnimateUp } from 'utils/animations';
import { withRouter, Link } from 'react-router-dom';
import { Paginator } from 'components/paginator/paginator';
import { RecordingImageEntity } from 'api/recording';
import { connect } from 'react-redux';
import { seoUpdate } from 'store/ducks/seo';
import { getKeywordsOfPerformers } from 'utils/seo';
import { Spinner } from 'components/spinner/spinner';
import { PerformerCollectionForm } from 'views/performer/forms/collection';

type Props = {
  location: { query: Object };
  seo: Function;
  search: {
    result: PerformerEntity[];
    meta: CollectionResultMeta;
    loading: boolean;
  };
};

type State = {
  animateUp: AnimateUp;
};

const mapStateToProps = state => ({
  search: state.performers.search.performers.embedded,
});

const mapDispatchToProps = (dispatch) => ({
  seo: (...args) => dispatch(seoUpdate(...args)),
});

@withRouter
@connect(mapStateToProps, mapDispatchToProps)
export default class extends Component<Props, State> {
  constructor(props: Props) {
    super(props);

    this.state = { animateUp: null };
  }

  componentDidUpdate(prevProps) {
    if (!prevProps.search.loading || this.props.search.loading) {
      return;
    }

    setTimeout(() => {
      this.setState({ animateUp: new AnimateUp() });
    }, 180);

    const params = queryString.parse(location.search);
    const page = parseInt(params.page) ? parseInt(params.page) : 1;

    this.props.seo(
      `Search results for ${params.query} - camtube.co`,
      `Search results for "${params.query}", page ${page} out of ${Math.ceil(this.props.search.meta.total / 30)}`,
      getKeywordsOfPerformers(this.props.search.result),
    );
  }

  render() {
    const params = queryString.parse(location.search);

    return (
      <div className="performer-search mt30">
        <h3 className="text-center mb20">Search results for "{ params.query }"</h3>
        <div className="container mt20">
          <PerformerCollectionForm/>
          <Spinner loading={ this.props.search.loading }/>
          { this.renderSearchResult() }
          <Paginator { ...this.props.search.meta } additionalQuery={ `&query=${params.query}` }/>
        </div>
      </div>
    );
  }

  renderSearchResult() {
    if (this.props.search.result.length === 0 && !this.props.search.loading) {
      return <div className="text-center">No results found</div>;
    }

    return this.props.search.result.map((p: PerformerEntity) => {
      let imageUrl = null;
      if (p.latestRecording) {
        imageUrl = `url(${ProxyService.getImageUrl(new RecordingImageEntity(p.latestRecording.collageUuid))})`;
      }

      return (
        <LazyLoad once>
          <div
            className="animate-up-opacity-base animate-up-opacity grid-4 tablet-grid-6 mobile-grid-12"
            key={ p.uuid }
          >
            <div className="image" style={ { backgroundImage: imageUrl, backgroundSize: 'contain' } }>
              <Link className="image-shade" to={ `/performer/${p.slug}` }>
                <span className="name">{ p.stageName }</span>
                <span className="other-videos">{ p.recordingCount }</span>
              </Link>
            </div>
          </div>
        </LazyLoad>
      );
    });
  }
}
