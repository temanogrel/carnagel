export class EmbeddedRecordingEntity {
  constructor(data) {
    this.uuid = data.uuid;
    this.videoUuid = data.videoUuid;
    this.collageUuid = data.collageUuid;
  }

  static create(data): EmbeddedRecordingEntity {
    return new EmbeddedRecordingEntity(data);
  }
}

export class PerformerEntity {
  constructor(data) {
    this.uuid = data.uuid;
    this.stageName = data.stageName;
    this.slug = data.slug;
    this.section = data.section;
    this.latestRecording = data.latestRecording ? EmbeddedRecordingEntity.create(data.latestRecording) : null;
    this.originService = data.originService;
    this.aliases = data.aliases;
    this.recordingCount = data.recordingCount;
  }

  static create(data): PerformerEntity {
    return new PerformerEntity(data);
  }
}
