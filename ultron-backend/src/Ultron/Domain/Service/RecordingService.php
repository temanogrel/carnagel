<?php
/**
 *
 *
 *
 */

namespace Ultron\Domain\Service;

use Cocur\Slugify\Slugify;
use DateTime;
use Doctrine\Common\Persistence\ObjectManager;
use DomainException;
use Ultron\Domain\Entity\PerformerEntity;
use Ultron\Domain\Entity\RecordingEntity;
use Ultron\Domain\Entity\ValueObject\Images;
use Ultron\Domain\Exception\RecordingMissingAudioException;
use Ultron\Domain\Exception\RecordingNotFoundException;
use Ultron\Infrastructure\Repository\RecordingRepositoryInterface;
use Ultron\Infrastructure\Service\PageCacheServiceInterface;
use Zend\Stdlib\Parameters;
use Zend\Stdlib\ParametersInterface;

final class RecordingService implements RecordingServiceInterface
{
    const DESCRIPTION_PARSING_PATTERN = '/rate:\s(?P<bitrate>\d+)\skb\/s\sAudio:\s(?P<audio>(.*))\sVideo:\s(?P<video>(.*))/';

    /**
     * @var Slugify
     */
    private $slugify;

    /**
     * @var RecordingRepositoryInterface
     */
    private $recordingRepository;

    /**
     * @var ObjectManager
     */
    private $objectManager;

    /**
     * @var CacheService
     */
    private $cacheService;

    /**
     * @var PageCacheServiceInterface
     */
    private $pageCacheService;

    /**
     * @param CacheService                 $cacheService
     * @param ObjectManager                $objectManager
     * @param RecordingRepositoryInterface $recordingRepository
     * @param Slugify                      $slugify
     * @param PageCacheServiceInterface    $pageCacheService
     */
    public function __construct(
        CacheService $cacheService,
        ObjectManager $objectManager,
        RecordingRepositoryInterface $recordingRepository,
        Slugify $slugify,
        PageCacheServiceInterface $pageCacheService
    )
    {
        $this->slugify = $slugify;
        $this->cacheService = $cacheService;
        $this->objectManager = $objectManager;
        $this->recordingRepository = $recordingRepository;
        $this->pageCacheService = $pageCacheService;
    }

    /**
     * Convert the recording to a post title
     *
     * @param RecordingEntity $recording
     *
     * @return string
     */
    public static function generatePostTitle(RecordingEntity $recording)
    {
        $parts = [];

        $parts[] = $recording->getStageName();
        $parts[] = $recording->getCreatedAt()->format('dmy Hi');
        $parts[] = $recording->getPerformer()->getService(true);

        if ($recording->getPerformer()->getService() === 'cbc') {
            $parts[] = $recording->getPerformer()->getSection();
        }

        return implode(' ', $parts);
    }

    /**
     * Get keywords that are related to the given recording
     *
     * @param RecordingEntity $recording
     *
     * @return string[]
     */
    public static function getRecordingKeywords(RecordingEntity $recording)
    {
        $keywords = [
            $recording->getStageName(),
            $recording->getPerformer()->getService(true),
        ];

        if ($recording->getStageName() != $recording->getPerformer()->getStageName()) {
            $keywords[] = $recording->getPerformer()->getStageName();
        }

        foreach ($keywords as $keyword) {
            if (strpos($keyword, '_') !== false) {
                $keywords[] = trim(str_replace('_', ' ', $keyword));
            }
        }

        return $keywords;
    }

    /**
     * Remove a recording
     *
     * @param RecordingEntity $recording
     *
     * @return void
     */
    public function remove(RecordingEntity $recording)
    {
        $this->objectManager->remove($recording);
        $this->objectManager->flush();

        $this->cacheService->removeRecording($recording);
    }

    /**
     * @param ParametersInterface $data
     * @param PerformerEntity     $performer
     *
     * @throws RecordingMissingAudioException
     *
     * @return void
     */
    public function create(ParametersInterface $data, PerformerEntity $performer)
    {
        if ($data->get('audio') === null) {
            throw new RecordingMissingAudioException();
        }

        $recording = new RecordingEntity();
        $recording->setPerformer($performer);
        $recording->setUid($data['id']);
        $recording->setStageName($performer->getStageName());
        $recording->setDuration((int)$data->get('duration', null));
        $recording->setVideoUrl($data->get('videoUrl'));
        $recording->setBitRate($data->get('bitRate'));
        $recording->setAudio($data->get('audio'));
        $recording->setVideo($data->get('video'));
        $recording->setImageUrls(Images::fromParameters(new Parameters($data->get('imageUrls'))));
        $recording->setCreatedAt(new DateTime($data->get('createdAt')));
        $recording->setUpdatedAt(new DateTime($data->get('updatedAt')));

        if ($data->get('encoding') === 'h264') {
            $recording->setSize((int)$data->get('size264'));
        } else {
            $recording->setSize((int)$data->get('size265'));
        }

        $performer->incrementRecordingCount();
        $performer->setUpdatedAt(new DateTime());

        $this->updateWithSlug($recording);

        $this->objectManager->persist($recording);
        $this->objectManager->flush();

        $this->cacheService->addRecording($recording);
        $this->pageCacheService->addRecording($recording);
    }

    /**
     * {@inheritdoc}
     */
    private function updateWithSlug(RecordingEntity $recording)
    {
        $slug = $this->slugify->slugify(self::generatePostTitle($recording));
        $occurrences = 1;

        if ($recording->getSlug() === $slug) {
            return;
        }

        while (true) {

            try {

                // Try and get by it's slug, if we fail it's available
                $this->recordingRepository->getBySlug($slug);

                // Generate a new slug and try again
                $slug = $this->slugify->slugify(self::generatePostTitle($recording) . '-' . $occurrences++);
            } catch (RecordingNotFoundException $e) {
                $recording->setSlug($slug);
                break;
            }
        }
    }
}
