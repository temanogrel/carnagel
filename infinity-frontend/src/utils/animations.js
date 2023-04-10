export class AnimateUp {

  scrollTop = 0;
  interval = null;
  animating = false;
  instant = false;

  constructor(instant) {
    this.instant = instant; // Animate all elements without waiting for scroll
    this.addListeners();
    this.handleMoveEvent();
    window.setInterval(this.handleMoveEvent.bind(this), 300);
  }

  addListeners() {
    window.addEventListener('scroll', this.handleMoveEvent.bind(this), false);
    window.addEventListener('touchmove', this.handleMoveEvent.bind(this), false);
    window.clearInterval(this.getPaymentTransactionInterval);
  }

  handleMoveEvent() {
    this.scrollTop = window.scrollY ? window.scrollY : window.pageYOffset;

    let elements = this.findElementsToAnimateUp();

    if (elements.length > 0 && !this.animating) {
      this.animate(elements);
    }
  }

  findElementsToAnimateUp() {
    let visibleElements = [];
    let elements = document.querySelectorAll('.animate-up-opacity');

    if (this.instant) {
      return elements;
    }

    for (let elem of elements) {
      let offset = this.getOffset(elem);
      if (offset.top < window.innerHeight) {
        visibleElements.push(elem);
      }
    }

    return visibleElements;
  }

  getOffset(el) {
    let viewportOffset = el.getBoundingClientRect();

    return {
      top: viewportOffset.top,
      left: viewportOffset.left
    };
  }

  animate(elems) {
    let current, interval, removeClass;

    current = 0;
    removeClass = function() {
      elems[current].classList.remove('animate-up-opacity');

      current++;

      if (current === elems.length) {
        window.clearInterval(interval);
        this.animating = false;
      }
    };

    interval = window.setInterval(removeClass.bind(this), 10);
  }

  destroy() {
    window.removeEventListener('scroll', this.handleMoveEvent, false);
    window.removeEventListener('touchmove', this.handleMoveEvent, false);
  }
}
