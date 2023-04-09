<?php
/**
 *
 *
 */

declare(strict_types = 1);

use Ultron\Infrastructure\Console\Command\BuildCacheCommand;
use Ultron\Infrastructure\Console\Command\BuildPageCacheCommand;
use Ultron\Infrastructure\Console\Command\GenerateSitemapCommand;
use Ultron\Infrastructure\Console\Command\RebuildPerformerRecordingCountCommand;

return [
    'console' => [
        'commands' => [
            BuildPageCacheCommand::class,
            BuildCacheCommand::class,
            GenerateSitemapCommand::class,
            RebuildPerformerRecordingCountCommand::class,
        ],
    ]
];
