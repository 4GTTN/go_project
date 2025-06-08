INSERT INTO maps (name, price, description, location, image, is_deleted, created_at, updated_at)
VALUES
    ('Desert Arena', 10.99, 'A sandy battleground with ancient ruins', 'Desert', '/images/map1.jpg', false, NOW(), NOW()),
    ('Urban City', 15.99, 'A modern cityscape for intense battles', 'City', '/images/map2.jpg', false, NOW(), NOW()),
    ('Forest Hideout', 12.50, 'A dense forest with hidden paths', 'Forest', '/images/map3.jpg', false, NOW(), NOW());

INSERT INTO players (username, email, password, role, token, created_at, updated_at)

INSERT INTO games (creator_id, status, created_at, updated_at)
VALUES
    (1, 'draft', NOW(), NOW()),
    (1, 'formed', NOW(), NOW());

INSERT INTO game_maps (game_id, map_id)
VALUES
    (1, 1),
    (1, 2),
    (2, 3);