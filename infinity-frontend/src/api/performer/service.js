import { httpClient } from 'api';
import { CollectionResult } from 'api/stdlib';
import { PerformerEntity } from 'api/performer/entity';

export class PerformerService {
  static getByUuidOrSlug(id: string): Promise<PerformerEntity> {
    return httpClient
      .get(`backend://performers/${id}`)
      .then(({ data }) => PerformerEntity.create(data));
  }

  static search(params): Promise<CollectionResult<PerformerEntity>> {
    return httpClient
      .get('backend://performers', { params })
      .then(({ data }) => {
        const items = data.data.map(PerformerEntity.create);
        const meta = data.meta;

        return {items, meta};
      });
  }
}
