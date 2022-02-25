create table questions
(
    id   integer
        constraint questions_pk
            primary key autoincrement,
    body text default ''
);

create unique index questions_id_uindex
    on questions (id);

create table options
(
    id         integer
        constraint options_pk
            primary key autoincrement,
    questionId integer,
    body       text    default '',
    correct    integer default 0,
    optionOrder    integer default 0
);

create unique index options_id_uindex
    on options (id);