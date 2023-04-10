export function isPrerender() {
  return /(PhantomJS|ChromeHeadless|Prerender)/.test(window.navigator.userAgent);
}
