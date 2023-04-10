import React, { Component } from 'react';
import { Link, NavLink, withRouter } from 'react-router-dom';
import { connect } from 'react-redux';
import { MODAL_ACCOUNT, MODAL_BANDWIDTH_EXPLAINED, MODAL_LOGIN, pushModal } from 'store/ducks/modals';
import { PerformerSearch } from '../search/performer';
import fileSize from 'filesize';

type Props = {};

type State = {
  visible: boolean;
};

@withRouter
class BackToTop extends Component<Props, State> {
  static displayOnPathnames = [
    '/',
    '/most-viewed',
    '/most-popular',
    '/favorites'
  ];

  constructor(props: Props) {
    super(props);

    this.state = {
      visible: false
    };
  }

  componentDidMount() {
    document.addEventListener('scroll', () => {
      const visible = window.scrollY || document.documentElement.scrollTop > 0;

      if (visible) {
        clearTimeout(this.timeout);

        this.timeout = setTimeout(() => {
          this.setState({visible: false});
        }, 2500);
      }

      this.setState({visible: visible});
    });
  }

  shouldComponentUpdate(nextProps, nextState) {
    return nextState.visible !== this.state.visible;
  }

  render() {
    const pathname = this.props.location.pathname;

    if (BackToTop.displayOnPathnames.indexOf(pathname) < 0 && pathname.indexOf('/performer') !== 0) {
      return null;
    }

    return <a className={`backtotop ${this.state.visible ? 'active' : ''}`} onClick={this.backToTop}>Back to top</a>;
  }

  backToTop = () => {
    clearTimeout(this.timeout);

    const interval = setInterval(() => {
      let offsetY = window.scrollY || document.documentElement.scrollTop;
      window.scrollTo(0, offsetY - 25);

      if (offsetY === 0) {
        clearInterval(interval);
      }
    }, 5);
  };
}

export const Header = ({user, login, account, bandwidth, explainBandwidth}) => {
  const menuText = user ? 'Account' : 'Login';
  const menuButton = user ? account : login;

  // pretty print the remaining bandwidth
  const [size, notation] = fileSize(bandwidth.remaining, {output: 'array', round: 0});

  // calculate the remaining bandwidth
  const percentage = Math.floor(((bandwidth.total - bandwidth.remaining) / bandwidth.total) * 100);

  // set the percentage class
  const classes = 'radial-loader-container float-right p-' + percentage;

  return (
    <header>
      <BackToTop/>
      <div className='header mobile-hidden'>
        <Link to='/' className='logo float-left'/>
        <a onClick={menuButton} className='float-right button red large mr30 user-button'>{menuText}</a>
        <div className={classes}>
          <div className='radial-loader' onClick={explainBandwidth}>
            <span>{size}</span>
            <div>{notation}</div>
          </div>
        </div>
        <PerformerSearch reduxId="header" className={'search button red large mr15 transparent collapsable text-left'}/>
        <ul className='header-menu mobile-hidden'>
          <li><NavLink to='/'>Newest</NavLink></li>
          <li><NavLink to='/most-viewed'>Most viewed</NavLink></li>
          <li><NavLink to='/most-popular'>Most popular</NavLink></li>
        </ul>
        <div className='clear'/>
      </div>

      <div className='clear mobile-header hidden mobile-block'>
        <button className='float-left button red medium mt25 ml25 mb25 hamburger-button'>
          <ul className='float-left mobile-hidden'>
            <li><NavLink to='/'>Newest</NavLink></li>
            <li><NavLink to='/most-viewed'>Most viewed</NavLink></li>
            <li><NavLink to='/most-popular'>Most popular</NavLink></li>
          </ul>
        </button>
        <Link to='/' className='logo'/>
        <a onClick={menuButton} className='float-right button red medium mt25 mr25 mb25 user-button'/>
        <div className='band-width-mobile clear'>
          <span>Bandwidth left {size} {notation}</span>
          <div style={{width: percentage + '%'}}/>
        </div>
        <div className='search-container'>
          <PerformerSearch reduxId="header" className={'search button white small transparent text-left'}/>
          <div className='clear'/>
        </div>
      </div>
    </header>
  );
};

const mapStateToProps = (state) => ({
  user: state.user,
  bandwidth: state.bandwidth,
});

const mapDispatchToProps = (dispatch) => ({
  login: () => dispatch(pushModal(MODAL_LOGIN)),
  account: () => dispatch(pushModal(MODAL_ACCOUNT)),
  explainBandwidth: () => dispatch(pushModal(MODAL_BANDWIDTH_EXPLAINED)),
});

export const HeaderContainer = connect(mapStateToProps, mapDispatchToProps)(Header);
