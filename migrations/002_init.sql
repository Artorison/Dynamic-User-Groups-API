INSERT INTO users (id, name) VALUES (1000, 'USER1');
INSERT INTO users (id, name) VALUES (1002, 'USER2');
INSERT INTO users (id, name) VALUES (1004, 'USER3');

INSERT INTO segments (slug) VALUES ('DISCOUNT_50');
INSERT INTO segments (slug) VALUES ('DISCOUNT_30');
INSERT INTO segments (slug) VALUES ('SUPPORT');
INSERT INTO segments (slug) VALUES ('MESSAGES');
INSERT INTO segments (slug) VALUES ('VIDEO');

INSERT INTO user_segments (user_id, segment_id) 
VALUES (1000, (SELECT id FROM segments WHERE slug = 'SUPPORT'));

INSERT INTO user_segments (user_id, segment_id) 
VALUES (1000, (SELECT id FROM segments WHERE slug = 'DISCOUNT_30'));

INSERT INTO user_segments (user_id, segment_id) 
VALUES (1000, (SELECT id FROM segments WHERE slug = 'MESSAGES'));

INSERT INTO user_segments (user_id, segment_id) 
VALUES (1002, (SELECT id FROM segments WHERE slug = 'DISCOUNT_50'));

INSERT INTO user_segments (user_id, segment_id) 
VALUES (1004, (SELECT id FROM segments WHERE slug = 'VIDEO'));