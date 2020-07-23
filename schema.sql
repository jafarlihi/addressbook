DROP TABLE IF EXISTS contact_list_entries;
DROP TABLE IF EXISTS contacts;
DROP TABLE IF EXISTS contact_lists;
DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id serial NOT NULL,
    username character varying NOT NULL UNIQUE,
    email character varying NOT NULL UNIQUE,
    password character varying NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE contact_lists (
    id serial NOT NULL,
    user_id INTEGER NOT NULL,
    name character varying NOT NULL,
    PRIMARY KEY(id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE contacts (
    id serial NOT NULL,
    user_id INTEGER NOT NULL,
    name character varying NOT NULL,
    surname character varying NOT NULL,
    phone INTEGER NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE contact_list_entries (
    contact_list INTEGER NOT NULL,
    contact INTEGER NOT NULL,
    UNIQUE (contact_list, contact),
    FOREIGN KEY (contact_list) REFERENCES contact_lists (id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (contact) REFERENCES contacts (id) ON UPDATE CASCADE ON DELETE CASCADE
)
