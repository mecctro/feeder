/* GetSearch */
USE `feeder`;
DROP PROCEDURE IF EXISTS `GetSearch`;

DELIMITER //
CREATE PROCEDURE `feeder` . `GetSearch` (IN terms TEXT, IN thisLimit INT)
 BEGIN
   SELECT
    feeditemhash,
    feeds.FeedType,
    feeds.title as feedtitle,
    feed_items.title,
    Content,
    Categories,
    feed_items.description,
    feeds.FeedLink as feedlink,
    feed_items.link,
    feed_items.updated,
    feed_items.published,
    feeds.author as feedauthor,
    feed_items.author,
    language,
    copyright,
    MATCH(feed_items.title) AGAINST(terms IN BOOLEAN MODE) AS FeedTitleScore,
    MATCH(Content) AGAINST(terms IN BOOLEAN MODE) AS FeedContentScore,
    MATCH(feed_items.description) AGAINST(terms IN BOOLEAN MODE) AS FeedDescriptionScore
   FROM feed_items
   INNER JOIN feeds ON feeds.ID=feed_items.FeedID
   WHERE MATCH(feed_items.title) AGAINST(terms IN BOOLEAN MODE)
   ORDER BY FeedTitleScore DESC, FeedContentScore DESC, FeedDescriptionScore DESC, feed_items.Published DESC
   LIMIT thisLimit;
 END //
DELIMITER ;
