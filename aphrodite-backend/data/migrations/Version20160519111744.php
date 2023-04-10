<?php

namespace Aphrodite\DbMigrations;

use Doctrine\DBAL\Migrations\AbstractMigration;
use Doctrine\DBAL\Schema\Schema;

/**
 * Auto-generated Migration: Please modify to your needs!
 */
class Version20160519111744 extends AbstractMigration
{
    /**
     * @param Schema $schema
     */
    public function up(Schema $schema)
    {
        // this up() migration is auto-generated, please modify it to your needs
        $this->abortIf($this->connection->getDatabasePlatform()->getName() != 'mysql', 'Migration can only be executed safely on \'mysql\'.');

        $this->addSql('CREATE TABLE oauth_access_tokens (token VARCHAR(40) NOT NULL, client_id VARCHAR(60) DEFAULT NULL, owner_id INT DEFAULT NULL, scopes LONGTEXT NOT NULL COMMENT \'(DC2Type:json_array)\', expires_at DATETIME DEFAULT NULL, INDEX IDX_CA42527C19EB6921 (client_id), INDEX IDX_CA42527C7E3C61F9 (owner_id), PRIMARY KEY(token)) DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci ENGINE = InnoDB');
        $this->addSql('CREATE TABLE oauth_authorization_codes (token VARCHAR(40) NOT NULL, client_id VARCHAR(60) DEFAULT NULL, owner_id INT DEFAULT NULL, expires_at DATETIME NOT NULL, scopes LONGTEXT NOT NULL COMMENT \'(DC2Type:json_array)\', redirect_uri VARCHAR(1000) NOT NULL, INDEX IDX_98A471C419EB6921 (client_id), INDEX IDX_98A471C47E3C61F9 (owner_id), PRIMARY KEY(token)) DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci ENGINE = InnoDB');
        $this->addSql('CREATE TABLE oauth_clients (id VARCHAR(60) NOT NULL, secret VARCHAR(60) NOT NULL, redirect_uris LONGTEXT NOT NULL COMMENT \'(DC2Type:json_array)\', PRIMARY KEY(id)) DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci ENGINE = InnoDB');
        $this->addSql('CREATE TABLE oauth_refresh_tokens (token VARCHAR(40) NOT NULL, client_id VARCHAR(60) DEFAULT NULL, owner_id INT DEFAULT NULL, scopes LONGTEXT NOT NULL COMMENT \'(DC2Type:json_array)\', expires_at DATETIME DEFAULT NULL, INDEX IDX_5AB68719EB6921 (client_id), INDEX IDX_5AB6877E3C61F9 (owner_id), PRIMARY KEY(token)) DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci ENGINE = InnoDB');
        $this->addSql('CREATE TABLE oauth_scopes (id INT AUTO_INCREMENT NOT NULL, name VARCHAR(40) NOT NULL, description VARCHAR(500) NOT NULL, is_default TINYINT(1) NOT NULL, PRIMARY KEY(id)) DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci ENGINE = InnoDB');
        $this->addSql('CREATE TABLE performers (id INT AUTO_INCREMENT NOT NULL, service_id VARCHAR(255) NOT NULL, stage_name VARCHAR(255) NOT NULL, blacklisted TINYINT(1) NOT NULL, online TINYINT(1) NOT NULL, is_recording TINYINT(1) NOT NULL, is_pending_recording TINYINT(1) NOT NULL, recording_count INT NOT NULL, current_viewers INT NOT NULL, peak_viewer_count INT NOT NULL, aliases LONGTEXT DEFAULT NULL COMMENT \'(DC2Type:simple_array)\', updated_at DATETIME NOT NULL, created_at DATETIME NOT NULL, service VARCHAR(255) NOT NULL, video_state SMALLINT DEFAULT NULL, cam_score SMALLINT DEFAULT NULL, cam_server SMALLINT DEFAULT NULL, miss_mfc_rank SMALLINT DEFAULT NULL, access_level SMALLINT DEFAULT NULL, section VARCHAR(255) DEFAULT NULL, INDEX serviceId (service_id, service, id), INDEX section (section), INDEX service (service), INDEX online (online, service, id), INDEX blacklisted (blacklisted, service), INDEX recording (is_recording, service, id), INDEX pending_recording (is_pending_recording, service, id), INDEX stageName (stage_name), INDEX createdAt (created_at), INDEX updatedAt (updated_at), UNIQUE INDEX service_and_id_unq (service_id, service), PRIMARY KEY(id)) DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci ENGINE = InnoDB');
        $this->addSql('CREATE TABLE recordings (id INT AUTO_INCREMENT NOT NULL, performer_id INT DEFAULT NULL, state VARCHAR(255) NOT NULL, stage_name VARCHAR(255) NOT NULL, section VARCHAR(255) DEFAULT NULL, service VARCHAR(255) DEFAULT NULL, duration INT NOT NULL, size INT NOT NULL, description VARCHAR(255) DEFAULT NULL, storage_server VARCHAR(255) DEFAULT NULL, storage_path VARCHAR(255) DEFAULT NULL, storage_path_thumb VARCHAR(255) DEFAULT NULL, bit_rate VARCHAR(255) DEFAULT NULL, encoding VARCHAR(255) DEFAULT NULL, audio VARCHAR(255) DEFAULT NULL, video VARCHAR(255) DEFAULT NULL, video_url VARCHAR(255) DEFAULT NULL COLLATE utf8_bin, video_url_valid TINYINT(1) NOT NULL, image_urls_valid TINYINT(1) NOT NULL, views INT NOT NULL, orphaned INT NOT NULL, created_at DATETIME NOT NULL, updated_at DATETIME NOT NULL, last_checked_at DATETIME DEFAULT NULL, image_urls_thumb VARCHAR(255) DEFAULT NULL, image_urls_large VARCHAR(255) DEFAULT NULL, image_urls_gallery VARCHAR(255) DEFAULT NULL, type VARCHAR(255) NOT NULL, INDEX IDX_E9D79C6E6C6B33F3 (performer_id), INDEX type (id, type), INDEX state (id, state, type), INDEX state_within_period (id, state, created_at, type), INDEX service (service), INDEX stageName (stage_name), INDEX storageServer (storage_server), INDEX storagePath (storage_path), INDEX videoUrl (video_url, type), INDEX location (storage_server, storage_path, type), INDEX createdAt (created_at, type), INDEX updatedAt (updated_at, type), INDEX lastCheckedAt (last_checked_at, type), PRIMARY KEY(id)) DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci ENGINE = InnoDB');
        $this->addSql('CREATE TABLE death_file_entries (id INT AUTO_INCREMENT NOT NULL, recording_id INT DEFAULT NULL, death_file_id INT DEFAULT NULL, url VARCHAR(255) NOT NULL, hermes_id INT DEFAULT NULL, state VARCHAR(255) DEFAULT NULL, created_at DATETIME NOT NULL, updated_at DATETIME NOT NULL, UNIQUE INDEX UNIQ_8B0F778AF47645AE (url), UNIQUE INDEX UNIQ_8B0F778A8CA9A845 (recording_id), INDEX IDX_8B0F778A7A2DE58C (death_file_id), PRIMARY KEY(id)) DEFAULT CHARACTER SET utf8 COLLATE utf8_bin ENGINE = InnoDB');
        $this->addSql('CREATE TABLE death_file_entity (id INT AUTO_INCREMENT NOT NULL, location VARCHAR(255) NOT NULL, entries INT DEFAULT NULL, ignored INT NOT NULL, pending INT NOT NULL, created_at DATETIME NOT NULL, updated_at DATETIME NOT NULL, UNIQUE INDEX UNIQ_45BB6CCD5E9E89CB (location), PRIMARY KEY(id)) DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci ENGINE = InnoDB');
        $this->addSql('CREATE TABLE servers (hostname VARCHAR(255) NOT NULL, internal_ip VARCHAR(255) DEFAULT NULL, available_space INT DEFAULT NULL, total_space INT DEFAULT NULL, enabled TINYINT(1) NOT NULL, provisioned TINYINT(1) NOT NULL, network VARCHAR(255) DEFAULT NULL, updated_at DATETIME NOT NULL, created_at DATETIME NOT NULL, type VARCHAR(255) NOT NULL, UNIQUE INDEX UNIQ_4F8AF5F7E551C011 (hostname), UNIQUE INDEX UNIQ_4F8AF5F7A50560A5 (internal_ip), PRIMARY KEY(hostname)) DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci ENGINE = InnoDB');
        $this->addSql('CREATE TABLE data_transfer_job (id INT AUTO_INCREMENT NOT NULL, source_id VARCHAR(255) DEFAULT NULL, target_id VARCHAR(255) DEFAULT NULL, amount INT NOT NULL, transferred INT DEFAULT NULL, ongoing TINYINT(1) NOT NULL, created_at DATETIME NOT NULL, updated_at DATETIME NOT NULL, INDEX IDX_40C02E8B953C1C61 (source_id), INDEX IDX_40C02E8B158E0B66 (target_id), PRIMARY KEY(id)) DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci ENGINE = InnoDB');
        $this->addSql('CREATE TABLE settings (name VARCHAR(255) NOT NULL, value VARCHAR(255) NOT NULL, created_at DATETIME NOT NULL, updated_at DATETIME NOT NULL, UNIQUE INDEX UNIQ_E545A0C55E237E06 (name), PRIMARY KEY(name)) DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci ENGINE = InnoDB');
        $this->addSql('CREATE TABLE post_association (id INT AUTO_INCREMENT NOT NULL, recording_id INT DEFAULT NULL, site_id INT DEFAULT NULL, post_id INT NOT NULL, INDEX IDX_11C5ABE38CA9A845 (recording_id), INDEX IDX_11C5ABE3F6BD1646 (site_id), PRIMARY KEY(id)) DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci ENGINE = InnoDB');
        $this->addSql('CREATE TABLE sites (id INT AUTO_INCREMENT NOT NULL, enabled TINYINT(1) NOT NULL, name VARCHAR(255) NOT NULL, api_uri VARCHAR(255) NOT NULL, username VARCHAR(255) NOT NULL, password VARCHAR(255) NOT NULL, sources LONGTEXT NOT NULL COMMENT \'(DC2Type:json_array)\', created_at DATETIME NOT NULL, updated_at DATETIME NOT NULL, PRIMARY KEY(id)) DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci ENGINE = InnoDB');
        $this->addSql('CREATE TABLE users (id INT AUTO_INCREMENT NOT NULL, email VARCHAR(255) NOT NULL, name VARCHAR(255) NOT NULL, surname VARCHAR(255) NOT NULL, password VARCHAR(255) NOT NULL, activated TINYINT(1) NOT NULL, blocked TINYINT(1) NOT NULL, created_at DATETIME NOT NULL, updated_at DATETIME NOT NULL, UNIQUE INDEX UNIQ_1483A5E9E7927C74 (email), PRIMARY KEY(id)) DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci ENGINE = InnoDB');
        $this->addSql('ALTER TABLE oauth_access_tokens ADD CONSTRAINT FK_CA42527C19EB6921 FOREIGN KEY (client_id) REFERENCES oauth_clients (id)');
        $this->addSql('ALTER TABLE oauth_access_tokens ADD CONSTRAINT FK_CA42527C7E3C61F9 FOREIGN KEY (owner_id) REFERENCES users (id) ON DELETE CASCADE');
        $this->addSql('ALTER TABLE oauth_authorization_codes ADD CONSTRAINT FK_98A471C419EB6921 FOREIGN KEY (client_id) REFERENCES oauth_clients (id)');
        $this->addSql('ALTER TABLE oauth_authorization_codes ADD CONSTRAINT FK_98A471C47E3C61F9 FOREIGN KEY (owner_id) REFERENCES users (id) ON DELETE CASCADE');
        $this->addSql('ALTER TABLE oauth_refresh_tokens ADD CONSTRAINT FK_5AB68719EB6921 FOREIGN KEY (client_id) REFERENCES oauth_clients (id)');
        $this->addSql('ALTER TABLE oauth_refresh_tokens ADD CONSTRAINT FK_5AB6877E3C61F9 FOREIGN KEY (owner_id) REFERENCES users (id) ON DELETE CASCADE');
        $this->addSql('ALTER TABLE recordings ADD CONSTRAINT FK_E9D79C6E6C6B33F3 FOREIGN KEY (performer_id) REFERENCES performers (id)');
        $this->addSql('ALTER TABLE death_file_entries ADD CONSTRAINT FK_8B0F778A8CA9A845 FOREIGN KEY (recording_id) REFERENCES recordings (id) ON DELETE CASCADE');
        $this->addSql('ALTER TABLE death_file_entries ADD CONSTRAINT FK_8B0F778A7A2DE58C FOREIGN KEY (death_file_id) REFERENCES death_file_entity (id) ON DELETE CASCADE');
        $this->addSql('ALTER TABLE data_transfer_job ADD CONSTRAINT FK_40C02E8B953C1C61 FOREIGN KEY (source_id) REFERENCES servers (hostname)');
        $this->addSql('ALTER TABLE data_transfer_job ADD CONSTRAINT FK_40C02E8B158E0B66 FOREIGN KEY (target_id) REFERENCES servers (hostname)');
        $this->addSql('ALTER TABLE post_association ADD CONSTRAINT FK_11C5ABE38CA9A845 FOREIGN KEY (recording_id) REFERENCES recordings (id) ON DELETE CASCADE');
        $this->addSql('ALTER TABLE post_association ADD CONSTRAINT FK_11C5ABE3F6BD1646 FOREIGN KEY (site_id) REFERENCES sites (id) ON DELETE CASCADE');
    }

    /**
     * @param Schema $schema
     */
    public function down(Schema $schema)
    {
        // this down() migration is auto-generated, please modify it to your needs
        $this->abortIf($this->connection->getDatabasePlatform()->getName() != 'mysql', 'Migration can only be executed safely on \'mysql\'.');

        $this->addSql('ALTER TABLE oauth_access_tokens DROP FOREIGN KEY FK_CA42527C19EB6921');
        $this->addSql('ALTER TABLE oauth_authorization_codes DROP FOREIGN KEY FK_98A471C419EB6921');
        $this->addSql('ALTER TABLE oauth_refresh_tokens DROP FOREIGN KEY FK_5AB68719EB6921');
        $this->addSql('ALTER TABLE recordings DROP FOREIGN KEY FK_E9D79C6E6C6B33F3');
        $this->addSql('ALTER TABLE death_file_entries DROP FOREIGN KEY FK_8B0F778A8CA9A845');
        $this->addSql('ALTER TABLE post_association DROP FOREIGN KEY FK_11C5ABE38CA9A845');
        $this->addSql('ALTER TABLE death_file_entries DROP FOREIGN KEY FK_8B0F778A7A2DE58C');
        $this->addSql('ALTER TABLE data_transfer_job DROP FOREIGN KEY FK_40C02E8B953C1C61');
        $this->addSql('ALTER TABLE data_transfer_job DROP FOREIGN KEY FK_40C02E8B158E0B66');
        $this->addSql('ALTER TABLE post_association DROP FOREIGN KEY FK_11C5ABE3F6BD1646');
        $this->addSql('ALTER TABLE oauth_access_tokens DROP FOREIGN KEY FK_CA42527C7E3C61F9');
        $this->addSql('ALTER TABLE oauth_authorization_codes DROP FOREIGN KEY FK_98A471C47E3C61F9');
        $this->addSql('ALTER TABLE oauth_refresh_tokens DROP FOREIGN KEY FK_5AB6877E3C61F9');
        $this->addSql('DROP TABLE oauth_access_tokens');
        $this->addSql('DROP TABLE oauth_authorization_codes');
        $this->addSql('DROP TABLE oauth_clients');
        $this->addSql('DROP TABLE oauth_refresh_tokens');
        $this->addSql('DROP TABLE oauth_scopes');
        $this->addSql('DROP TABLE performers');
        $this->addSql('DROP TABLE recordings');
        $this->addSql('DROP TABLE death_file_entries');
        $this->addSql('DROP TABLE death_file_entity');
        $this->addSql('DROP TABLE servers');
        $this->addSql('DROP TABLE data_transfer_job');
        $this->addSql('DROP TABLE settings');
        $this->addSql('DROP TABLE post_association');
        $this->addSql('DROP TABLE sites');
        $this->addSql('DROP TABLE users');
    }
}
