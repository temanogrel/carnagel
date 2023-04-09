<?php
use Hermes\Service\UpstoreServiceInterface;

/**
 *
 *
 *
 */

namespace Hermes\Command\Sync;

use Doctrine\Common\Collections\Criteria;
use Doctrine\Common\Persistence\ObjectManager;
use Hermes\Entity\UrlEntity;
use Hermes\Repository\UrlRepositoryInterface;
use Hermes\Service\UpstoreServiceInterface;
use Knp\Command\Command;
use Symfony\Component\Console\Helper\ProgressBar;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;

class UpstoreSync extends Command
{
    protected function configure()
    {
        $this
            ->setName('hermes:sync:upstore')
            ->setDescription('Sync the downlink hash for all the entries');
    }

    protected function execute(InputInterface $input, OutputInterface $output)
    {
        /* @var $objectManager ObjectManager */
        $objectManager = $this->getSilexApplication()['objectManager'];

        /* @var $repository UrlRepositoryInterface */
        $repository = $objectManager->getRepository(UrlEntity::class);

        /* @var $upStoreService UpstoreServiceInterface */
        $upStoreService = $this->getSilexApplication()['service.upstore'];

        $criteria = Criteria::create();
        $criteria->andWhere($criteria->expr()->isNull('upstoreDownloadHash'));
        $criteria->andWhere($criteria->expr()->eq('isUpstore', true));

        $count = (int) $repository->getCount(clone $criteria);
        if ($count === 0) {
            return $output->writeln('<info>Nothing to sync</info>');
        }

        $segments = ceil($count / 1000);

        $progress = new ProgressBar($output, $count);
        $progress->setRedrawFrequency($count  / 1000);
        $progress->start();

        for ($i = 0; $i <= $segments; $i++) {

            $urls = $repository->findBy(['upstoreDownloadHash' => null, 'isUpstore' => true], null, 1000, $i * 1000);

            $upStoreService->syncUrlCollection($urls);

            $objectManager->flush();
            $objectManager->clear();

            $progress->setProgress($i * 1000);
        }

        $progress->finish();
    }
}
