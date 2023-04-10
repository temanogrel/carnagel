import React, { Component } from 'react';
import { connect } from 'react-redux';
import { withRouter } from 'react-router-dom';
import { Helmet } from 'react-helmet';
import ReactGA from 'react-ga';

type Props = {
  title: string;
  description: string;
  keywords: string[];
  location: Object;
};

const mapStateToProps = state => ({
  title: state.seo.title,
  description: state.seo.description,
  keywords: state.seo.keywords,
});

@withRouter
@connect(mapStateToProps)
export class Seo extends Component<Props> {
  static defaultProps = {
    description: 'The largest collection of camtube videos',
  };

  componentDidMount() {
    if (__PROD__) {
      // Initialize google analytics if on prod
      ReactGA.initialize('UA-106135020-1');
    }
  }

  componentDidUpdate({ location }) {
    if (this.props.location.pathname !== location.pathname) {
      ReactGA.set({page: window.location.pathname});
      ReactGA.pageview(window.location.pathname);
    }
  }

  render() {
    return (
      <Helmet>
        { this.props.title && <title>{ this.props.title } </title> }
        { this.props.keywords.length > 0 && <meta name="keywords" content={this.props.keywords.join(',')}/> }
        { this.props.description && <meta name="description" content={this.props.description}/> }
      </Helmet>
    );
  }
}
