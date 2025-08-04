-- Удаляем индексы
DROP INDEX IF EXISTS idx_orders_user_id;
DROP INDEX IF EXISTS idx_products_price;
DROP INDEX IF EXISTS idx_products_name;

-- Удаляем таблицы в обратном порядке зависимости
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS cart_items;
DROP TABLE IF EXISTS carts;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS users;
