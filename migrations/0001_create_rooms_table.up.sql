CREATE TABLE IF NOT EXISTS rooms (
    id SERIAL PRIMARY KEY,
    slug VARCHAR(255) UNIQUE NOT NULL,
    space_id VARCHAR(255) UNIQUE NOT NULL
);

INSERT INTO rooms (slug, space_id) VALUES ('test1', 'enh-pjwc-nvm');
INSERT INTO rooms (slug, space_id) VALUES ('test2', 'fqt-mmbh-nvq');
