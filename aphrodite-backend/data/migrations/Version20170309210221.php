<?php

namespace Aphrodite\DbMigrations;

use Doctrine\DBAL\Migrations\AbstractMigration;
use Doctrine\DBAL\Schema\Schema;

/**
 * Auto-generated Migration: Please modify to your needs!
 */
class Version20170309210221 extends AbstractMigration
{
    /**
     * @param Schema $schema
     */
    public function up(Schema $schema)
    {
        // this up() migration is auto-generated, please modify it to your needs
        $this->abortIf($this->connection->getDatabasePlatform()->getName() !== 'mysql', 'Migration can only be executed safely on \'mysql\'.');

        $this->addSql('DROP TABLE recording_images');
        $this->addSql('DROP INDEX type ON recordings');
        $this->addSql('DROP INDEX storageServer ON recordings');
        $this->addSql('DROP INDEX storagePath ON recordings');
        $this->addSql('DROP INDEX location ON recordings');
        $this->addSql('DROP INDEX state ON recordings');
        $this->addSql('DROP INDEX state_within_period ON recordings');
        $this->addSql('DROP INDEX videoUrl ON recordings');
        $this->addSql('DROP INDEX createdAt ON recordings');
        $this->addSql('DROP INDEX updatedAt ON recordings');
        $this->addSql('DROP INDEX lastCheckedAt ON recordings');
        $this->addSql('ALTER TABLE recordings ADD video_uuid VARCHAR(255) DEFAULT NULL, ADD wordpress_collage_uuid VARCHAR(255) DEFAULT NULL, ADD infinity_collage_uuid VARCHAR(255) DEFAULT NULL, ADD sprites LONGTEXT DEFAULT NULL COMMENT \'(DC2Type:json_array)\', ADD images LONGTEXT DEFAULT NULL COMMENT \'(DC2Type:json_array)\', DROP storage_server, DROP storage_path, DROP storage_path_thumb, DROP views, DROP type, DROP storage_path_collage');
        $this->addSql('CREATE INDEX state ON recordings (state)');
        $this->addSql('CREATE INDEX state_within_period ON recordings (state, created_at)');
        $this->addSql('CREATE INDEX videoUrl ON recordings (video_url)');
        $this->addSql('CREATE INDEX createdAt ON recordings (created_at)');
        $this->addSql('CREATE INDEX updatedAt ON recordings (updated_at)');
        $this->addSql('CREATE INDEX lastCheckedAt ON recordings (last_checked_at)');
    }

    /**
     * @param Schema $schema
     */
    public function down(Schema $schema)
    {
        // this down() migration is auto-generated, please modify it to your needs
        $this->abortIf($this->connection->getDatabasePlatform()->getName() !== 'mysql', 'Migration can only be executed safely on \'mysql\'.');

        $this->addSql('CREATE TABLE recording_images (id INT AUTO_INCREMENT NOT NULL, recording_id INT DEFAULT NULL, storage_server VARCHAR(255) NOT NULL COLLATE utf8_unicode_ci, storage_path VARCHAR(255) NOT NULL COLLATE utf8_unicode_ci, captured_at DOUBLE PRECISION NOT NULL, INDEX IDX_D3BA23B68CA9A845 (recording_id), PRIMARY KEY(id)) DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci ENGINE = InnoDB');
        $this->addSql('ALTER TABLE recording_images ADD CONSTRAINT FK_D3BA23B68CA9A845 FOREIGN KEY (recording_id) REFERENCES recordings (id)');
        $this->addSql('DROP INDEX state ON recordings');
        $this->addSql('DROP INDEX state_within_period ON recordings');
        $this->addSql('DROP INDEX videoUrl ON recordings');
        $this->addSql('DROP INDEX createdAt ON recordings');
        $this->addSql('DROP INDEX updatedAt ON recordings');
        $this->addSql('DROP INDEX lastCheckedAt ON recordings');
        $this->addSql('ALTER TABLE recordings ADD storage_server VARCHAR(255) DEFAULT NULL COLLATE utf8_unicode_ci, ADD storage_path VARCHAR(255) DEFAULT NULL COLLATE utf8_unicode_ci, ADD storage_path_thumb VARCHAR(255) DEFAULT NULL COLLATE utf8_unicode_ci, ADD views INT NOT NULL, ADD type VARCHAR(255) NOT NULL COLLATE utf8_unicode_ci, ADD storage_path_collage VARCHAR(255) DEFAULT NULL COLLATE utf8_unicode_ci, DROP video_uuid, DROP wordpress_collage_uuid, DROP infinity_collage_uuid, DROP sprites, DROP images');
        $this->addSql('CREATE INDEX type ON recordings (id, type)');
        $this->addSql('CREATE INDEX storageServer ON recordings (storage_server)');
        $this->addSql('CREATE INDEX storagePath ON recordings (storage_path)');
        $this->addSql('CREATE INDEX location ON recordings (storage_server, storage_path, type)');
        $this->addSql('CREATE INDEX state ON recordings (id, state, type)');
        $this->addSql('CREATE INDEX state_within_period ON recordings (id, state, created_at, type)');
        $this->addSql('CREATE INDEX videoUrl ON recordings (video_url, type)');
        $this->addSql('CREATE INDEX createdAt ON recordings (created_at, type)');
        $this->addSql('CREATE INDEX updatedAt ON recordings (updated_at, type)');
        $this->addSql('CREATE INDEX lastCheckedAt ON recordings (last_checked_at, type)');
    }
}
