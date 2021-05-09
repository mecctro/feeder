/* GetCommoditiesExchange */
USE `feeder`;
DROP PROCEDURE IF EXISTS `GetCommoditiesExchange`;

DELIMITER //
CREATE PROCEDURE `feeder` . `GetCommoditiesExchange` (IN thisTimestamp INT)
 BEGIN
  /* Get exchange rate from timestamp til current */
  SELECT
    ID,
    Timestamp,
    Base,
	  Aluminum,
    Bismuth,
    Coal,
    Cobalt,
	  Copper,
    Crude,
    Gold,
    Indium,
	  Iridium,
    Iron,
    Lead,
    Magnesium,
	  Manganese,
    NaturalGas,
    Nickel,
    Palladium,
	  Platinum,
    Rhenium,
    Rhodium,
    Ruthenium,
	  Silver,
    Tantalum,
    Tin,
    Titanium,
	  Tungsten,
    Zinc
  FROM exchange_rates_commodities
	WHERE Timestamp >= thisTimestamp
  ORDER BY Timestamp;
 END //
DELIMITER ;
