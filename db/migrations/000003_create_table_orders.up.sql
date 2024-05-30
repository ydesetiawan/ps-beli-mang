CREATE TABLE orders (
    id char(26) PRIMARY KEY,
    user_id char(26) NOT NULL,
    total_price NUMERIC NOT NULL,
    delivery_time INT NOT NULL,
    is_order BOOLEAN NOT NULL default false,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE order_items (
    id char(26) PRIMARY KEY,
    user_id char(26) NOT NULL,
    order_id char(26) NOT NULL ,
    merchant_id char(26) NOT NULL ,
    merchant_item_id char(26) NOT NULL ,
    quantity INT NOT NULL ,
    price NUMERIC NOT NULL ,
    created_at timestamp,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (order_id) REFERENCES orders(id),
    FOREIGN KEY (merchant_id) REFERENCES merchants(id),
    FOREIGN KEY (merchant_item_id) REFERENCES merchant_items(id)
);
