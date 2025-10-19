-- 1. Включаем расширение pgvector
CREATE EXTENSION IF NOT EXISTS vector;

-- 2. Пользователи (для аутентификации)
CREATE TABLE users (
    id serial PRIMARY KEY,
    email text UNIQUE NOT NULL,
    password_hash text NOT NULL,
    full_name text,
    created_at timestamptz DEFAULT now()
);

-- 3. Счета (Карта)
CREATE TABLE accounts (
    id serial PRIMARY KEY,
    user_id integer UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    card_number_mock text, -- (Например "KZ****1234")
    balance numeric(15, 2) NOT NULL DEFAULT 0.00,
    next_payday date       -- Ключевое поле для советов!
);

-- 4. Транзакции
CREATE TABLE transactions (
    id serial PRIMARY KEY,
    account_id integer REFERENCES accounts(id) ON DELETE CASCADE,
    amount numeric(15, 2) NOT NULL,
    merchant text,       -- ("ИП 'МАГАЗИН'")
    category text,       -- ("Продукты", "Одежда", "Такси")
    type text NOT NULL,  -- 'income' (доход) or 'expense' (расход)
    created_at timestamptz DEFAULT now()
);

-- 5. Финансирование (Халяльная замена "кредитам")
CREATE TABLE financing (
    id serial PRIMARY KEY,
    user_id integer REFERENCES users(id) ON DELETE CASCADE,
    product_name text, -- 'Мурабаха на авто'
    total_amount numeric(15, 2) NOT NULL,
    remaining_amount numeric(15, 2) NOT NULL,
    monthly_payment numeric(15, 2) NOT NULL
);

-- 6. Таблица для RAG (наши знания)
CREATE TABLE product_embeddings (
    id bigserial PRIMARY KEY,
    content text NOT NULL,
    embedding vector(1536) -- 1536 - это размер для text-embedding-3-small
);

-- 7. Таблица для истории чатов (память ассистента)
CREATE TABLE chat_history (
    id bigserial PRIMARY KEY,
    user_id integer REFERENCES users(id) ON DELETE CASCADE,
    session_id text NOT NULL, -- (Например, UUID)
    role text NOT NULL, -- 'user' или 'ai'
    content text NOT NULL,
    created_at timestamptz DEFAULT now()
);

-- (Примерные данные для теста)
-- Вставляем тестового юзера (пароль: 123456)
INSERT INTO users (email, password_hash, full_name) VALUES 
('test@zaman.kz', '$2a$10$fwhg9nE/iKmoNM.yECuAWeMa3mRkSCg/BEs2K2rftP.2/aNA1sThu', 'Тестовый Клиент');

INSERT INTO accounts (user_id, card_number_mock, balance, next_payday) VALUES
(1, 'KZ****1234', 150000.00, (current_date + interval '15 days'));

INSERT INTO financing (user_id, product_name, total_amount, remaining_amount, monthly_payment) VALUES
(1, 'Мурабаха на авто', 5000000.00, 2500000.00, 150000.00);