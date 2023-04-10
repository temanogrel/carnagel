import React, { Component } from 'react';
import { RecordingEntity } from 'api/recording';
import { Link } from 'react-router-dom';
import { connect } from 'react-redux';
import { addRecordingView } from 'store/ducks/recordings';
import { isPrerender } from 'utils/prerender';

type Props = {
  bandwidthRemaining: number;
  recording: RecordingEntity;
  addView: addRecordingView;
};

const mapStateToProps = (state) => ({
  bandwidthRemaining: state.bandwidth.remaining,
});

const mapDispatchToProps = dispatch => ({
  addView: recording => dispatch(addRecordingView(recording)),
});

@connect(mapStateToProps, mapDispatchToProps)
export class VideoPlayer extends Component<Props> {
  inFullScreen = false;

  constructor(props) {
    super(props);

    this.state = {
      hasPlayed: false,
    };
  }

  shouldComponentUpdate(nextProps, nextState) {
    if (this.props.recording.uuid !== nextProps.recording.uuid) {
      return true;
    }

    return this.props.bandwidthRemaining !== nextProps.bandwidthRemaining;
  }

  componentDidMount() {
    this.player = new Clappr.Player({
      parent: this.container,
      source: this.props.recording.manifest,
      muted: false,
      width: '100%',
      height: '100%',
      poster: this.props.recording.collageRawUrl,
      hlsjsConfig: {
        enableWorker: true,
        xhrSetup: (xhr: XMLHttpRequest, url) => {
          xhr.withCredentials = true;

          // Abort request if out of bandwidth
          if (this.props.bandwidthRemaining <= 0) {
            xhr.destroy();

            return;
          }

          xhr.addEventListener('load', (event) => {
            if (event.target.status === 402 && this.inFullScreen) {
              this.player.core.mediaControl.toggleFullscreen();
            }
          });
        },
        initialLiveManifestSize: 3,
      },
      plugins: {
        core: [ClapprThumbnailsPlugin],
      },
      scrubThumbnails: {
        backdropHeight: 64,
        spotlightHeight: 84,
        thumbs: [],
      },
      events: {
        onPlay: () => {
          if (!this.state.hasPlayed) {
            this.setState({ hasPlayed: true }, () => {
              this.props.addView(this.props.recording);
            });
          }
        },
      },
    });

    this.thumbnailsPlugin = this.player.getPlugin('scrub-thumbnails');
    this.setThumbnails();

    const container = this.player.core.getCurrentContainer();
    container.on(Clappr.Events.CONTAINER_FULLSCREEN, () => this.inFullScreen = !this.inFullScreen);
  }

  /**
   * Generate an array of scrub thumbnails for the ClapprThumbnailsPlugin
   *
   * @returns {Array}
   */
  setThumbnails(): Array {
    if (isPrerender()) {
      // Prerender will never access the thumbnails so no need to load them
      return [];
    }

    let time = 0;

    this.props.recording.sprites.map((sprite) => {
      const url = sprite.rawUrl;

      let rows = 10;
      let columns = 10;
      let scrubHeight = 100;
      let scrubWidth = 178;

      if (this.props.recording.isMinified) {
        columns = 11;
        scrubHeight = 65;
        scrubWidth = 101;
        rows = Math.ceil(this.props.recording.duration / (11 * 5));
      }

      for (let row = 0; row < rows; row++) {
        for (let column = 0; column < columns; column++) {
          const thumb = {
            w: scrubWidth,
            h: scrubHeight,
            time: time,
            url: url,
            x: column * scrubWidth,
            y: row * scrubHeight,
          };

          time += 5;

          this
            .thumbnailsPlugin
            .addThumbnail(thumb)
            .catch(() => {
            });
        }
      }
    });
  }

  render() {
    return (
      <div className='video-player-container'>
        { this.props.bandwidthRemaining <= 0 && (
          <div className='out-of-bandwidth'>
            <div className='text-center text-red margin-auto mt10'>
              <h3>Out of bandwidth</h3>
              Wait until tomorrow to continue watching or upgrade your payment plan.
              <Link
                to='/payment-plans'
                onClick={ close }
                className='button red medium user-button mt10 mb10 margin-auto'>
                Upgrade plan
              </Link>
            </div>
          </div>
        ) }
        <div className='video-player' ref={ container => this.container = container }/>
      </div>
    );
  }
}
