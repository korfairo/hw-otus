-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.events
(
    uuid uuid NOT NULL,
    user_uuid uuid NOT NULL,
    title character varying(256) COLLATE pg_catalog."default" NOT NULL,
    description character varying COLLATE pg_catalog."default",
    date timestamp with time zone NOT NULL,
    duration interval,
    notice interval,
    CONSTRAINT events_pkey PRIMARY KEY (uuid)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.events;
-- +goose StatementEnd
