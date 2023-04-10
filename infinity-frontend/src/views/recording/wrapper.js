import { RecordingEntity } from 'api/recording';
import { Spinner } from 'components/spinner/spinner';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { withRouter } from 'react-router-dom';
import { BandwidthPollInterval, setBandwidthPollInterval } from 'store/ducks/bandwidth';
import { getCurrentRecording, setCurrentRecording } from 'store/ducks/recordings';

type Props = {
  children: Object;

  setBandwidthPollInterval: setBandwidthPollInterval;
  setCurrentRecording: setCurrentRecording;
  recording: RecordingEntity;

  match: {
    params: {
      id: number,
    };
  };
};

const mapStateToProps = state => ({
  recording: getCurrentRecording(state),
});

const mapDispatchToProps = dispatch => ({
  setBandwidthPollInterval: interval => dispatch(setBandwidthPollInterval(interval)),
  setCurrentRecording: id => dispatch(setCurrentRecording(id)),
});

@withRouter
@connect(mapStateToProps, mapDispatchToProps)
export default class extends Component<Props> {
  componentDidMount() {
    // Set bandwidth poll interval for recording view
    this.props.setBandwidthPollInterval(BandwidthPollInterval.RECORDING_VIEW);
    this.props.setCurrentRecording(this.props.match.params.id);
  }

  componentWillUnmount() {
    // Reset bandwidth poll interval to default value
    this.props.setBandwidthPollInterval(BandwidthPollInterval.DEFAULT);
  }

  componentDidUpdate(prevProps) {
    if (this.props.match.params.id !== prevProps.match.params.id) {
      this.props.setCurrentRecording(this.props.match.params.id);
    }
  }

  render() {
    if (!this.props.recording) {
      return <div className='container'><Spinner loading/></div>;
    }

    return this.props.children;
  }
}
