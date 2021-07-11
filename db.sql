CREATE TABLE users (
     userID int,
     uniqueID varchar(255),
     ethosKey varchar(25),
     keyExpiry datetime,
     discordID varchar(255)
);

INSERT INTO users 
VALUES (userID, uniqueID, ethosKey, keyExpiry, discordID);

DELETE FROM table_name WHERE condition; 

DATE_ADD(CURRENT_TIMESTAMP, INTERVAL 30 DAY);



CREATE TABLE users (
     ethosKey varchar(30) NOT NULL UNIQUE,
     userId varchar(255) NOT NULL UNIQUE,
     expiration datetime NOT NULL,
     PRIMARY KEY (ethosKey)
);

CREATE TABLE sigma (
     sigmaKey varchar(40) NOT NULL UNIQUE,
     ethosKey varchar(30) NOT NULL UNIQUE,
     PRIMARY KEY (sigmaKey),
     FOREIGN KEY (ethosKey) REFERENCES users(ethosKey)

);