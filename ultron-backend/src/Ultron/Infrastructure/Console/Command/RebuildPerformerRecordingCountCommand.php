<?php
/**
 *
 *
 */

declare(strict_types = 1);

namespace Ultron\Infrastructure\Console\Command;

use DateTime;
use Doctrine\ORM\EntityManager;
use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Helper\ProgressBar;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;
use Ultron\Infrastructure\Repository\PerformerRepositoryInterface;
use Ultron\Infrastructure\Repository\RecordingRepositoryInterface;

final class RebuildPerformerRecordingCountCommand extends Command
{
    /**
     * @var RecordingRepositoryInterface
     */
    private $recordingRepository;

    /**
     * @var PerformerRepositoryInterface
     */
    private $performerRepository;

    /**
     * @var EntityManager
     */
    private $entityManager;

    /**
     * RebuildPerformerRecordingCountCommand constructor.
     * @param RecordingRepositoryInterface $recordingRepository
     * @param PerformerRepositoryInterface $performerRepository
     * @param EntityManager $entityManager
     * @param null $name
     */
    public function __construct(
        RecordingRepositoryInterface $recordingRepository,
        PerformerRepositoryInterface $performerRepository,
        EntityManager $entityManager,
        $name = null
    ) {
        parent::__construct($name);

        $this->recordingRepository = $recordingRepository;
        $this->performerRepository = $performerRepository;
        $this->entityManager       = $entityManager;
    }

    protected function configure()
    {
        $this
            ->setName('ultron:performer-recording-count:rebuild')
            ->setDescription('Rebuild recording count for performers');
    }

    protected function execute(InputInterface $input, OutputInterface $output)
    {
        $output->writeln('<info>Started rebuilding performer recording count</info>');
        $output->writeln('');
        $output->writeln('<info>Calculating recording count of performers</info>');
        $output->writeln('');

        $result = $this->recordingRepository->getRecordingCountOfPerformers();

        $output->writeln('<info>Updating recording count for all performers</info>');
        $output->writeln('<info>There are ' . count($result) . ' performers with at least one recording</info>');
        $output->writeln('');

        $progress = new ProgressBar($output, count($result));

        foreach ($result as $performerRecordingCount) {
            $performer = $this->performerRepository->getById($performerRecordingCount['performerId']);
            $performer->setRecordingCount((int) $performerRecordingCount['recordingCount']);
            $performer->setUpdatedAt(new DateTime());

            $progress->advance();
        }

        $progress->finish();

        $output->writeln('');
        $output->writeln('Writing updated recording count to database');
        $output->writeln('');

        $this->entityManager->flush();
    }
}
