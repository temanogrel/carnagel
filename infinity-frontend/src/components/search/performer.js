// @flow
import type { CollectionResultMeta } from 'api/stdlib';
import React, { Component } from 'react';
import Autocomplete from 'react-autocomplete';
import { connect } from 'react-redux';
import { push } from 'react-router-redux';
import { fetchPerformers, fetchPerformersReset } from 'store/ducks/performers';
import { isMobile } from 'utils/device';
import uuidv4 from 'uuid/v4';
import { PerformerEntity } from 'api';

type Props = {
  className: string;
  push: push;
  reduxId: string;
  reset: fetchPerformersReset;
  fetchPerformers: fetchPerformers;
  search: {
    result: PerformerEntity[];
    meta: CollectionResultMeta;
    loading: boolean;
  };
};

type State = {
  id: string;
  value: string;
  menuVisible: boolean,
  hoverActive: boolean;
};

const mapStateToProps = state => ({
  search: state.performers.search.performers,
});

const mapDispatchToProps = dispatch => ({
  reset: id => dispatch(fetchPerformersReset(id)),
  fetchPerformers: (id, params) => dispatch(fetchPerformers(id, params)),
  push: url => dispatch(push(url)),
});

@connect(mapStateToProps, mapDispatchToProps)
export class PerformerSearch extends Component<Props, State> {
  constructor(props) {
    super(props);

    this.state = {
      id: uuidv4(),
      value: '',
      menuVisible: false,
      hoverActive: false,
    };
  }

  componentDidUpdate(prevProps) {
    const search = this.props.search[this.props.reduxId];

    if (search.result.length === 1) {
      this.props.reset(this.props.reduxId);
      this.setState({ value: '' });

      return this.props.push(`/performer/${search.result[0].slug}`);
    }
  }

  componentDidMount() {
    document.addEventListener('keyup', this.onKeyDown);
  }

  componentWillUnmount() {
    document.removeEventListener('keyup', this.onKeyDown);

    if (this.input) {
      this.input.removeEventListener('mouseenter', this.onMouseEnter);
      this.input.removeEventListener('mouseout', this.onMouseOut);
    }
  }

  setInputRef = (input) => {
    if (!input) {
      return;
    }

    this.input = input.refs.input;

    this.input.addEventListener('mouseenter', this.onMouseEnter);
    this.input.addEventListener('mouseout', this.onMouseOut);
  };

  onMenuVisibilityChange = menuVisible => this.setState({ menuVisible });

  onMouseEnter = () => this.setState({ hoverActive: true });
  onMouseOut = () => this.setState({ hoverActive: false });

  render() {
    const { result, loading } = this.props.search[this.props.reduxId];

    return (
      <Autocomplete
        ref={ this.setInputRef }
        inputProps={ {
          className: this.props.className,
          id: this.state.id,
        } }
        wrapperStyle={ {
          float: 'right',
          display: 'inline-block',
        } }
        onMenuVisibilityChange={ this.onMenuVisibilityChange }
        open={ result.length > 1 && (this.state.hoverActive || this.state.menuVisible) }
        value={ this.state.value }
        items={ result }
        getItemValue={ (item) => {
          return item.stageName;
        } }
        autoHighlight={ false }
        onSelect={ (value, item) => {
          this.props.push(`/performer/${item.slug}`);
          this.props.reset(this.props.reduxId);
          this.setState({ value: '' });
        } }
        onChange={ (event, value) => {
          if (value.length < 3) {
            this.props.reset(this.props.reduxId);
            return this.setState({ value });
          } else {
            this.setState({ value });
          }

          this.props.fetchPerformers(this.props.reduxId, {
            query: value.replace(' ', '_'),
            limit: isMobile() ? 10 : 20,
            offset: 0,
            includeLatestRecording: 0,
          });
        } }
        renderMenu={ (items, value, style) => {
          if (!loading && items.length === 0) {
            items.push(<div className='suggestion'>No results found</div>);
          }

          return <div className='auto-complete-suggestions' children={ items }/>;
        } }
        renderItem={ (item: PerformerEntity, isHighlighted) => {
          const aliases = item.aliases.filter(a => {
            return (
              a.toLocaleLowerCase() !== item.stageName.toLocaleLowerCase() &&
              a.toLocaleLowerCase().indexOf(this.state.value.toLocaleLowerCase()) !== -1
            );
          });

          return (
            <div className={ 'suggestion ' + (isHighlighted ? 'highlighted' : '') } key={ item.uuid }>
              { item.stageName } { aliases.length > 0 && `(${aliases.join(', ')})` }
            </div>
          );
        } }
      />
    );
  }

  onKeyDown = (evt) => {
    const inputElement = document.getElementById(this.state.id);

    if (evt.which === 13 && inputElement === document.activeElement && this.state.value.length > 0) {
      evt.preventDefault();
      this.props.push({ pathname: '/search', search: `?query=${this.state.value}` });
      this.props.reset(this.props.reduxId);
      this.setState({ value: '' });
    }
  };
}
