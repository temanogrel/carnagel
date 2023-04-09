<?php
/**
 *
 *
 *
 */

declare(strict_types = 1);

namespace Ultron\Domain\Service;

use Ultron\Domain\Entity\PerformerEntity;
use Zend\Stdlib\ParametersInterface;

interface PerformerServiceInterface
{
    public function create(ParametersInterface $data):PerformerEntity;
    public function update(PerformerEntity $performer);
}
