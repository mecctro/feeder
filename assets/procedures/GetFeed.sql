/* GetFeed */
USE `feeder`;
DROP PROCEDURE IF EXISTS `GetFeed`;

DELIMITER //
CREATE PROCEDURE `feeder` . `GetFeed` (IN thisFeedHash CHAR(40), IN thisFeedType VARCHAR(255))
 BEGIN
  /* Get all feeds of type */
  IF thisFeedType = '0' THEN
    SELECT
      id,
      feedhash,
      feedtype,
      title,
      description,
      link,
      feedlink,
      updated,
      published,
      author,
      language,
      image,
      copyright
    FROM feeds
    ORDER BY title;
  /* Get specific feed by type */
  ELSEIF thisFeedType IN ('Feed', 'Forum', 'Video') THEN
    SELECT
    id,
    feedhash,
    feedtype,
    title,
    description,
    link,
    feedlink,
    updated,
    published,
    author,
    language,
    image,
    copyright
    FROM feeds
    WHERE feedtype = thisFeedType
    ORDER BY title;
  ELSE
    SELECT
    id,
    feedhash,
    feedtype,
    title,
    description,
    link,
    feedlink,
    updated,
    published,
    author,
    language,
    image,
    copyright
    FROM feeds
    WHERE feedhash = thisFeedHash
    ORDER BY title;
  END IF;
 END //
DELIMITER ;
