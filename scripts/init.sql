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
	logo VARCHAR(255),
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
	performancer VARCHAR(255),
	year INT,
	audio VARCHAR(255),
	duration BIGINT
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
	id SERIAL PRIMARY KEY,
	channelid INT REFERENCES channels(id) ON DELETE CASCADE,
	trackid INT REFERENCES tracks(id) ON DELETE CASCADE,
	startdate TIMESTAMP,
	enddate TIMESTAMP,
	UNIQUE (channelid, startdate)
);

CREATE TABLE news (
	id SERIAL PRIMARY KEY,
	title VARCHAR(255),
	content TEXT,
	image VARCHAR(255),
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

INSERT INTO users (name, email, password, role, avatar) VALUES 
  ('Admin', 'admin@yandex.ru', '$2a$10$NN9YDCE77tR/BhahmtwoWOcpC9zjs3InRrtQamBzqBQvN0YgR3U3q', 2, NULL),
  ('User 1', 'user1@yandex.ru', '$2a$10$H3RLgGWd6I872d1JTPiEXOyuL3x6SyfdvMoAEPQaGCupbZPuvJs6a', 1, 'files/avatar-1.jpg'),
  ('User 2', 'user2@yandex.ru', '$2a$10$DQnNxoRFOV2PEfKTMj5Y0uOSKwaLDu4PShMnIoVRTgnRlI.EJJ87y', 1, 'files/avatar-2.jpg'),
  ('User 3', 'user3@yandex.ru', '$2a$10$2hXFbpNkU4Zc7Q4hOUf/7Oa0IsKyaonQA.baVa24yt9DV3rUk.Um6', 1, 'files/avatar-3.jpg'),
  ('User 4', 'user4@yandex.ru', '$2a$10$V4iRIco2MBZzmlbSbrZzwuahAq3uTr9MG556jwzk2yDhDD0c4WicG', 1, 'files/avatar-4.jpg'),
  ('User 5', 'user5@yandex.ru', '$2a$10$kbjg4vTm6WBU73iR9T7KNeXdKps9OcxBqx21Z9kLkE3D4HLsdS6xu', 1, NULL);

INSERT INTO channels (title, status, logo, description) VALUES 
  ('Рэп музыка', 1, 'files/logo1.jpg', 'Рэп или исполнение рэпа (англ. rap, rapping), также известное как рифмование (rhyming), фристайл (freestyling, или в просторечии spitting)[1], эмсиинг (emceeing)[2] или MC’инг (MCing, M.C.’ing)[2][3] — музыкальная форма вокальной подачи, включающая в себя «рифму, ритмичную речь и уличный жаргон»[4], которая исполняется или произносится нараспев различными способами, как правило, под фоновый бит или музыкальный аккомпанемент[4]. Компоненты рэпа включают в себя «содержание» (то, что говорится), «флоу» (ритм, рифма) и «подачу» (каденция, тон)[5]. Рэп отличается от устной поэзии тем, что обычно исполняется в свободное время под музыкальное сопровождение[6]. Рэп — основной элемент хип-хоп-музыки, обычно ассоциируется именно с этим жанром и даже служит его метонимией; однако истоки рэпа на много лет предшествовали появлению хип-хоп-культуры.'),
  ('Русский рок', 1, 'files/logo2.jpg', 'Рок-музыка (англ. Rock music) — обобщающее название ряда направлений популярной музыки. Слово rock (в переводе с английского «качать», «укачивать», «качаться») в данном случае указывает на характерные для этих направлений ритмические ощущения, связанные с определённой формой движения, по аналогии с roll, twist, swing, shake…\nРок-музыка имеет самое большое количество направлений (около 180): от достаточно «лёгких», таких как танцевальный рок-н-ролл, поп-рок, мерсибит, до агрессивных — хэви-метала, глэм-метала, трэш-метала, блэк-метала и брутальных дэт-метала, грайндкора. '),
  ('Шансон', 1, 'files/logo3.jpg', 'Шансо́н (фр. chanson — песня) — французская эстрадная песня в стиле кабаре.\nЖанры шансона использовали певцы французских кабаре в конце XIX века — первой половине XX века. Наиболее известными из них являются Аристид Брюан, Мистингетт. Из кабаре данная модификация шансона перешла во французскую эстрадную музыку XX века.\nВ 1950-х годах оформились два главных направления оригинальной франкоязычной песни, существующие до настоящего времени:\n    Классический шансон, где первостепенное значение придается поэтической компоненте песни и автор, как правило, сам является исполнителем. Этот жанр связывается в первую очередь с именами Мориса Шевалье, Шарля Трене и Эдит Пиаф, которая продолжала традицию реалистической песни.\n    «Новый шансон», где используются новейшие приёмы современной легкой музыки, но по-прежнему очень требовательно относятся к текстам своих песен. Начало «нового шансона» связывают с именем Доминика А и относят к последнему десятилетию XX века.'),
  ('Поп музыка', 1, 'files/logo4.jpg', 'Впервые термин «pop song» в английском языке прозвучал ещё в 1926 году, однако корни поп-музыки уходят в историю глубже. Непосредственным предшественником поп-музыки была народная музыка, а также более поздние уличные романсы и баллады.\nСовременная поп-музыка формировалась параллельно с другими жанрами, такими как рок-музыка, и не всегда была отделима от них. В 1950-е и 1960-е её наиболее типичной формой был т. н. «традиционный поп» (traditional pop), который в СССР было принято называть «эстрадная музыка», «эстрада». Традиционный поп исполняется певцом-солистом под фоновый аккомпанемент. В США эстрада была тесно связана с джазом (Фрэнк Синатра), во Франции — с шансоном. Аналогичные исполнители были популярны и в СССР — Леонид Утёсов, Клавдия Шульженко, Марк Бернес, Владимир Трошин. Значительную часть поп-музыкальной сцены США составляют чернокожие исполнители в жанре соул.\nНастоящим прорывом в поп-музыке стало появление в 1970-е стиля «диско» (евродиско) и таких групп, как ABBA, Boney M, Dschinghis Khan, Bee Gees. Поп-музыка отныне вытесняет рок-н-ролл в качестве основной танцевальной музыки на дискотеках, и с этого времени танцевальная музыка (dance music) является одним из основных направлений в поп-музыке.'),
  ('Русские народные песни', 1, 'files/logo5.jpg', 'Русская народная песня — фольклорное произведение, которое сохраняется в народной памяти и передаётся в устной форме, продукт коллективного устного творчества русского народа. Относится к народному искусству.\nЧаще всего у народной песни нет определённого автора, или автор неизвестен, но известны и народные песни литературного происхождения. Существенная черта большинства жанров русской народной песни — непосредственная связь народной песни с бытом и трудовой деятельностью (например, песни трудовые, сопровождающие различные виды труда — бурлацкие, покосные, прополочные, жатвенные, молотильные и др.), обрядовые, сопровождающие земледельческие и семейные обряды и празднества, — колядки, масленичные, веснянки, купальские, свадебные, похоронные, игровые, календарные и т. п.). ');

INSERT INTO tracks (title, performancer, year, audio, duration) VALUES
  ('Lose Yourself', 'Eminem', 2002, 'files/1.ogg', 327000000000),
  ('Беспечный ангел', 'Ария', 2004, 'files/2.ogg', 238000000000),
  ('Владимирский централ', 'Михаил круг', 1998, 'files/3.ogg', 268000000000),
  ('VICIOUS', 'Thomas Day', 2023, 'files/4.ogg', 132000000000),
  ('Шел казак на побывку домой', 'Пелагея', 2014, 'files/5.ogg', 273000000000);

INSERT INTO schedule (channelid, trackid, startdate, enddate) VALUES
  (1, 1, NOW(), NOW() + interval '1' day),
  (2, 2, NOW(), NOW() + interval '1' day),
  (3, 3, NOW(), NOW() + interval '1' day),
  (4, 4, NOW(), NOW() + interval '1' day),
  (5, 5, NOW(), NOW() + interval '1' day);

INSERT INTO news (title, publication_date, image, content) VALUES
  ('Новость 1', NOW(), 'files/image1.jpg', '# Adspicit domum

## Venit et suae

Lorem markdownum inpia et fatendo ad ignotas! Tu acrior manuum! Ciboque virgo
vidi tamen, per cubat et aequoris capaci inire adolescere qui. Esse Arcton
ignavi secutum. Pondus *nascentur fronde* loci quae, nascentia coniuge,
[falsaque](http://hiccaelum.com/ipsesinat) defensus intus, multa petit redit
metit laboris.

Troiae sanguinis, nubila dat [tuum medicamine
mihi](http://potest.io/saevo-quas.html) Cecropidum deum suppressis conchis
pulcherrima pendet perhorruit iubeo Rutulum. Foedere corpora, rotat arbor enim,
vinctoque rogat huic humo plaga, **aurea**? Quod tamen passa Syron vixisti
suique: fundae tota, de cum Rex! Siquis patebat emerguntque suasit praecluserat
ab cvrrvs iussus aequare inposita bello ter pennatis quamquam me vincla sic
fingetur, vivacemque!

## Est iecur cursu

Mortalem aera; vilibus Thessala rettulit et et omnes delia. Alis Caicum, lucent
oculis: nimbos et addit factisque terrore dementer quidem incaluitque solet?

    if (simplex_flash_degauss != drmPptpVideo.vlog(utility_pmu, 71)) {
        snow_bankruptcy_dial(parityPermalinkBittorrent);
        utf_piracy += rgbProgramMemory - zipHfsSuperscalar;
    } else {
        cleanSystem.play -= 925592;
    }
    search = algorithm / digital_optical + android * flatbedBootMicrocomputer *
            firewire_faq_flat + softwareRawAnimated;
    broadband_servlet_asp(33);
    if (phreaking + bug_tag == ospf_path_operating - associationSku) {
        typeface_duplex_optical(backside_irq);
        cursor_smartphone.rowSkinRuntime += reality_bounce(833825);
    }
    trojanExcel = packet;

## Fluctus mihi auctor fratres despectus manifesta cunas

Veste flexi quare arces deprensi et parte. Vos illis, prensoque tegitur, et
vulnera nostro, quid. Vultu fessa nat fronte hiatus?

    if (pipelineBootName) {
        tweakBankruptcy = 71;
        vpiWwwGis += mashup_dvr;
        smb_memory = 2;
    }
    type(youtubeShell, phpPrebindingToslink(5));
    var server = multimedia_wordart_hashtag + nat_gnu_prebinding;
    if (expansionSmbCharacter + active + data >= pitchAddressMeta + unit) {
        bar.softwareGigaflops = menu(file_soft + halftone, software_point,
                megabyteIvrPrint);
        gatewayW = clipboardLag / cycleIcmp + down_install_jumper;
        compactLan(program_pptp + wDrag);
    }

## Arva amas intendit

Ire pes lacrimas manus regia, nisi mihi aures rabiem imis inque fore geminae.
*Sibi concussit papavera* fuit pectora?

*Non circuit arcebatque* stramine lata hospes Cepheni facilem totosque,
traxerunt motibus deorum fama eum forte sub et Penthea. Querellae terraque
crepitantia manus terribilesque manus putes; suas, et ipse si igne texique,
siqua, et. Totum bracchia duris vivitur doleret? Ruinae ventos ac Finierat agmen
gemitus quodque meus eam procul patefecit? Misero fratrisque, materiem corpore,
oculos qua.'),
  ('Новость 2', NOW() - interval '1' hour, 'files/image2.jpg', '# Заголовок'),
  ('Новость 3', NOW() - interval '2' hour, 'files/image3.jpg', '# Заголовок'),
  ('Новость 4', NOW() - interval '3' hour, 'files/image1.jpg', '# Заголовок'),
  ('Новость 5', NOW() - interval '4' hour, 'files/image2.jpg', '# Заголовок'),
  ('Новость 6', NOW() - interval '5' hour, 'files/image3.jpg', '# Заголовок'),
  ('Новость 7', NOW() - interval '6' hour, 'files/image1.jpg', '# Заголовок'),
  ('Новость 8', NOW() - interval '7' hour, 'files/image2.jpg', '# Заголовок'),
  ('Новость 9', NOW() - interval '8' hour, 'files/image3.jpg', '# Заголовок'),
  ('Новость 10', NOW() - interval '9' hour, 'files/image1.jpg', '# Заголовок'),
  ('Новость 11', NOW() - interval '10' hour, 'files/image2.jpg', '# Заголовок'),
  ('Новость 12', NOW() - interval '11' hour, 'files/image3.jpg', '# Заголовок'),
  ('Новость 13', NOW() - interval '12' hour, 'files/image1.jpg', '# Заголовок'),
  ('Новость 14', NOW() - interval '13' hour, 'files/image2.jpg', '# Заголовок'),
  ('Новость 15', NOW() - interval '14' hour, 'files/image3.jpg', '# Заголовок'),
  ('Новость 16', NOW() - interval '15' hour, 'files/image1.jpg', '# Заголовок'),
  ('Новость 17', NOW() - interval '16' hour, 'files/image2.jpg', '# Заголовок'),
  ('Новость 18', NOW() - interval '17' hour, 'files/image3.jpg', '# Заголовок'),
  ('Новость 19', NOW() - interval '18' hour, 'files/image1.jpg', '# Заголовок'),
  ('Новость 20', NOW() - interval '19' hour, 'files/image2.jpg', '# Заголовок'),
  ('Новость 21', NOW() - interval '20' hour, 'files/image3.jpg', '# Заголовок'),
  ('Новость 22', NOW() - interval '21' hour, 'files/image1.jpg', '# Заголовок'),
  ('Новость 23', NOW() - interval '22' hour, 'files/image2.jpg', '# Заголовок'),
  ('Новость 24', NOW() - interval '23' hour, 'files/image3.jpg', '# Заголовок'),
  ('Новость 25', NOW() - interval '24' hour, 'files/image1.jpg', '# Заголовок'),
  ('Новость 26', NOW() - interval '25' hour, 'files/image2.jpg', '# Заголовок'),
  ('Новость 27', NOW() - interval '26' hour, 'files/image3.jpg', '# Заголовок'),
  ('Новость 28', NOW() - interval '27' hour, 'files/image1.jpg', '# Заголовок'),
  ('Новость 29', NOW() - interval '28' hour, 'files/image2.jpg', '# Заголовок'),
  ('Новость 30', NOW() - interval '29' hour, 'files/image3.jpg', '# Заголовок'),
  ('Новость 31', NOW() - interval '30' hour, 'files/image1.jpg', '# Заголовок'),
  ('Новость 32', NOW() - interval '31' hour, 'files/image2.jpg', '# Заголовок'),
  ('Новость 33', NOW() - interval '32' hour, 'files/image3.jpg', '# Заголовок');

INSERT INTO comments (userid, parent, text, time) VALUES
  (1, NULL, 'Я первый', NOW() - interval '20' minute),
  (2, NULL, 'Кто это вообще слушает?', NOW() - interval '5' minute),
  (3, NULL, 'Лайк не глядя', NOW() - interval '3' minute),
  (4, NULL, 'Жаль подписки не добавили', NOW() - interval '8' minute),
  (5, NULL, 'Ура! Комментарии работают!', NOW() - interval '12' minute),
  (2, 1, 'Боже чел', NOW());

INSERT INTO tracks_comments (trackid, commentid) VALUES
  (1, 1), (1, 6), (1, 2), (1, 3), (1, 4), (1, 5),
  (2, 1), (2, 6), (2, 2), (2, 3), (2, 5), 
  (3, 1), (3, 6), (3, 2), (3, 4), (3, 5), 
  (4, 1), (4, 6), (4, 3), (4, 4), (4, 5), 
  (5, 2), (5, 3), (5, 4), (5, 5);

INSERT INTO news_comments (newsid, commentid) VALUES
  (1, 2), (1, 3), (1, 4), (1, 5), 
  (2, 1), (2, 6), (2, 3), (2, 4), (2, 5), 
  (3, 1), (3, 6), (3, 2), (3, 4), (3, 5), 
  (4, 1), (4, 6), (4, 2), (4, 3), (4, 5);

