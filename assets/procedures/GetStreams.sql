/* GetStreams */
USE `feeder`;
DROP PROCEDURE IF EXISTS `GetStreams`;

DELIMITER //
CREATE PROCEDURE `feeder` . `GetStreams` ()
 BEGIN
   SELECT id, name, streams.desc, streams.host_id, owner, type, head, proto, domain, url, url_end, stream_hosts.interval FROM cams.streams JOIN cams.stream_hosts ON streams.Host_ID=stream_hosts.Host_ID;
 END //
DELIMITER ;
