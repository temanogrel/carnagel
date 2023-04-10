<?php
/**
 *
 *
 *  AB
 */

use Aphrodite\Site\Entity\PostAssociation;
use Zend\Stdlib\ArrayUtils;

$publishedOn = array_map(function (PostAssociation $association) {
    return [
        'id'     => $association->getId(),
        'site'   => $association->getSite()->getId(),
        'postId' => $association->getPostId(),
    ];

}, ArrayUtils::iteratorToArray($this->recording->getPublishedOn()));

$lastCheckedAt = $this->recording->getLastCheckedAt();

$data = [
    'id'          => $this->recording->getId(),
    'oldId'       => $this->recording->getOldId(),
    'performerId' => $this->recording->getPerformer() !== null ? $this->recording->getPerformer()->getId() : null,
    'state'       => $this->recording->getState(),
    'stageName'   => $this->recording->getStageName(),
    'section'     => $this->recording->getSection(),
    'service'     => $this->recording->getService(),
    'duration'    => $this->recording->getDuration(),

    // todo: remove the following fields when we have updated the rest of the system
    'size'        => $this->recording->getSize264(),
    'videoUuid'   => $this->recording->getVideoMp4Uuid(),

    // new style size versions due to AP-11
    'size264'     => $this->recording->getSize264(),
    'size265'     => $this->recording->getSize265(),
    'bitRate'     => $this->recording->getBitRate(),
    'encoding'    => $this->recording->getEncoding(),
    'audio'       => $this->recording->getAudio(),
    'video'       => $this->recording->getVideo(),
    'videoUrl'    => $this->recording->getVideoUrl(),
    'orphaned'    => $this->recording->isOrphaned(),

    'videoMp4Uuid'         => $this->recording->getVideoMp4Uuid(),
    'videoHlsUuid'         => $this->recording->getVideoHlsUuid(),
    'videoManifest'        => $this->recording->getVideoManifest(),
    'wordpressCollageUuid' => $this->recording->getWordpressCollageUuid(),
    'infinityCollageUuid'  => $this->recording->getInfinityCollageUuid(),
    'images'               => $this->recording->getImages(),
    'sprites'              => $this->recording->getSprites(),
    'upstoreHash'          => $this->recording->getUpstoreHash(),

    'imageUrls' => [
        'thumb'   => $this->recording->getImageUrls()->getThumb(),
        'large'   => $this->recording->getImageUrls()->getLarge(),
        'gallery' => $this->recording->getImageUrls()->getGallery(),
    ],

    'publishedOn'   => $publishedOn,
    'lastCheckedAt' => $lastCheckedAt !== null ? $lastCheckedAt->format(DateTime::RFC3339) : null,

    'createdAt' => $this->recording->getCreatedAt()->format(DateTime::RFC3339),
    'updatedAt' => $this->recording->getUpdatedAt()->format(DateTime::RFC3339),
];

return $data;
