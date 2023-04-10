import environment from 'api/environment';
import { ProxyService } from 'api';
import moment from 'moment';

export class EmbeddedPerformerEntity {
  constructor(data) {
    this.uuid = data.uuid;
    this.stageName = data.stageName;
    this.slug = data.slug;
    this.recordingCount = data.recordingCount;
    this.originService = data.originService;
    this.section = data.section;

    return Object.freeze(this);
  }

  static create(data): EmbeddedPerformerEntity {
    return new EmbeddedPerformerEntity(data);
  }
}

export class RecordingImageEntity {
  constructor(uuid) {
    this.uuid = uuid;

    return Object.freeze(this);
  }

  get url() {
    return `url(${ProxyService.getImageUrl(this)})`;
  }

  get rawUrl() {
    return ProxyService.getImageUrl(this);
  }

  static create(uuid): RecordingImageEntity {
    return new RecordingImageEntity(uuid);
  }
}

export class RecordingEntity {
  constructor(data) {
    this.uuid = data.uuid;
    this.videoUuid = data.videoUuid;
    this.collageUuid = data.collageUuid;
    this.performer = EmbeddedPerformerEntity.create(data.performer);
    this.images = data.images.map(RecordingImageEntity.create);
    this.sprites = data.sprites.map(RecordingImageEntity.create);
    this.stageName = data.stageName;
    this.duration = data.duration;
    this.viewCount = data.viewCount;
    this.likeCount = data.likeCount;
    this.isLiked = data.isLiked;
    this.isFavorite = data.isFavorite;
    this.createdAt = moment(data.createdAt);
    this.slug = data.slug;

    return Object.freeze(this);
  }

  get manifest() {
    return `${environment.apiUri}/recordings/${this.uuid}/manifest.m3u8`;
  }

  get collageRawUrl() {
    return ProxyService.getImageUrl(RecordingImageEntity.create(this.collageUuid));
  }

  get collageUrl() {
    return `url(${ProxyService.getImageUrl(RecordingImageEntity.create(this.collageUuid))})`;
  }

  /**
   * The recording is minified if we only have a single image as a thumb
   *
   * @returns {boolean}
   */
  get isMinified(): boolean {
    return this.images.length === 1;
  }

  durationAsString() {
    const hours = Math.floor(this.duration / 60 / 60),
      minutes = Math.floor((this.duration - (hours * 60 * 60)) / 60),
      seconds = Math.round(this.duration - (hours * 60 * 60) - (minutes * 60));

    return hours + ':' + ((minutes < 10) ? '0' + minutes : minutes) + ':' + ((seconds < 10) ? '0' + seconds : seconds);
  }

  toggleFavorite(isFavorite: boolean) {
    this.isFavorite = isFavorite;
  }

  addView() {
    this.viewCount++;
  }

  toggleLike(isLiked: boolean) {
    this.isLiked = isLiked;

    if (isLiked) {
      this.likeCount++;
    } else {
      this.likeCount--;
    }
  }

  static create(data): RecordingEntity {
    return new RecordingEntity(data);
  }
}
