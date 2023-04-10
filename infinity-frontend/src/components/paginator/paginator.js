import React, { Component } from 'react';
import { connect } from 'react-redux';
import { push } from 'react-router-redux';

type Props = {
  total: number;
  offset: number;
  limit: number;
  push: push;
  additionalQuery?: string;
};

type State = {
  pagesInRange: number[];
  currentPage: number;
  pages: number;
  pageSelection: number;
};

const mapDispatchToProps = dispatch => ({
  push: url => dispatch(push(url)),
});

@connect(null, mapDispatchToProps)
export class Paginator extends Component<Props, State> {
  constructor(props) {
    super(props);

    this.state = {
      pagesInRange: [],
      currentPage: 0,
      pages: 0,
      pageSelection: 0
    };

    this.goTo = this.goTo.bind(this);
    this.goToPageSelection = this.goToPageSelection.bind(this);
    this.onPageSelectionKeyDown = this.onPageSelectionKeyDown.bind(this);
  }

  componentDidMount() {
    this.calculate(this.props);
  }

  componentWillReceiveProps(nextProps) {
    this.calculate(nextProps);
  }

  calculate(props) {
    const pageCount = Math.ceil(props.total / props.limit);
    const pageNumber = props.offset === 0 ? 1 : (props.offset / props.limit) + 1;

    let pageRange = 5;
    let lowerBound = 0;
    let upperBound = 1;

    if (pageRange > pageCount) {
      pageRange = pageCount;
    }

    let delta = Math.ceil(pageRange / 2);

    if (pageNumber - delta > pageCount - pageRange) {
      lowerBound = pageCount - pageRange + 1;
      upperBound = pageCount;
    } else {
      if (pageNumber - delta < 0) {
        delta = pageNumber;
      }

      const offset = pageNumber - delta;
      lowerBound = offset + 1;
      upperBound = offset + pageRange;
    }

    const range = [];
    for (; lowerBound <= upperBound; lowerBound++) {
      range.push(lowerBound);
    }

    this.setState({ pages: pageCount, currentPage: pageNumber, pagesInRange: range });
  }

  goTo(page: number) {
    const additionalQuery = this.props.additionalQuery ? this.props.additionalQuery : '';

    return () => {
      this.props.push(`${window.location.pathname}?page=${page}${additionalQuery}`);
    };
  }

  renderPagesInRange() {
    return this.state.pagesInRange.map(p => (
      <div className='grid-1 mobile-grid-2' key={p}>
        <a onClick={this.goTo(p)} className={this.state.currentPage === p ? 'active' : ''}>{p}</a>
      </div>
    ));
  }

  renderPrevious() {
    const anchor = this.state.currentPage === 1 ? <a>&lt;</a>
      : <a onClick={this.goTo(this.state.currentPage - 1)}>&lt;</a>;

    return (
      <div className='grid-1 mobile-grid-1'>{anchor}</div>
    );
  }

  renderNext() {
    const anchor = this.state.currentPage === this.state.pages ? <a>&gt;</a>
      : <a onClick={this.goTo(this.state.currentPage + 1)}>&gt;</a>;

    return (
      <div className='grid-1 mobile-grid-1'>{anchor}</div>
    );
  }

  onPageSelectionKeyDown(evt) {
    if (evt.keyCode === 13) {
      this.goTo(evt.target.value)();
    } else {
      this.setState({ pageSelection: parseInt(evt.target.value) });
    }
  }

  goToPageSelection() {
    this.goTo(this.state.pageSelection)();
  }

  render() {
    if (isNaN(this.state.pages) || this.state.pages <= 1) {
      return false;
    }

    return (
      <div className='footer'>
        { this.renderPrevious() }
        { this.renderPagesInRange() }
        { this.renderNext() }
        <div className='grid-2 mobile-grid-4'><p className='numberofpages special'>{'..' + this.state.pages}</p></div>
        <div className='grid-2 mobile-grid-4'>
          <p className='special'>
            <input type='number' placeholder='Page nr.' onKeyUp={this.onPageSelectionKeyDown}/>
          </p>
        </div>
        <div className='grid-1 mobile-grid-4'><a className='special' onClick={this.goToPageSelection}>Go</a></div>
      </div>
    );
  }
}
