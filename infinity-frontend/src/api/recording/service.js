import { httpClient } from 'api';
import { RecordingEntity } from 'api/recording/entity';
import { CollectionResult } from 'api/stdlib';

export class RecordingService {
  static getByUuidOrSlug(uuidOrSlug: string): Promise<RecordingEntity> {
    return httpClient
      .get(`backend://recordings/${uuidOrSlug}`)
      .then(({ data }) => RecordingEntity.create(data));
  }

  static getByPerformer(uuidOrSlug: string, params: Object): Promise<CollectionResult<RecordingEntity>> {
    return httpClient
      .get(`backend://performers/${uuidOrSlug}/recordings`, { params })
      .then(({ data }) => {
        const items = data.data.map(RecordingEntity.create);

        return { items, meta: data.meta };
      });
  }

  static getAll(params): Promise<CollectionResult<RecordingEntity>> {
    return httpClient
      .get('backend://recordings', { params })
      .then(({ data }) => {
        const items = data.data.map(RecordingEntity.create);

        return { items, meta: data.meta };
      });
  }

  static getUserFavorites(id, params = {}): Promise<CollectionResult<RecordingEntity>> {
    return httpClient
      .get(`backend://users/${id}/favorites`, { params })
      .then(({ data }) => {
        const items = data.data.map(RecordingEntity.create);

        return { items, meta: data.meta };
      });
  }

  static addView(uuid: string): Promise<boolean> {
    return httpClient
      .post(`backend://rpc/recording.view`, { uuid })
      .then(({ data }) => data.added);
  }

  static toggleLike(uuid: string): Promise<boolean> {
    return httpClient
      .post('backend://rpc/recording.like', { uuid })
      .then(({ data }) => data.present);
  }

  static toggleFavorite(uuid: string): Promise<boolean> {
    return httpClient
      .post('backend://rpc/recording.favorite', { uuid })
      .then(({ data }) => data.present);
  }
}
