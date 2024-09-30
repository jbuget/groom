CREATE TABLE IF NOT EXISTS rooms (
    id SERIAL PRIMARY KEY,
    slug VARCHAR(255) UNIQUE NOT NULL,
    space_name VARCHAR(255) UNIQUE NOT NULL
);

INSERT INTO rooms (slug, space_name) VALUES ('test1', 'enh-pjwc-nvm');
INSERT INTO rooms (slug, space_name) VALUES ('test2', 'fqt-mmbh-nvq');
