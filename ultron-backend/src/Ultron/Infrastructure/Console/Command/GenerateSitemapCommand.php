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
use Ultron\Infrastructure\Service\SitemapServiceInterface;

final class GenerateSitemapCommand extends Command
{
    /**
     * @var SitemapServiceInterface
     */
    private $sitemapService;

    /**
     * GenerateSitemapCommand constructor.
     * @param SitemapServiceInterface $recordingRepository
     * @param $name
     */
    public function __construct(SitemapServiceInterface $recordingRepository, $name = null)
    {
        parent::__construct($name);

        $this->sitemapService = $recordingRepository;
    }

    protected function configure()
    {
        $this
            ->setName('ultron:sitemap-generate')
            ->setDescription('Create the sitemap for the performer, recordings');
    }

    protected function execute(InputInterface $input, OutputInterface $output)
    {
        $output->writeln('<info>Started performer sitemap generation</info>');
        $this->sitemapService->createPerformerSitemaps();

        $output->writeln('<info>Started recording sitemap generation</info>');
        $this->sitemapService->createRecordingSitemaps();
    }
}
