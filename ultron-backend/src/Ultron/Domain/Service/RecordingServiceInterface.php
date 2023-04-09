<?php
/**
 *
 *
 *
 */

namespace Ultron\Domain\Service;


use Ultron\Domain\Entity\PerformerEntity;
use Ultron\Domain\Entity\RecordingEntity;
use Zend\Stdlib\ParametersInterface;

interface RecordingServiceInterface
{
    /**
     * Remove a recording
     *
     * @param RecordingEntity $recording
     *
     * @return void
     */
    public function remove(RecordingEntity $recording);

    /**
     * @param ParametersInterface $data
     * @param PerformerEntity     $performer
     *
     * @return void
     */
    public function create(ParametersInterface $data, PerformerEntity $performer);
}
