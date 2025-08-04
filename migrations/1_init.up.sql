-- users
CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     name TEXT NOT NULL,
                                     email TEXT NOT NULL UNIQUE,
                                     password_hash TEXT NOT NULL,
                                     role TEXT NOT NULL DEFAULT 'user',  -- 'user', 'seller', 'admin'
                                     created_at TIMESTAMP NOT NULL DEFAULT now()
    );

-- categories
CREATE TABLE IF NOT EXISTS categories (
                                          id SERIAL PRIMARY KEY,
                                          name TEXT NOT NULL UNIQUE,
                                          description TEXT
);

-- products
CREATE TABLE IF NOT EXISTS products (
                                        id SERIAL PRIMARY KEY,
                                        seller_id INT NOT NULL
                                        REFERENCES users(id) ON DELETE CASCADE,
    category_id INT
    REFERENCES categories(id) ON DELETE SET NULL,
    name TEXT NOT NULL,
    description TEXT,
    price NUMERIC(10,2) NOT NULL CHECK (price >= 0),
    stock INT NOT NULL DEFAULT 0 CHECK (stock >= 0),
    created_at TIMESTAMP NOT NULL DEFAULT now()
    );

-- carts
CREATE TABLE IF NOT EXISTS carts (
                                     id SERIAL PRIMARY KEY,
                                     user_id INT NOT NULL
                                     REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT now()
    );

-- cart_items
CREATE TABLE IF NOT EXISTS cart_items (
                                          id SERIAL PRIMARY KEY,
                                          cart_id INT NOT NULL
                                          REFERENCES carts(id) ON DELETE CASCADE,
    product_id INT NOT NULL
    REFERENCES products(id) ON DELETE CASCADE,
    quantity INT NOT NULL CHECK (quantity > 0),
    UNIQUE (cart_id, product_id)
    );

-- orders
CREATE TABLE IF NOT EXISTS orders (
                                      id SERIAL PRIMARY KEY,
                                      user_id INT NOT NULL
                                      REFERENCES users(id) ON DELETE CASCADE,
    status TEXT NOT NULL DEFAULT 'pending',  -- 'pending','shipped','completed','canceled'
    total NUMERIC(10,2) NOT NULL CHECK (total >= 0),
    created_at TIMESTAMP NOT NULL DEFAULT now()
    );

-- order_items
CREATE TABLE IF NOT EXISTS order_items (
                                           id SERIAL PRIMARY KEY,
                                           order_id INT NOT NULL
                                           REFERENCES orders(id) ON DELETE CASCADE,
    product_id INT NOT NULL
    REFERENCES products(id) ON DELETE CASCADE,
    quantity INT NOT NULL CHECK (quantity > 0),
    price NUMERIC(10,2) NOT NULL CHECK (price >= 0)
    );

-- Индексы для ускорения поиска и фильтрации
CREATE INDEX IF NOT EXISTS idx_products_name
    ON products USING gin (to_tsvector('simple', name));

CREATE INDEX IF NOT EXISTS idx_products_price
    ON products (price);

CREATE INDEX IF NOT EXISTS idx_orders_user_id
    ON orders (user_id);
