/* SetLog */
USE `feeder`;
DROP PROCEDURE IF EXISTS `SetLog`;

DELIMITER //
CREATE PROCEDURE `feeder` . `SetLog` (IN thisType VARCHAR(25), IN thisMessage TEXT)
 BEGIN
  INSERT IGNORE INTO feeder.log (Type, Message)
  VALUES (thisType, thisMessage);
 END //
DELIMITER ;
