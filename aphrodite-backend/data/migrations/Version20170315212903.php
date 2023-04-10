<?php

namespace Aphrodite\DbMigrations;

use Doctrine\DBAL\Migrations\AbstractMigration;
use Doctrine\DBAL\Schema\Schema;

/**
 * Auto-generated Migration: Please modify to your needs!
 */
class Version20170315212903 extends AbstractMigration
{
    /**
     * @param Schema $schema
     */
    public function up(Schema $schema)
    {
        // this up() migration is auto-generated, please modify it to your needs
        $this->abortIf($this->connection->getDatabasePlatform()->getName() !== 'mysql', 'Migration can only be executed safely on \'mysql\'.');

        $this->addSql('ALTER TABLE data_transfer_job DROP FOREIGN KEY FK_40C02E8B158E0B66');
        $this->addSql('ALTER TABLE data_transfer_job DROP FOREIGN KEY FK_40C02E8B953C1C61');
        $this->addSql('DROP TABLE data_transfer_job');
        $this->addSql('DROP TABLE servers');
        $this->addSql('DROP TABLE settings');
    }

    /**
     * @param Schema $schema
     */
    public function down(Schema $schema)
    {
        // this down() migration is auto-generated, please modify it to your needs
        $this->abortIf($this->connection->getDatabasePlatform()->getName() !== 'mysql', 'Migration can only be executed safely on \'mysql\'.');

        $this->addSql('CREATE TABLE data_transfer_job (id INT AUTO_INCREMENT NOT NULL, target_id VARCHAR(255) DEFAULT NULL COLLATE utf8_unicode_ci, source_id VARCHAR(255) DEFAULT NULL COLLATE utf8_unicode_ci, amount INT NOT NULL, transferred INT DEFAULT NULL, ongoing TINYINT(1) NOT NULL, created_at DATETIME NOT NULL, updated_at DATETIME NOT NULL, INDEX IDX_40C02E8B953C1C61 (source_id), INDEX IDX_40C02E8B158E0B66 (target_id), PRIMARY KEY(id)) DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci ENGINE = InnoDB');
        $this->addSql('CREATE TABLE servers (hostname VARCHAR(255) NOT NULL COLLATE utf8_unicode_ci, internal_ip VARCHAR(255) DEFAULT NULL COLLATE utf8_unicode_ci, available_space INT DEFAULT NULL, total_space INT DEFAULT NULL, enabled TINYINT(1) NOT NULL, provisioned TINYINT(1) NOT NULL, network VARCHAR(255) DEFAULT NULL COLLATE utf8_unicode_ci, updated_at DATETIME NOT NULL, created_at DATETIME NOT NULL, type VARCHAR(255) NOT NULL COLLATE utf8_unicode_ci, UNIQUE INDEX UNIQ_4F8AF5F7E551C011 (hostname), UNIQUE INDEX UNIQ_4F8AF5F7A50560A5 (internal_ip), PRIMARY KEY(hostname)) DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci ENGINE = InnoDB');
        $this->addSql('CREATE TABLE settings (name VARCHAR(255) NOT NULL COLLATE utf8_unicode_ci, value VARCHAR(255) NOT NULL COLLATE utf8_unicode_ci, created_at DATETIME NOT NULL, updated_at DATETIME NOT NULL, UNIQUE INDEX UNIQ_E545A0C55E237E06 (name), PRIMARY KEY(name)) DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci ENGINE = InnoDB');
        $this->addSql('ALTER TABLE data_transfer_job ADD CONSTRAINT FK_40C02E8B158E0B66 FOREIGN KEY (target_id) REFERENCES servers (hostname)');
        $this->addSql('ALTER TABLE data_transfer_job ADD CONSTRAINT FK_40C02E8B953C1C61 FOREIGN KEY (source_id) REFERENCES servers (hostname)');
    }
}
