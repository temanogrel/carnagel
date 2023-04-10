export const cloneObject = obj => Object.assign(Object.create(Object.getPrototypeOf(obj)), obj);
