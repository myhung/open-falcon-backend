set names utf8;

CREATE TABLE IF NOT EXISTS `placard` (
  id INT NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (id),
  title varchar(100) NOT NULL,
  content varchar(1000) NOT NULL,
  createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
  updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  authorId int(10) unsigned NOT NULL,
  FOREIGN KEY (authorId)
  REFERENCES uic.user(id)
  ON UPDATE RESTRICT ON DELETE CASCADE
) ENGINE=INNODB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `readpath` (
  id int(10) unsigned NOT NULL,
  PRIMARY KEY(id),
  readAt DATETIME ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (id)
  REFERENCES uic.user(id)
  ON UPDATE RESTRICT ON DELETE CASCADE
) ENGINE=INNODB DEFAULT CHARSET=utf8;
