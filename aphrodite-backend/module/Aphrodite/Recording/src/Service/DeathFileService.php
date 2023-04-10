<?php
/**
 *
 *
 *  AB
 */

namespace Aphrodite\Recording\Service;

use Aphrodite\Recording\Entity\DeathFileEntity;
use BsbFlysystem\Filter\File\RenameUpload;
use DateTime;
use Doctrine\Common\Persistence\ObjectManager;
use League\Flysystem\FileNotFoundException;
use League\Flysystem\Filesystem;
use Rhubarb\Rhubarb;
use Zend\Stdlib\ParametersInterface;

class DeathFileService implements DeathFileServiceInterface
{
    /**
     * @var ObjectManager
     */
    private $objectManager;

    /**
     * @var Filesystem
     */
    private $filesystem;

    /**
     * @var Rhubarb
     */
    private $rhubarb;

    /**
     * DeathFileService constructor.
     *
     * @param ObjectManager $objectManager
     * @param Filesystem    $filesystem
     * @param Rhubarb       $rhubarb
     */
    public function __construct(ObjectManager $objectManager, Filesystem $filesystem, Rhubarb $rhubarb)
    {
        $this->rhubarb       = $rhubarb;
        $this->filesystem    = $filesystem;
        $this->objectManager = $objectManager;
    }

    /**
     * {@inheritdoc}
     */
    public function placeInQueue(DeathFileEntity $deathFile)
    {
        $task = $this->rhubarb->sendTask('central.tasks.process_death_file', [$deathFile->getId()]);
        $task->delay(['queue' => 'death_file']);
    }

    /**
     * {@inheritdoc}
     */
    public function createFromFile(ParametersInterface $file)
    {
        $rename = new RenameUpload([
            'useUploadExtension' => true,
            'target'             => 'death-files/',
            'randomize'          => true,
            'filesystem'         => $this->filesystem,
        ]);

        $target = $rename->filter($file->get('tmp_name'));

        $entry = new DeathFileEntity();
        $entry->setLocation($target);
        $entry->setCreatedAt(new DateTime());
        $entry->setUpdatedAt(new DateTime());

        $this->objectManager->persist($entry);
        $this->objectManager->flush();

        // Must be done after flushing the object manager
        $this->placeInQueue($entry);

        return $entry;
    }

    /**
     * {@inheritdoc}
     */
    public function update(DeathFileEntity $file)
    {
        $file->setUpdatedAt(new DateTime());
        $this->objectManager->flush();
    }

    /**
     * {@inheritdoc}
     */
    public function delete(DeathFileEntity $file)
    {
        try {
            // Remove the file from the filesystem
            $this->filesystem->delete($file->getLocation());
        } catch (FileNotFoundException $e) {
            // we don't care if it didn't exist, we wanted it removed anyway
        }

        $this->objectManager->remove($file);
        $this->objectManager->flush();
    }
}
