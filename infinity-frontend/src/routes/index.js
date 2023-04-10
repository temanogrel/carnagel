// We only need to import the modules necessary for initial render
import { AppContainer } from 'containers/app';
import React, { Component } from 'react';
import L from 'react-loadable';
import { Redirect } from 'react-router';
import { Route, Switch } from 'react-router-dom';
import { ScrollToTop, GoogleAnalytics, RecordingListSortingWrapper } from './utils';

/* eslint-disable */
export const Loadable = opts =>
  L({
    loading: () => {
      return (
        <div className="loading-cover">
          <div className="loading-spinner"/>
        </div>
      );
    },
    ...opts,
  });

const RecordingList = Loadable({ loader: () => import(/* webpackChunkName: 'recording-list' */ '../views/recording/list') });
const RecordingWrapper = Loadable({ loader: () => import(/* webpackChunkName: 'recording-wrapper' */ '../views/recording/wrapper') });
const RecordingView = Loadable({ loader: () => import(/* webpackChunkName: 'recording-view' */ '../views/recording/view') });
const PerformerWrapper = Loadable({ loader: () => import(/* webpackChunkName: 'performer-wrapper' */ '../views/performer/wrapper') });
const PerformerRecordings = Loadable({ loader: () => import(/* webpackChunkName: 'performer-recordings' */ '../views/performer/recordings') });
const PerformerSearchResult = Loadable({ loader: () => import(/* webpackChunkName: 'performer-search-result' */ '../views/performer/search-result') });
const PageNotFound = Loadable({ loader: () => import(/* webpackChunkName: 'page-not-found' */ '../views/404/not-found') });
const UserFavorites = Loadable({ loader: () => import(/* webpackChunkName: 'user-favorites */ '../views/user/favorites') });
const UserWrapper = Loadable({ loader: () => import(/* webpackChunkName: 'user-wrapper' */ '../views/user/wrapper') });
const PaymentPlans = Loadable({ loader: () => import(/* webpackChunkName: 'payment-plans' */ '../views/user/payment-plans') });
const SetNewPasswordWrapper = Loadable({ loader: () => import(/* webpackChunkName: 'set-new-password-wrapper' */ '../views/user/set-new-password-wrapper') });
const SessionWrapper = Loadable({ loader: () => import(/* webpackChunkName: 'session-wrapper' */ '../views/session/wrapper') });

export default class Router extends Component {
  render() {
    return (
      <AppContainer>
        <GoogleAnalytics>
          <ScrollToTop>
            <SessionWrapper>
              <Switch>
                <Route exact path="/recording/:id">
                  <RecordingWrapper>
                    <Route exact path="/recording/:id" component={ RecordingView }/>
                  </RecordingWrapper>
                </Route>

                <Route path="/performer/:id">
                  <PerformerWrapper>
                    <Route exact path="/performer/:id" component={ PerformerRecordings }/>
                  </PerformerWrapper>
                </Route>

                <Route exact path="/">
                  <RecordingListSortingWrapper>
                    <Route exact path="/" component={ RecordingList }/>
                  </RecordingListSortingWrapper>
                </Route>

                <Route exact path="/most-viewed">
                  <RecordingListSortingWrapper>
                    <Route exact path="/most-viewed" component={ RecordingList }/>
                  </RecordingListSortingWrapper>
                </Route>

                <Route exact path="/most-popular">
                  <RecordingListSortingWrapper>
                    <Route exact path="/most-popular" component={ RecordingList }/>
                  </RecordingListSortingWrapper>
                </Route>

                <Route exact path="/search" component={ PerformerSearchResult }/>

                <Route exact path="/favorites">
                  <UserWrapper>
                    <Route exact path="/favorites" component={ UserFavorites }/>
                  </UserWrapper>
                </Route>

                <Route exact path="/payment-plans" component={ PaymentPlans }/>

                <Route exact path="/password-reset">
                  <SetNewPasswordWrapper>
                    <Route exact path="/password-reset" render={ () => <Redirect to="/"/> }/>
                  </SetNewPasswordWrapper>
                </Route>

                <Route path="/404" component={ PageNotFound }/>
                <Route path="*" render={ () => <Redirect to="/404"/> }/>
              </Switch>
            </SessionWrapper>
          </ScrollToTop>
        </GoogleAnalytics>
      </AppContainer>
    );
  }
}
