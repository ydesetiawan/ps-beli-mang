CREATE TABLE merchants (
   id char(26) PRIMARY KEY,
   name VARCHAR(30) NOT NULL,
   merchant_category  VARCHAR(30) NOT NULL,
   image_url VARCHAR(2048) NOT NULL,
   loc_lat FLOAT NOT NULL,
   loc_long FLOAT NOT NULL,
   created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE merchant_items (
    id char(26) PRIMARY KEY,
    merchant_id char(26) NOT NULL ,
    name VARCHAR(30) NOT NULL ,
    category VARCHAR(30),
    image_url VARCHAR(2048),
    price NUMERIC NOT NULL ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (merchant_id) REFERENCES merchants(id)
);
