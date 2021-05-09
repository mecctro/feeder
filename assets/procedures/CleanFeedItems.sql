/* CleanFeedItems */
USE `feeder`;
DROP PROCEDURE IF EXISTS `CleanFeedItems`;

DELIMITER //
CREATE PROCEDURE `feeder` . `CleanFeedItems` (IN days INT)
 BEGIN
  DELETE FROM feed_items WHERE published < DATE_SUB(CURDATE(),INTERVAL days DAY);
 END//
DELIMITER ;
