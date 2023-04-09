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
use Ultron\Domain\Sites;
use Ultron\Infrastructure\Service\PageCacheServiceInterface;

final class BuildPageCacheCommand extends Command
{
    /**
     * @var PageCacheServiceInterface
     */
    private $pageCacheService;

    /**
     * @var Sites
     */
    private $sites;

    /**
     * GeneratePageCacheCommand constructor.
     * @param PageCacheServiceInterface $pageCacheService
     * @param Sites $sites
     * @param param null $name
     */
    public function __construct(PageCacheServiceInterface $pageCacheService, Sites $sites, $name = null)
    {
        parent::__construct($name);

        $this->pageCacheService = $pageCacheService;
        $this->sites            = $sites;
    }

    protected function configure()
    {
        $this
            ->setName('ultron:build-page-cache')
            ->setDescription('Build the cache for pages');
    }

    protected function execute(InputInterface $input, OutputInterface $output)
    {
        $output->writeln('<info>Building page cache</info>');
        $output->writeln('');

        $sitesConfig = $this->sites->getSiteConfigurations();

        foreach ($sitesConfig as $site) {
            if (!$site->isEnabled()) {
                $output->writeln(sprintf('<warning>Skipping %s because it\'s not enabled</warning>', $site->getDomain()));
                continue;
            }

            $output->writeln('<info>Creating page cache for domain: ' . $site->getDomain() . '</info>');

            $this->pageCacheService->create($site);
        }
    }
}
