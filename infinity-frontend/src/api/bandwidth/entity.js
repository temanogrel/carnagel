export class BandwidthStatusEntity {
  constructor(data) {
    this.total = data.total;
    this.remaining = data.remaining;

    return Object.freeze(this);
  }

  static create(data): BandwidthStatusEntity {
    return new BandwidthStatusEntity(data);
  }
}
