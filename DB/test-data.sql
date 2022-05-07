use one_shot_url_test;

CREATE TABLE 00_urls(id int NOT NULL PRIMARY KEY AUTO_INCREMENT, long_url varchar(1000), short_url VARCHAR(100), updated_at DATETIME, created_at DATETIME, deleted_at DATETIME);

INSERT INTO 00_urls VALUES(1,"https://example.com", "00hu8jgt", "2022-05-07 02:14:01", "2022-05-07 02:14:01", NULL);
INSERT INTO 00_urls VALUES(2,"https://example.com/~canicani", "00bgtuq2", "2022-05-07 02:14:01", "2022-05-07 02:14:01", NULL);
INSERT INTO 00_urls VALUES(3,"https://exmaple.com/cani", "00Hy7Koi", "2022-05-07 02:14:01", "2022-05-07 02:14:01", NULL);
