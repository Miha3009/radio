INSERT INTO users (name, email, password, role) VALUES 
  ('Admin', 'admin@yandex.ru', '$2a$10$NN9YDCE77tR/BhahmtwoWOcpC9zjs3InRrtQamBzqBQvN0YgR3U3q', 2);

INSERT INTO channels (id, title, description, status) VALUES 
  (1, 'Channel 1', 'This is channel 1', 1),
  (2, 'Channel 2', 'This is channel 2', 1);

INSERT INTO tracks (id, title, perfomancer, year, audio) VALUES
  (1, 'Track 1', 'Someone', 2019, 'files/1.ogg'),
  (2, 'Track 2', 'Another', 2021, 'files/2.ogg');

INSERT INTO schedule (channelid, trackid, startdate, enddate) VALUES
  (1, 1, NOW(), NOW() + interval '1' day),
  (2, 2, NOW(), NOW() + interval '1' day);
