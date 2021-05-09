/* GetFeedItems */
USE `feeder`;
DROP PROCEDURE IF EXISTS `GetFeedItems`;

DELIMITER //
CREATE PROCEDURE `feeder` . `GetFeedItems` (IN thisFeedID INT, IN thisfeedtype VARCHAR(255), IN thislimit INT, IN thisoffset INT)
 BEGIN
  CASE thisFeedID
    /* Get all feeds of type */
    WHEN 0 THEN
      SELECT
        feed_items.id,
        feeditemhash,
        feeds.id as feedid,
        feeds.title as feedtitle,
        feed_items.title,
        Content,
        Categories,
        feed_items.description,
        feeds.link as feedlink,
        feed_items.link,
        feed_items.updated,
        feed_items.published,
        feeds.author as feedauthor,
        feed_items.author, language,
        feeds.image as feedimage,
        feed_items.image,
        copyright
      FROM feed_items
      INNER JOIN feeds ON feeds.ID=feed_items.FeedID
      WHERE feedtype = thisfeedtype
      ORDER BY feed_items.published
      DESC LIMIT thislimit OFFSET thisoffset;
    /* Get specific feed by type */
    ELSE
      SELECT
        feed_items.id,
        feeditemhash,
        feeds.id as feedid,
        feeds.title as feedtitle,
        feed_items.title,
        Content,
        Categories,
        feed_items.description,
        feeds.link as feedlink,
        feed_items.link,
        feed_items.updated,
        feed_items.published,
        feeds.author as feedauthor,
        feed_items.author,
        language,
        feeds.image as feedimage,
        feed_items.image,
        copyright
      FROM feed_items
      INNER JOIN feeds ON feeds.ID=feed_items.FeedID
      WHERE feeds.ID=thisFeedID AND feedtype = thisfeedtype
      ORDER BY feed_items.published
      DESC LIMIT thislimit OFFSET thisoffset;
  END CASE;
 END //
DELIMITER ;
