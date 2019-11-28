--
-- PostgreSQL database dump
--

-- Dumped from database version 12.1
-- Dumped by pg_dump version 12.1

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

DROP DATABASE IF EXISTS tankyou_poc;
--
-- Name: tankyou_poc; Type: DATABASE; Schema: -; Owner: tankyou_poc
--

CREATE DATABASE tankyou_poc WITH TEMPLATE = template0 ENCODING = 'UTF8' LC_COLLATE = 'en_US.utf8' LC_CTYPE = 'en_US.utf8';


ALTER DATABASE tankyou_poc OWNER TO tankyou_poc;

\connect tankyou_poc

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: tasks; Type: TABLE; Schema: public; Owner: tankyou_poc
--

CREATE TABLE public.tasks
(
    id          integer           NOT NULL,
    resume      character varying NOT NULL,
    content     text,
    reporter_id integer,
    worker_id   integer,
    status      integer           NOT NULL default 0
);


ALTER TABLE public.tasks
    OWNER TO tankyou_poc;

--
-- Name: tasks_id_seq; Type: SEQUENCE; Schema: public; Owner: tankyou_poc
--

CREATE SEQUENCE public.tasks_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.tasks_id_seq
    OWNER TO tankyou_poc;

--
-- Name: tasks_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: tankyou_poc
--

ALTER SEQUENCE public.tasks_id_seq OWNED BY public.tasks.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: tankyou_poc
--

CREATE TABLE public.users
(
    id    integer           NOT NULL,
    name  character varying NOT NULL,
    email character varying NOT NULL
);


ALTER TABLE public.users
    OWNER TO tankyou_poc;

--
-- Name: user_id_seq; Type: SEQUENCE; Schema: public; Owner: tankyou_poc
--

CREATE SEQUENCE public.user_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.user_id_seq
    OWNER TO tankyou_poc;

--
-- Name: user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: tankyou_poc
--

ALTER SEQUENCE public.user_id_seq OWNED BY public.users.id;


--
-- Name: tasks id; Type: DEFAULT; Schema: public; Owner: tankyou_poc
--

ALTER TABLE ONLY public.tasks
    ALTER COLUMN id SET DEFAULT nextval('public.tasks_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: tankyou_poc
--

ALTER TABLE ONLY public.users
    ALTER COLUMN id SET DEFAULT nextval('public.user_id_seq'::regclass);

--
-- Name: tasks_id_seq; Type: SEQUENCE SET; Schema: public; Owner: tankyou_poc
--

SELECT pg_catalog.setval('public.tasks_id_seq', 1, false);


--
-- Name: user_id_seq; Type: SEQUENCE SET; Schema: public; Owner: tankyou_poc
--

SELECT pg_catalog.setval('public.user_id_seq', 2, true);


--
-- Name: tasks tasks_pk; Type: CONSTRAINT; Schema: public; Owner: tankyou_poc
--

ALTER TABLE ONLY public.tasks
    ADD CONSTRAINT tasks_pk PRIMARY KEY (id);


--
-- Name: users user_pk; Type: CONSTRAINT; Schema: public; Owner: tankyou_poc
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT user_pk PRIMARY KEY (id);


--
-- Name: tasks_id_uindex; Type: INDEX; Schema: public; Owner: tankyou_poc
--

CREATE UNIQUE INDEX tasks_id_uindex ON public.tasks USING btree (id);


--
-- Name: tasks_reporter_idx; Type: INDEX; Schema: public; Owner: tankyou_poc
--

CREATE INDEX tasks_reporter_idx ON public.tasks USING btree (reporter_id);


--
-- Name: tasks_resume_uindex; Type: INDEX; Schema: public; Owner: tankyou_poc
--

CREATE UNIQUE INDEX tasks_resume_uindex ON public.tasks USING btree (resume);


--
-- Name: tasks_worker_idx; Type: INDEX; Schema: public; Owner: tankyou_poc
--

CREATE INDEX tasks_worker_idx ON public.tasks USING btree (worker_id);


--
-- Name: user_email_uindex; Type: INDEX; Schema: public; Owner: tankyou_poc
--

CREATE UNIQUE INDEX user_email_uindex ON public.users USING btree (email);


--
-- Name: user_id_uindex; Type: INDEX; Schema: public; Owner: tankyou_poc
--

CREATE UNIQUE INDEX user_id_uindex ON public.users USING btree (id);


--
-- Name: user_name_uindex; Type: INDEX; Schema: public; Owner: tankyou_poc
--

CREATE UNIQUE INDEX user_name_uindex ON public.users USING btree (name);


--
-- PostgreSQL database dump complete
--

