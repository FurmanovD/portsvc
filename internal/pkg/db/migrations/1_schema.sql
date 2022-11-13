-- +goose Up
CREATE DATABASE IF NOT EXISTS `portsvc` /*!40100 DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci */;

CREATE TABLE IF NOT EXISTS `ports` (
  `id` int(13) NOT NULL AUTO_INCREMENT,
  `portid` varchar(100) COLLATE utf8_unicode_ci NOT NULL,
  `name` varchar(100) COLLATE utf8_unicode_ci NOT NULL,
  `city` varchar(100) COLLATE utf8_unicode_ci NOT NULL,
  `country` varchar(100) COLLATE utf8_unicode_ci NOT NULL,
  `province` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  `timezone` varchar(100) COLLATE utf8_unicode_ci NOT NULL,
  `code` varchar(100) COLLATE utf8_unicode_ci NOT NULL,
  `latitude` varchar(15) COLLATE utf8_unicode_ci NOT NULL,
  `longitude` varchar(15) COLLATE utf8_unicode_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `portid_UNIQUE` (`portid`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

CREATE TABLE IF NOT EXISTS `ports_unlocks` (
  `id` int(13) NOT NULL AUTO_INCREMENT,
  `port_id` int(13) NOT NULL,
  `unlockid` varchar(100) COLLATE utf8_unicode_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_ports_unlocks_ports_idx` (`unlockid`),
  CONSTRAINT `fk_ports_unlocks_ports` FOREIGN KEY (`unlockid`) REFERENCES `ports` (`portid`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

CREATE TABLE IF NOT EXISTS `ports_aliases` (
  `id` int(13) NOT NULL AUTO_INCREMENT,
  `port_id` int(13) NOT NULL,
  `alias` varchar(100) COLLATE utf8_unicode_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_ports_aliases_ports_idx` (`port_id`),
  CONSTRAINT `fk_ports_aliases_ports` FOREIGN KEY (`port_id`) REFERENCES `ports` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- +goose Down
-- DROP TABLE `ports`;
-- DROP TABLE `ports_unlocks`;
-- DROP TABLE `ports_aliases`;
