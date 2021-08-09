CREATE TABLE IF NOT EXISTS `books` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `isbn` text NOT NULL,
  `title` varchar(45) NOT NULL,
  `authors` varchar(45) NOT NULL,
  `imageUrl` text,
  `smallImageUrl` text,
  `publicationYear` int(4) NOT NULL,
  `publisher` text,
  `userId` varchar(45) NOT NULL,
  `status` int(11) NOT NULL,
  `description` text,
  `pageCount` int(11) DEFAULT NULL,
  `categories` text,
  `language` text,
  `source` text NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
