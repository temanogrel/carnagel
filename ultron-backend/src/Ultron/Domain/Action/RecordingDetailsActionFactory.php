<?php
/**
 *
 *
 *
 */

declare(strict_types = 1);

namespace Ultron\Domain\Action;

use Doctrine\ORM\EntityManager;
use Interop\Container\ContainerInterface;
use Ultron\Domain\Entity\RecordingEntity;
use Zend\Expressive\Template\TemplateRendererInterface;

class RecordingDetailsActionFactory
{
    /**
     * @param ContainerInterface $container
     *
     * @return RecordingDetailsAction
     */
    public function __invoke(ContainerInterface $container):RecordingDetailsAction
    {
        return new RecordingDetailsAction(
            $container->get(TemplateRendererInterface::class),
            $container->get(EntityManager::class)->getRepository(RecordingEntity::class)
        );
    }
}
