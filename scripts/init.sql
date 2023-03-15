CREATE TABLE users (
	id SERIAL PRIMARY KEY,
	name VARCHAR(255),
	email VARCHAR(255),
	password VARCHAR(255),
	avatar VARCHAR(255),
	role INT
);

CREATE TABLE sessions (
	userid INT REFERENCES users(id) ON DELETE CASCADE,
	ip VARCHAR(32),
	refresh_token VARCHAR(255),
	expires TIMESTAMP,
	PRIMARY KEY (userid, ip)
);

CREATE TABLE verification_codes (
	email VARCHAR(255) PRIMARY KEY,
	value VARCHAR(255),
	expires TIMESTAMP
);

CREATE TABLE channels (
	id SERIAL PRIMARY KEY,
	title VARCHAR(255),
	description TEXT,
	status INT
);

CREATE TABLE comments (
	id SERIAL PRIMARY KEY,
	userid INT REFERENCES users(id) ON DELETE CASCADE,
	parent INT REFERENCES comments(id) ON DELETE CASCADE,
	text TEXT,
	time TIMESTAMP
);

CREATE TABLE tracks (
	id SERIAL PRIMARY KEY,
	title VARCHAR(255),
	perfomancer VARCHAR(255),
	year INT,
	audio VARCHAR(255)
);

CREATE TABLE tracks_likes (
	userid INT REFERENCES users(id) ON DELETE CASCADE,
	trackid INT REFERENCES tracks(id) ON DELETE CASCADE,
	time TIMESTAMP,
	PRIMARY KEY (userid, trackid)
);

CREATE TABLE tracks_comments (
	trackid INT REFERENCES tracks(id) ON DELETE CASCADE,
	commentid INT REFERENCES comments(id) ON DELETE CASCADE,
	PRIMARY KEY (trackid, commentid)	
);

CREATE TABLE schedule (
	channelid INT REFERENCES channels(id) ON DELETE CASCADE,
	trackid INT REFERENCES tracks(id) ON DELETE CASCADE,
	startdate TIMESTAMP,
	enddate TIMESTAMP,
	PRIMARY KEY (channelid, trackid, startdate)
);

CREATE TABLE news (
	id SERIAL PRIMARY KEY,
	title VARCHAR(255),
	content TEXT,
	publication_date TIMESTAMP
);

CREATE TABLE news_likes (
	userid INT REFERENCES users(id) ON DELETE CASCADE,
	newsid INT REFERENCES news(id) ON DELETE CASCADE,
	time TIMESTAMP,
	PRIMARY KEY (userid, newsid)
);

CREATE TABLE news_comments (
	newsid INT REFERENCES news(id) ON DELETE CASCADE,
	commentid INT REFERENCES comments(id) ON DELETE CASCADE,
	PRIMARY KEY (newsid, commentid)	
);


