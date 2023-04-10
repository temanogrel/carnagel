import { PerformerEntity } from 'api/performer';
import { Spinner } from 'components/spinner/spinner';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { withRouter } from 'react-router-dom';
import { getCurrentPerformer, setCurrentPerformer } from 'store/ducks/performers';

type Props = {
  children: Object;

  setCurrentPerformer: setCurrentPerformer;
  performer: PerformerEntity;

  match: {
    params: {
      id: number,
    };
  };
};

const mapStateToProps = state => ({
  performer: getCurrentPerformer(state),
});

const mapDispatchToProps = dispatch => ({
  setCurrentPerformer: id => dispatch(setCurrentPerformer(id)),
});

@withRouter
@connect(mapStateToProps, mapDispatchToProps)
export default class extends Component<Props> {
  componentDidMount() {
    this.props.setCurrentPerformer(this.props.match.params.id);
  }

  componentDidUpdate(prevProps) {
    if (this.props.match.params.id !== prevProps.match.params.id) {
      this.props.setCurrentPerformer(this.props.match.params.id);
    }
  }

  render() {
    const { match, performer } = this.props;

    if (!performer || (performer.uuid !== match.params.id && performer.slug !== match.params.id)) {
      return <div className='container'><Spinner loading/></div>;
    }

    return this.props.children;
  }
}
