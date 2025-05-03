CREATE TABLE Tag (
    id SERIAL PRIMARY KEY,
    name VARCHAR(512) NOT NULL,
    color VARCHAR(7) NOT NULL
);

CREATE TABLE Category (
    id SERIAL PRIMARY KEY,
    name VARCHAR(512) NOT NULL,
    color VARCHAR(7) NOT NULL
);

CREATE TYPE StatusProduct AS ENUM ('bought', 'need_buying', 'not_need_buying');

CREATE TABLE Product (
    id SERIAL PRIMARY KEY,
    name VARCHAR(512) NOT NULL,
    priority INTEGER NOT NULL Default 0,
    status StatusProduct NOT NULL DEFAULT 'need_buying',
    icon VARCHAR(1) NOT NULL DEFAULT '🍴',
    category_id INTEGER NOT NULL,
    FOREIGN KEY (category_id) REFERENCES Category(id)  
);

CREATE Table TagInProduct (
    id SERIAL PRIMARY KEY,
    tag_id INTEGER NOT NULL,
    FOREIGN KEY (tag_id) REFERENCES Tag(id)  
    product_id INTEGER NOT NULL,
    FOREIGN KEY (product_id) REFERENCES Product(id)     
)
