import createHistory from 'history/createBrowserHistory';

// This helper file is needed because we can only use createHistory() once, but we need it
// when creating react-router-redux middleware, and when creating our Router component */
export default createHistory();
