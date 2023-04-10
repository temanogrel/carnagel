// @flow
import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import { lazyload } from 'react-lazyload';
import { connect } from 'react-redux';
import { RecordingEntity } from 'api/recording';
import { toggleRecordingFavorite } from 'store/ducks/recordings';
import { displayToaster } from 'store/ducks/toaster';

type Props = {
  recording: RecordingEntity;
  isAuthenticated: boolean;
  popup: displayToaster;
  toggleFavorite: toggleRecordingFavorite;
}

const mapStateToProps = (state) => ({
  isAuthenticated: state.user !== null,
});

const mapDispatchToProps = (dispatch) => ({
  toggleFavorite: recording => dispatch(toggleRecordingFavorite(recording)),
  popup: (message: string) => dispatch(displayToaster(message)),
});

@lazyload({
  once: true,
})
@connect(mapStateToProps, mapDispatchToProps)
export class RecordingPreview extends Component<Props> {
  toggleFavorite = () => this.props.toggleFavorite(this.props.recording);

  render() {
    const { recording } = this.props;

    // Only do the render up if we are not running inside PhantomJS
    const animateUp = !/PhantomJS/.test(window.navigator.userAgent) ? 'animate-up-opacity-base animate-up-opacity ' : '';

    let style = {
      backgroundImage: recording.collageUrl,
      backgroundSize: 'cover',
      paddingBottom: '75%',
    };

    if (recording.isMinified) {
      style.backgroundSize = 'cover';
      style.paddingBottom = '75%';
    }

    return (
      <div className={ animateUp + 'grid-4 tablet-grid-6 mobile-grid-12' } key={ recording.uuid }>
        <div className="image" style={ style }>
          <Link to={ '/recording/' + recording.slug } className="video-play-overlay"/>
          { this.renderToggleFavorite() }
          <Link className="image-shade" to={ `/performer/${recording.performer.slug}` }>
            <span className="name">{ recording.stageName }</span>
            <span className="other-videos">{ recording.performer.recordingCount }</span>
          </Link>
        </div>
        <div className="video-info">
          <div className="grid-4"><span className="views">{ recording.viewCount }</span></div>
          <div className="grid-4"><span className="likes">{ recording.likeCount }</span></div>
          <div className="grid-4"><span className="duration">{ recording.durationAsString() }</span></div>
          <div className="grid-12"><span>{ recording.createdAt.format('MMMM DD, YYYY') }</span></div>
          <div className="clear"/>
        </div>
      </div>
    );
  }

  renderToggleFavorite() {
    // If we"re not signed in we cannot add something to our collection
    if (!this.props.isAuthenticated) {
      return null;
    }

    if (this.props.recording.isFavorite) {
      return (<a onClick={ this.toggleFavorite } className="collection-remove"><span/>Remove video</a>);
    } else {
      return (<a onClick={ this.toggleFavorite } className="collection-add"><span/>Add to collection</a>);
    }
  }
}
