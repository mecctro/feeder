/* GetLinks */
USE `feeder`;
DROP PROCEDURE IF EXISTS `GetLinks`;

DELIMITER //
CREATE PROCEDURE `feeder` . `GetLinks` ()
 BEGIN
   SELECT id, title, url FROM cams.links ORDER BY title;
 END //
DELIMITER ;
