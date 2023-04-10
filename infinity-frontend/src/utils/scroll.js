export function scrollTo(x: number = 0, y: number = 0) {
  if (typeof window.scroll === 'function') {
    window.scroll({
      top: y,
      left: x,
      behavior: 'smooth',
    });
  } else if (typeof window.scrollTo === 'function') {
    window.scrollTo(x, y);
  }
}
