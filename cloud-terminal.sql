CREATE TABLE UserGroups (
                            GroupID INT PRIMARY KEY,
                            GroupName VARCHAR(100),
    -- 其他用户组信息字段
);

CREATE TABLE Users (
                       UserID INT PRIMARY KEY,
                       UserName VARCHAR(100),
                       Password VARCHAR(100),
                       Email VARCHAR(100),
                       Nickname VARCHAR(100),
                       TOTP_Secret VARCHAR(100),
                       Online BOOLEAN,
                       EnableType ENUM('Enabled', 'Disabled'),
                       GroupID INT,
                       UserType ENUM('Admin', 'Auditor', 'SuperUser', 'User'),
                       FOREIGN KEY(GroupID) REFERENCES UserGroups(GroupID)
    -- 其他用户信息字段
);


CREATE TABLE AssetGroups (
                             GroupID INT PRIMARY KEY,
                             GroupName VARCHAR(100),
    -- 其他资产组信息字段
);

CREATE TABLE Assets (
                        AssetID INT PRIMARY KEY,
                        AssetName VARCHAR(100),
                        GroupID INT,
                        FOREIGN KEY(GroupID) REFERENCES AssetGroups(GroupID)
    -- 其他资产信息字段
);

CREATE TABLE UserAssets (
                            UserID INT,
                            AssetID INT,
                            PRIMARY KEY(UserID, AssetID),
                            FOREIGN KEY(UserID) REFERENCES Users(UserID),
                            FOREIGN KEY(AssetID) REFERENCES Assets(AssetID)
);
