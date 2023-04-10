import { Spinner } from 'components/spinner/spinner';
import React, { Component } from 'react';
import { withRouter } from 'react-router-dom';
import { connect } from 'react-redux';
import { RecordingSortMode, setRecordingsSorting } from 'store/ducks/recordings';

type Props = {
  children: Object;
  sorting: string;
  setSorting: setRecordingsSorting;
};

const mapStateToProps = state => ({
  sorting: state.recordings.sorting,
});

const mapDispatchToProps = dispatch => ({
  setSorting: sorting => dispatch(setRecordingsSorting(sorting)),
});

@withRouter
@connect(mapStateToProps, mapDispatchToProps)
export class RecordingListSortingWrapper extends Component<Props> {
  componentDidMount() {
    this.updateSortingBasedOnRoute(this.props.location);
  }

  componentWillReceiveProps(nextProps) {
    if (this.props.location.pathname !== nextProps.location.pathname) {
      this.updateSortingBasedOnRoute(nextProps.location)
    }
  }

  updateSortingBasedOnRoute(location) {
    return this.props.setSorting(this.getRecordingSortModeForPathname(location.pathname));
  }

  getRecordingSortModeForPathname(): string {
    switch (location.pathname) {
      case '/':
        return RecordingSortMode.LATEST;

      case '/most-viewed':
        return RecordingSortMode.VIEWS;

      case '/most-popular':
        return RecordingSortMode.POPULARITY;
    }
  }

  render() {
    if (this.props.sorting !== this.getRecordingSortModeForPathname(this.props.location.pathname)) {
      return <div className='container'><Spinner loading/></div>;;
    }

    return this.props.children;
  }
}
