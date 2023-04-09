<?php
/**
 *
 *
 */

declare(strict_types = 1);

namespace Ultron\Infrastructure\Console\Command;

use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;
use Ultron\Domain\Service\CacheService;
use Ultron\Infrastructure\Repository\RecordingRepositoryInterface;

final class BuildCacheCommand extends Command
{
    /**
     * @var CacheService
     */
    private $cacheService;

    /**
     * @var RecordingRepositoryInterface
     */
    private $recordingRepository;

    /**
     * BuildCacheCommand constructor.
     * @param CacheService $cacheService
     * @param RecordingRepositoryInterface $recordingRepository
     * @param null $name
     */
    public function __construct(
        CacheService $cacheService,
        RecordingRepositoryInterface $recordingRepository,
        $name = null
    ) {
        parent::__construct($name);

        $this->cacheService        = $cacheService;
        $this->recordingRepository = $recordingRepository;
    }

    protected function configure()
    {
        $this
            ->setName('ultron:build-cache')
            ->setDescription('Build the cache');
    }

    protected function execute(InputInterface $input, OutputInterface $output)
    {
        $output->writeln('<info>Building cache</info>');
        $output->writeln('');

        $this->cacheService->buildCache($this->recordingRepository);
    }
}
