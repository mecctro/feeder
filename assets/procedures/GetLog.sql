/* GetLog */
USE `feeder`;
DROP PROCEDURE IF EXISTS `GetLog`;

DELIMITER //
CREATE PROCEDURE `feeder` . `GetLog` (IN thisType VARCHAR(25), IN thisLimit INT)
 BEGIN
  SELECT
    ID,
    Type,
    Message
  FROM feed.log
  WHERE Type = thisType
  ORDER BY ID DESC
  LIMIT thisLimit;
 END //
DELIMITER ;
