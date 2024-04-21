CREATE TABLE IF NOT EXISTS public.users
(
    user_id bigserial NOT NULL,
    username varchar(30)  NOT NULL,
    password_hash varchar(100),
    PRIMARY KEY (user_id)
);

ALTER TABLE IF EXISTS public.users
    OWNER to postgres;

ALTER TABLE IF EXISTS public.users
    ADD CONSTRAINT users_username_uk UNIQUE (username);

CREATE TABLE IF NOT EXISTS public.tokenbl
(
    token TEXT NOT NULL,
    PRIMARY KEY (token)
);

ALTER TABLE IF EXISTS public.tokenbl
    OWNER to postgres;

-- Table: public.tasks

-- DROP TABLE IF EXISTS public.tasks;

CREATE TABLE IF NOT EXISTS public.tasks
(
    task_id bigserial NOT NULL,
    data jsonb,
    user_id bigint,
    CONSTRAINT tasks_pkey PRIMARY KEY (task_id),
    CONSTRAINT tasks_users_fk FOREIGN KEY (user_id)
        REFERENCES public.users (user_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.tasks
    OWNER to postgres;
-- Index: fki_tasks_users_fk

-- DROP INDEX IF EXISTS public.fki_tasks_users_fk;

CREATE INDEX IF NOT EXISTS fki_tasks_users_fk
    ON public.tasks USING btree
    (user_id ASC NULLS LAST)
    TABLESPACE pg_default;