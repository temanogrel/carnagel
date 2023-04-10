// ------------------------------------
// Constants
// ------------------------------------
export const SEO_UPDATE = 'seo:update';

// ------------------------------------
// Actions
// ------------------------------------
export const seoUpdate = (title = '', description = '', keywords = []) =>
  ({ type: SEO_UPDATE, payload: { title, description, keywords } });

// ------------------------------------
// Reducer
// ------------------------------------
const initialState = {
  title: 'Largest collection of cams - camtube.co',
  description: '',
  keywords: [],
};

export default (state = initialState, action) => {
  switch (action.type) {
    case SEO_UPDATE:
      return action.payload;

    default:
      return state;
  }
};
