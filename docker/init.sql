CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    telegram_id TEXT UNIQUE NOT NULL
);

CREATE TABLE groups (
    group_id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE group_members (
    user_id INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    group_id INT NOT NULL REFERENCES groups(group_id) ON DELETE CASCADE,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (user_id, group_id)
);


CREATE TABLE shopping_items (
    group_id INT NOT NULL REFERENCES groups(group_id) ON DELETE CASCADE,
    item TEXT NOT NULL
);

