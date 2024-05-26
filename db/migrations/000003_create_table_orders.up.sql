CREATE TABLE orders (
    id char(26) PRIMARY KEY,
    total_price NUMERIC NOT NULL,
    delivery_time INT NOT NULL,
    is_order BOOLEAN NOT NULL default false,
    user_loc_lat FLOAT NOT NULL,
    user_loc_long FLOAT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE order_items (
    id char(26) PRIMARY KEY,
    order_id char(26) NOT NULL ,
    merchant_id char(26) NOT NULL ,
    is_starting_point BOOLEAN NOT NULL  default false,
    merchant_item_id char(26) NOT NULL ,
    quantity INT NOT NULL ,
    price NUMERIC NOT NULL ,
    amount NUMERIC NOT NULL ,
    created_at timestamp,
    FOREIGN KEY (order_id) REFERENCES orders(id),
    FOREIGN KEY (merchant_id) REFERENCES merchants(id),
    FOREIGN KEY (merchant_item_id) REFERENCES merchant_items(id)
);
