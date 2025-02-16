-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TABLE country
(
    id      integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    country text NOT NULL UNIQUE
);

CREATE TABLE city
(
    id         integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    city       text NOT NULL,
    region     text,
    country_id integer REFERENCES country (id)
);

CREATE TABLE user_data
(
    id      integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    email   text NOT NULL UNIQUE,
    passwrd text NOT NULL,
    salt    text NOT NULL
);

CREATE TABLE profile_data
(
    user_id  integer REFERENCES user_data (id) ON DELETE CASCADE,
    username text UNIQUE,
    avatar   text,
    bio      text
);

CREATE TABLE category
(
    id      integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name    text NOT NULL UNIQUE
);

CREATE TABLE sight
(
    id          integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    rating      float,
    name        text NOT NULL,
    description text,
    city_id     integer REFERENCES city (id),
    country_id  integer REFERENCES country (id),
    latitude    double precision,
    longitude   double precision,
    UNIQUE (name, city_id),
    category_id     integer REFERENCES category (id)
);

CREATE TABLE image_data
(
    id       integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    "path"   text NOT NULL UNIQUE,
    sight_id integer REFERENCES sight (id)
);

CREATE TABLE journey
(
    id          integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name        text NOT NULL UNIQUE,
    user_id     integer REFERENCES user_data (id),
    description text
);

CREATE TABLE journey_sight
(
    id         integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    journey_id integer REFERENCES journey (id),
    sight_id   integer REFERENCES sight (id),
    priority   integer NOT NULL,
    UNIQUE (journey_id, sight_id)
);

CREATE TABLE feedback
(
    id       integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id  integer REFERENCES user_data (id),
    sight_id integer REFERENCES sight (id),
    rating   integer NOT NULL CHECK (rating > 0 AND rating <= 5),
    feedback text    NOT NULL
);

CREATE TABLE question
(
    id integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    text text NOT NULL
);

CREATE TABLE quiz
(
    id integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id integer REFERENCES user_data (id),
    rating integer NOT NULL CHECK (rating > 0 AND rating <= 5),
    question_id integer REFERENCES question (id),
    created_at timestamptz
);

CREATE TABLE album 
(
    id integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id integer REFERENCES user_data(id),
    name text NOT NULL,
    description text,
    UNIQUE(user_id, name)
);

CREATE TABLE album_photo
(
    id integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    album_id integer REFERENCES album(id),
    path text UNIQUE,
    description text
);

INSERT INTO country(country)
VALUES ('Абхазия');


INSERT INTO city (city, country_id)
VALUES ('Гагра', 1),
       ('Пицунда', 1),
       ('Гудаута', 1),
       ('Новый Афон', 1),
       ('Сухум', 1),
       ('Очамчыра', 1),
       ('Ткуарчал', 1),
       ('Гал', 1);

INSERT INTO category (name)
VALUES ('Рестораны'),
       ('Отели'),
       ('Достопримечательности');

INSERT INTO sight(name, description, city_id, country_id, category_id, longitude, latitude)
VALUES ('Озеро Рица',
        'Рица — горное озеро ледниково-тектонического происхождения на Западном Кавказе, в Гудаутском районе Абхазии',
        3,
        1,
        3,
        43.480130, 
        40.542047),
        ('Колоннада',
        'Колоннада в Гагре возводилась в конце 40-x - начале 50-х годов. Колоннада в Гагре стала символом возрождения и уверенности в лучшем светлом будущем. А еще колоннада должна была стать необычным памятником: для этого в проект включили 45 колонн (что символизировало год завершения страшной войны).',
        1,
        1,
        3,
        43.321597,   
        40.237041),
        ('Новоафонский монастырь',
        'Новоафонский Симоно-Кананитский мужской монастырь находится в Абхазии, на горе Афон, на высоте в 75 метров над уровнем моря. Основан он в 1875 году монахами — выходцами из Греции. Это один из крупнейших религиозных центров всего Кавказа, место паломничества и одна из популярных туристических достопримечательностей Абхазии.',
        4,
        1,
        3,
        43.088063, 
        40.820555);


INSERT INTO image_data(path, sight_id)
VALUES ('public/ritsa.jpg', 1),
       ('public/colonnada.jpg', 2),
       ('public/afon_monastyr.jpg', 3),;


INSERT INTO question(text)
VALUES ('Насколько вы удовлетворены удобством КудаТуда?'),
       ('Насколько интуитивно понятен интерфейс?');


CREATE OR REPLACE FUNCTION create_profile()
    RETURNS TRIGGER AS
$$
BEGIN
    INSERT INTO profile_data (user_id, username, bio, avatar)
    VALUES (NEW.id, NEW.email, '', '');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER create_profile_trigger
    AFTER INSERT
    ON user_data
    FOR EACH ROW
EXECUTE FUNCTION create_profile();


CREATE OR REPLACE FUNCTION update_sight_rating()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE sight
    SET rating = (SELECT AVG(rating) FROM feedback WHERE sight_id = NEW.sight_id)
    WHERE id = NEW.sight_id;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER after_insert_feedback
AFTER INSERT ON feedback
FOR EACH ROW
EXECUTE FUNCTION update_sight_rating();

CREATE TRIGGER after_update_feedback
AFTER UPDATE ON feedback
FOR EACH ROW
EXECUTE FUNCTION update_sight_rating();

CREATE TRIGGER after_delete_feedback
AFTER DELETE ON feedback
FOR EACH ROW
EXECUTE FUNCTION update_sight_rating();

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

DROP TABLE IF EXISTS user_data CASCADE;
DROP TABLE IF EXISTS city CASCADE;
DROP TABLE IF EXISTS country CASCADE;
DROP TABLE IF EXISTS sight CASCADE;
DROP TABLE IF EXISTS journey CASCADE;
DROP TABLE IF EXISTS journey_sight CASCADE;
DROP TABLE IF EXISTS image_data CASCADE;
DROP TABLE IF EXISTS feedback CASCADE;
DROP TABLE IF EXISTS profile_data CASCADE;
DROP TABLE IF EXISTS question CASCADE;
DROP TABLE IF EXISTS quiz CASCADE;
DROP TABLE IF EXISTS category CASCADE;
DROP TABLE IF EXISTS album CASCADE;
DROP TABLE IF EXISTS album_photo CASCADE;
