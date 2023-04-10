<?php
/**
 *
 *         
 *
 */

declare(strict_types = 1);

namespace Aphrodite\DbMigrations;

use Doctrine\DBAL\Migrations\AbstractMigration;
use Doctrine\DBAL\Schema\Schema;

/**
 * Auto-generated Migration: Please modify to your needs!
 */
class Version20161004092815 extends AbstractMigration
{
    /**
     * @param Schema $schema
     */
    public function up(Schema $schema)
    {
        $this->addSql('CREATE TABLE recording_images (id INT AUTO_INCREMENT NOT NULL, recording_id INT DEFAULT NULL, storage_server VARCHAR(255) NOT NULL, storage_path VARCHAR(255) NOT NULL, captured_at FLOAT NOT NULL, INDEX IDX_D3BA23B68CA9A845 (recording_id), PRIMARY KEY(id)) DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci ENGINE = InnoDB');
        $this->addSql('ALTER TABLE recording_images ADD CONSTRAINT FK_D3BA23B68CA9A845 FOREIGN KEY (recording_id) REFERENCES recordings (id)');
    }

    /**
     * @param Schema $schema
     */
    public function down(Schema $schema)
    {
        $this->addSql('DROP TABLE recording_images');
    }
}
