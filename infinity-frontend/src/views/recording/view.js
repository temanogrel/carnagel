import { VideoPlayer } from 'components/recording/player';

// @flow
import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import { connect } from 'react-redux';
import { RecordingEntity } from 'api/recording';
import { getCurrentRecording, toggleRecordingFavorite, toggleRecordingLike } from 'store/ducks/recordings';
import { displayToaster } from 'store/ducks/toaster';
import { AnimateUp } from 'utils/animations';
import { seoUpdate } from 'store/ducks/seo';
import moment from 'moment';
import { originServiceToString } from 'utils/performer';
import { ImageCollection } from './image-collection';
import { getDescriptionOfRecording, getKeywordsOfRecording, getPostTitleOfRecording } from 'utils/seo';

type Props = {
  isAuthenticated: boolean;
  displayToaster: displayToaster;
  recording: RecordingEntity;
  toggleLike: toggleRecordingLike;
  toggleFavorite: toggleRecordingFavorite;
};

type State = {
  animateUp: Object;
};

const mapStateToProps = (state) => ({
  isAuthenticated: state.user !== null,
  recording: getCurrentRecording(state),
});

const mapDispatchToProps = (dispatch) => ({
  toggleLike: recording => dispatch(toggleRecordingLike(recording)),
  toggleFavorite: recording => dispatch(toggleRecordingFavorite(recording)),
  popup: (message: string) => dispatch(displayToaster(message)),
  seo: (...args) => dispatch(seoUpdate(...args)),
});

@connect(mapStateToProps, mapDispatchToProps)
export default class extends Component<Props, State> {
  constructor(props) {
    super(props);

    this.state = {
      animateUp: null,
    };
  }

  componentDidUpdate(prevProps) {
    if (this.props.recording.id === prevProps.recording.id) {
      return;
    }

    this.props.seo(
      `${getPostTitleOfRecording(this.props.recording)} - camtube.co`,
      getDescriptionOfRecording(this.props.recording),
      getKeywordsOfRecording(this.props.recording),
    );
  }

  componentWillUnmount() {
    if (this.state.animateUp) {
      this.state.animateUp.destroy();
    }
  }

  componentDidMount() {
    setTimeout(() => {
      this.setState({ animateUp: new AnimateUp() });
    }, 180);
  }

  toggleLike = () => {
    if (!this.props.isAuthenticated) {
      return;
    }

    this.props.toggleLike(this.props.recording);
  };

  toggleFavorite = () => {
    if (!this.props.isAuthenticated) {
      return;
    }

    this.props.toggleFavorite(this.props.recording);
  };

  render() {
    const { recording } = this.props;

    let likeButtonClasses = recording.isLiked ? 'blue' : 'red';
    likeButtonClasses = !this.props.isAuthenticated ? `${likeButtonClasses} disabled` : likeButtonClasses;

    return (
      <div className='container'>
        <div className='grid-12 recording'>
          <h3 className='text-center pb10'>
            { recording.stageName }&nbsp;
            { moment(recording.createdAt).format('DDMMYY HHMM') }&nbsp;
            { originServiceToString(recording.performer.originService) }&nbsp;
            { recording.performer.section }
          </h3>
          <h4 className='text-center pb20'>
            { moment(recording.createdAt).format('MMMM DD, YYYY') }
          </h4>
          <div className='image full-video'>
            <VideoPlayer recording={ recording }/>

            { this.renderToggleFavorite() }
          </div>
          <Link
            to={ `/performer/${recording.performer.slug}` }
            className='button medium square red grid-6 tablet-grid-12'
          >
            More from this performer ({ recording.performer.recordingCount })
          </Link>
          <a onClick={ this.toggleLike }
             className={ 'button medium square transparent grid-6 tablet-grid-12 ' + likeButtonClasses }>
            Like video
          </a>
          <div className='video-info clear'>
            <div className='grid-4'><span className='views'>{ recording.viewCount }</span></div>
            <div className='grid-4'><span className='likes'>{ recording.likeCount }</span></div>
            <div className='grid-4'><span className='duration'>{ recording.durationAsString() }</span></div>
            <div className='grid-12'><span>{ recording.createdAt.format('MMMM DD, YYYY') }</span></div>
            <div className='clear'/>
          </div>
        </div>

        { /*{ this.renderTagCollection() }*/ }
        <ImageCollection images={ recording.images }/>
      </div>
    );
  }

  renderToggleFavorite() {
    if (!this.props.isAuthenticated) {
      return null;
    }

    if (this.props.recording.isFavorite) {
      return (<a onClick={ this.toggleFavorite } className='collection-remove'><span/>Remove video</a>);
    } else {
      return (<a onClick={ this.toggleFavorite } className='collection-add'><span/>Add to collection</a>);
    }
  }

  renderTagCollection() {
    // Todo: implement this later
    return false;

    return (
      <div className='grid-12 pt0'>
        <div className='video-info clear tags-container'>
          <span className='tag'>Swag</span>
          <span className='tag'>Cool</span>

          <span className='tag add'>Add tag</span>
          <div className='clear'/>
          <form className='active'>
            <input type='text' className='white mr10 width-auto' placeholder='Tag name'/>
            <button className='red wide small user-button inline-block tablet-hidden mobile-hidden'>Add tag</button>
          </form>
        </div>
      </div>
    );
  }
}
