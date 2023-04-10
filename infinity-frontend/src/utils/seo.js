
import { RecordingEntity, PerformerEntity } from 'api';
import { originServiceToString } from 'utils/performer';
import moment from 'moment';
import { RecordingSortMode } from 'store/ducks/recordings';

export function getTitleForPathname(pathname: string): string {
  switch (pathname) {
    case '/':
      return 'Latest cams - camtube.co';

    case '/most-viewed':
      return 'Most viewed cams - camtube.co';

    case '/most-popular':
      return 'Most popular cams - camtube.co';

    default:
      return 'Largest collection of cams - camtube.co ';
  }
}

export function getKeywordsOfRecording(recording: RecordingEntity): string[] {
  let keywords = [
    recording.stageName,
    originServiceToString(recording.performer.originService)
  ];

  if (recording.stageName !== recording.performer.stageName) {
    keywords.push(recording.performer.stageName);
  }

  keywords.forEach((keyword) => {
    if (keyword.indexOf('_') > -1) {
      keywords.push(keyword.replace('_', ' ').trim());
    }
  });

  return keywords;
}

export function getKeywordsOfRecordings(recordings: RecordingEntity[]): string[] {
  let keywords = [
    'camgirl',
    'webcam',
    'videos',
    'chaturbate',
    'cam4'
  ];
  recordings.forEach((r: RecordingEntity) => {
    const recordingKeywords = getKeywordsOfRecording(r);

    recordingKeywords.forEach((keyword: string) => {
      if (keywords.indexOf(keyword) < 0) {
        keywords.push(keyword);
      }
    });
  });

  return keywords;
}

export function getPostTitleOfRecording(recording: RecordingEntity): string[] {
  let parts = [
    recording.stageName,
    moment(recording.createdAt).format('DDMMYY HHmm'),
    originServiceToString(recording.performer.originService)
  ];

  if (recording.performer.originService === 'cbc') {
    parts.push(recording.performer.section);
  }

  return parts.join(' ');
}

export function getKeywordsOfPerformer(performer: PerformerEntity): string[] {
  let keywords = [
    performer.stageName,
    originServiceToString(performer.originService),
    'camgirl',
    'webcam',
    'videos',
    'chaturbate',
    'cam4'
  ];

  keywords.forEach((keyword) => {
    if (keyword.indexOf('_') > -1) {
      keywords.push(keyword.replace('_', ' ').trim());
    }
  });

  return keywords;
}

export function getKeywordsOfPerformers(performers: PerformerEntity[]): string[] {
  let keywords = [
    'camgirl',
    'webcam',
    'videos',
    'chaturbate',
    'cam4'
  ];

  performers.forEach((p: PerformerEntity) => {
    const performerKeywords = getKeywordsOfPerformer(p);

    performerKeywords.forEach((keyword: string) => {
      if (keywords.indexOf(keyword) < 0) {
        keywords.push(keyword);
      }
    });
  });

  return keywords;
}

export function getDescriptionOfSortedRecordings(currentPage, maxPage, sortMode) {
  let sorting = 'most recent';

  switch (sortMode) {
    case RecordingSortMode.VIEWS:
      sorting = 'most viewed';
      break;
    case RecordingSortMode.POPULARITY:
      sorting = 'most popular';
      break;
  }

  return `Page ${currentPage} of ${maxPage} for ${sorting} recordings`;
}

export function getDescriptionOfRecording(recording: RecordingEntity): string {
  return `A webcam recording of ${recording.performer.stageName} from ${moment(recording.createdAt).format('DD/MM/YY HH:mm')}`;
}

export function getDescriptionForCollectionPage(currentPage, maxPage, performer: PerformerEntity) {
  return `Page ${currentPage} of ${maxPage} for all the recordings available for the ${performer.stageName} on ${originServiceToString(performer.originService)}`;
}
