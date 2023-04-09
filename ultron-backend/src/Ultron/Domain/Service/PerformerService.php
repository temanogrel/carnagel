<?php
/**
 *
 *
 *
 */

declare(strict_types = 1);

namespace Ultron\Domain\Service;

use Cocur\Slugify\Slugify;
use DateTime;
use Doctrine\Common\Persistence\ObjectManager;
use Ultron\Domain\Entity\PerformerEntity;
use Ultron\Domain\Exception\PerformerNotFoundException;
use Ultron\Infrastructure\Repository\PerformerRepositoryInterface;
use Zend\Stdlib\ParametersInterface;

class PerformerService implements PerformerServiceInterface
{
    /**
     * @var ObjectManager
     */
    private $objectManager;

    /**
     * @var PerformerRepositoryInterface
     */
    private $performerRepository;

    /**
     * @var Slugify
     */
    private $slugify;

    /**
     * PerformerService constructor.
     *
     * @param ObjectManager                $objectManager
     * @param PerformerRepositoryInterface $performerRepository
     * @param Slugify                      $slugify
     */
    public function __construct(
        ObjectManager $objectManager,
        PerformerRepositoryInterface $performerRepository,
        Slugify $slugify
    ) {
        $this->slugify             = $slugify;
        $this->objectManager       = $objectManager;
        $this->performerRepository = $performerRepository;
    }

    /**
     *
     *
     * @param PerformerEntity $performer
     *
     * @return void
     */
    private function updateWithSlug(PerformerEntity $performer)
    {
        $generator = function (PerformerEntity $performer) {

            $parts = [];

            switch ($performer->getService()) {
                case 'cam':
                    $parts[] = 'cam4';
                    break;

                case 'cbc':
                    $parts[] = 'chaturbate';
                    break;

                case 'mfc':
                    $parts[] = 'myfreecams';
                    break;
            }

            $parts[] = $performer->getStageName();

            return $this->slugify->slugify(implode('-', $parts));
        };

        $slug = $generator($performer);

        // short circuit if the slug is already set and matches.
        if ($performer->getSlug() === $slug) {
            return;
        }

        $occurrences = 1;

        while (true) {
            try {
                $this->performerRepository->getBySlug($slug);

                $slug = $generator($performer) . '-' . $occurrences++;

            } catch (PerformerNotFoundException $e) {
                $performer->setSlug($slug);
                break;
            }
        }
    }

    /**
     * {@inheritdoc}
     */
    public function create(ParametersInterface $data):PerformerEntity
    {
        $performer = new PerformerEntity();
        $performer->setUid((int)$data->get('id'));
        $performer->setStageName((string)$data->get('stageName'));
        $performer->setAliases($data->get('aliases'));
        $performer->setService((string)$data->get('service'));
        $performer->setSection((string)$data->get('section'));
        $performer->setCreatedAt(new DateTime());
        $performer->setUpdatedAt(new DateTime());

        $this->updateWithSlug($performer);

        $this->objectManager->persist($performer);
        $this->objectManager->flush();

        return $performer;
    }

    /**
     * {@inheritdoc}
     */
    public function update(PerformerEntity $performer)
    {
        $performer->setUpdatedAt(new DateTime());

        $this->updateWithSlug($performer);

        $this->objectManager->flush();
    }
}
