INSERT INTO walk_options (id, category, title)
VALUES
  (1, 'free', 'ふらっとそこらへん'),
  (2, 'minutes', '5分だけ外に出る'),
  (3, 'destination', 'コンビニまで');

INSERT INTO users (
    id,
    nickname
)
VALUES (
    '019e1cd3-8194-7a36-816b-2f38206ca52c',
    'dummy user'
);

INSERT INTO characters (
    id,
    user_id,
    supporter_type
)
VALUES (
    '019e2000-0000-7000-8000-000000000001',
    '019e1cd3-8194-7a36-816b-2f38206ca52c',
    'dog'
);