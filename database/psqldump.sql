--
-- PostgreSQL database dump
--

\restrict vfDaLqC7kFxuknK9dVSke5kDzqBOdITaIoDPsRExxBbXSd6fjQxiQVhonTw4qcq

-- Dumped from database version 18.0
-- Dumped by pg_dump version 18.0

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: site; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA site;


ALTER SCHEMA site OWNER TO postgres;

--
-- Name: SCHEMA site; Type: COMMENT; Schema: -; Owner: postgres
--

COMMENT ON SCHEMA site IS 'база данных сайта';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: categories; Type: TABLE; Schema: site; Owner: postgres
--

CREATE TABLE site.categories (
    id uuid NOT NULL,
    name text NOT NULL,
    parent_id uuid,
    description text NOT NULL
);


ALTER TABLE site.categories OWNER TO postgres;

--
-- Name: chats; Type: TABLE; Schema: site; Owner: postgres
--

CREATE TABLE site.chats (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    seller_id uuid NOT NULL,
    customer_id uuid NOT NULL,
    listing_id uuid NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE site.chats OWNER TO postgres;

--
-- Name: complaints; Type: TABLE; Schema: site; Owner: postgres
--

CREATE TABLE site.complaints (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    complainant_id uuid NOT NULL,
    listing_id uuid,
    target_user_id uuid,
    complaint_type text NOT NULL,
    description text,
    status text NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone,
    admin_id uuid,
    resolution_comment text
);


ALTER TABLE site.complaints OWNER TO postgres;

--
-- Name: email_confirmations; Type: TABLE; Schema: site; Owner: postgres
--

CREATE TABLE site.email_confirmations (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    token text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    expires_at timestamp with time zone DEFAULT (CURRENT_TIMESTAMP + '24:00:00'::interval) NOT NULL,
    confirmed boolean DEFAULT false NOT NULL,
    confirmed_at timestamp with time zone
);


ALTER TABLE site.email_confirmations OWNER TO postgres;

--
-- Name: listing_images; Type: TABLE; Schema: site; Owner: postgres
--

CREATE TABLE site.listing_images (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    listing_id uuid NOT NULL,
    image_id text NOT NULL,
    "position" integer NOT NULL,
    uploaded_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE site.listing_images OWNER TO postgres;

--
-- Name: listings; Type: TABLE; Schema: site; Owner: postgres
--

CREATE TABLE site.listings (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    category_id uuid NOT NULL,
    title text NOT NULL,
    description text NOT NULL,
    price numeric NOT NULL,
    location_id uuid NOT NULL,
    condition text NOT NULL,
    status text DEFAULT 'moderation'::text NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone,
    expiration_date timestamp with time zone,
    views_count integer DEFAULT 0 NOT NULL,
    contact_phone text
);


ALTER TABLE site.listings OWNER TO postgres;

--
-- Name: locations; Type: TABLE; Schema: site; Owner: postgres
--

CREATE TABLE site.locations (
    id uuid NOT NULL,
    country text NOT NULL,
    region text NOT NULL,
    city text NOT NULL,
    postal_code text NOT NULL
);


ALTER TABLE site.locations OWNER TO postgres;

--
-- Name: messages; Type: TABLE; Schema: site; Owner: postgres
--

CREATE TABLE site.messages (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    sender_id uuid NOT NULL,
    message_text text NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    chat_id uuid NOT NULL
);


ALTER TABLE site.messages OWNER TO postgres;

--
-- Name: password_reset_tokens; Type: TABLE; Schema: site; Owner: postgres
--

CREATE TABLE site.password_reset_tokens (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    token text NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    expires_at timestamp with time zone NOT NULL
);


ALTER TABLE site.password_reset_tokens OWNER TO postgres;

--
-- Name: user_favorites; Type: TABLE; Schema: site; Owner: postgres
--

CREATE TABLE site.user_favorites (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    listing_id uuid NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE site.user_favorites OWNER TO postgres;

--
-- Name: users; Type: TABLE; Schema: site; Owner: postgres
--

CREATE TABLE site.users (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    username text NOT NULL,
    email text NOT NULL,
    password_hash text NOT NULL,
    role text NOT NULL,
    status text NOT NULL,
    phone_number text,
    registration_date timestamp with time zone NOT NULL,
    last_login time with time zone,
    profile_picture_id text
);


ALTER TABLE site.users OWNER TO postgres;

--
-- Data for Name: categories; Type: TABLE DATA; Schema: site; Owner: postgres
--

COPY site.categories (id, name, parent_id, description) FROM stdin;
c9255afc-856c-4790-a89a-5e94cd87a71e	Услуги	\N	
d2f664ce-6e70-4c05-ba7b-a2606ea75f55	Авто	\N	
38058e35-8b30-48c6-a711-9a511fe81d10	Недвижимость	\N	
3a3c3593-006c-4eeb-bf2c-35a494306124	Электроника	\N	
44f960e8-07b9-475c-a9a3-99328481767e	Работа	\N	
3935cdc1-b21c-4fce-a866-654433456c55	Мода	\N	
5a3aa40e-c100-4843-bf0d-af19a42d692f	Для дома	\N	
95a56d54-9a54-42bf-a00d-0b9c6b563756	Хобби	\N	
0e0b4e4d-9d36-426a-95b3-e2444451232e	Ремонт техники	c9255afc-856c-4790-a89a-5e94cd87a71e	
7b6110f3-cc50-4b54-9880-395e2c31c11e	Красота и здоровье	c9255afc-856c-4790-a89a-5e94cd87a71e	
4971bdb4-b860-42ee-8353-8b5f3fefda6c	Образование	c9255afc-856c-4790-a89a-5e94cd87a71e	
f9d63541-a080-43f4-a460-642c5c60ea68	Транспортные услуги	c9255afc-856c-4790-a89a-5e94cd87a71e	
0d308470-1259-46a0-b1b2-893f389d51da	Ремонт и строительство	c9255afc-856c-4790-a89a-5e94cd87a71e	
bd6eafb0-ec4b-46c0-936b-1ab78b74ee3b	Легковые автомобили	d2f664ce-6e70-4c05-ba7b-a2606ea75f55	
373da658-39a0-4da4-83b3-3c8fb535930d	Мотоциклы	d2f664ce-6e70-4c05-ba7b-a2606ea75f55	
78674ca7-8245-423d-b7e2-d84966bb19ea	Грузовики	d2f664ce-6e70-4c05-ba7b-a2606ea75f55	
cd5d2bb6-43ad-42fc-9122-f213d18477cf	Спецтехника	d2f664ce-6e70-4c05-ba7b-a2606ea75f55	
0ed219a3-dfe1-498e-84a9-81466cba3213	Запчасти и аксессуары	d2f664ce-6e70-4c05-ba7b-a2606ea75f55	
3302af03-a3b4-426a-b369-afd59d5d0b83	Квартиры	38058e35-8b30-48c6-a711-9a511fe81d10	
2ebc2d48-ad9b-4156-bcac-d1600f180327	Дома и коттеджи	38058e35-8b30-48c6-a711-9a511fe81d10	
e6b2f6b2-ed72-42c6-a5cf-197a0f832cbe	Коммерческая недвижимость	38058e35-8b30-48c6-a711-9a511fe81d10	
1e5f4980-30c0-41d7-bee0-23c433db9792	Земельные участки	38058e35-8b30-48c6-a711-9a511fe81d10	
e6d5c834-0736-41c5-8e5b-0763f591320a	Гаражи и стоянки	38058e35-8b30-48c6-a711-9a511fe81d10	
d60f1f8c-bfa5-4868-8327-12b879fbb6bb	Смартфоны и планшеты	3a3c3593-006c-4eeb-bf2c-35a494306124	
36a2b144-08b7-477a-a59c-699d168648a3	Ноутбуки и компьютеры	3a3c3593-006c-4eeb-bf2c-35a494306124	
283b2754-528e-48b5-9a6d-88266dfcf439	Телевизоры и аудио	3a3c3593-006c-4eeb-bf2c-35a494306124	
14b5b78c-cfb2-4e01-9712-5a94445aaad3	Фото и видео	3a3c3593-006c-4eeb-bf2c-35a494306124	
596bb5bc-382b-4ae6-acaa-668c0c0970f3	Бытовая техника	3a3c3593-006c-4eeb-bf2c-35a494306124	
467f9cbe-5136-4666-90b7-df4277afb712	Вакансии	44f960e8-07b9-475c-a9a3-99328481767e	
62767c9d-5da9-4b76-a300-8072e1beaf43	Резюме	44f960e8-07b9-475c-a9a3-99328481767e	
d9eae735-1cc2-44b8-a185-bd44d8c0d898	Фриланс	44f960e8-07b9-475c-a9a3-99328481767e	
7cfff2ee-fd03-411f-8ed9-d104c8c470ea	Удаленная работа	44f960e8-07b9-475c-a9a3-99328481767e	
e375781c-e701-44c6-b0a1-de801efaa834	Подработка	44f960e8-07b9-475c-a9a3-99328481767e	
dc4a36da-025d-4bec-af82-3adea25e87af	Одежда	3935cdc1-b21c-4fce-a866-654433456c55	
1dea4f6d-0ae1-4fb7-972f-bd1a8a08a377	Обувь	3935cdc1-b21c-4fce-a866-654433456c55	
741ae748-147a-4a31-9de5-0da55ee5af42	Аксессуары	3935cdc1-b21c-4fce-a866-654433456c55	
68a7b50d-1f3c-42b5-af1d-0e5710c42c64	Часы и украшения	3935cdc1-b21c-4fce-a866-654433456c55	
5254086e-c579-45f7-908f-ce0fdb3c0775	Косметика	3935cdc1-b21c-4fce-a866-654433456c55	
cf49bea0-c025-4e78-aea2-356ad838feb5	Мебель	5a3aa40e-c100-4843-bf0d-af19a42d692f	
4a0836f8-40a7-4b7c-a3b0-ef2ab2dfbd9f	Интерьер	5a3aa40e-c100-4843-bf0d-af19a42d692f	
3583c014-15dc-4b42-866a-a9e3505e3e2c	Посуда	5a3aa40e-c100-4843-bf0d-af19a42d692f	
97171caa-22b9-4abe-b224-85760e98a2a2	Текстиль	5a3aa40e-c100-4843-bf0d-af19a42d692f	
7d3ec9ac-74d6-4681-bc97-62e80fa10dcb	Хозяйственные товары	5a3aa40e-c100-4843-bf0d-af19a42d692f	
0600b600-24ae-4c2e-a502-ec0c940a6d1d	Спорт и отдых	95a56d54-9a54-42bf-a00d-0b9c6b563756	
bd4dcc63-8009-403f-ac4b-4dc8f9bf8a6f	Книги и журналы	95a56d54-9a54-42bf-a00d-0b9c6b563756	
916a1568-b5ae-4899-aca9-8edc44f4f120	Коллекционирование	95a56d54-9a54-42bf-a00d-0b9c6b563756	
65ac947d-eeea-4796-b3c5-ee5152f681ba	Музыкальные инструменты	95a56d54-9a54-42bf-a00d-0b9c6b563756	
613b781e-0ce9-4c0e-80dc-2a8e841e156c	Туризм и рыбалка	95a56d54-9a54-42bf-a00d-0b9c6b563756	
\.


--
-- Data for Name: chats; Type: TABLE DATA; Schema: site; Owner: postgres
--

COPY site.chats (id, seller_id, customer_id, listing_id, created_at) FROM stdin;
\.


--
-- Data for Name: complaints; Type: TABLE DATA; Schema: site; Owner: postgres
--

COPY site.complaints (id, complainant_id, listing_id, target_user_id, complaint_type, description, status, created_at, updated_at, admin_id, resolution_comment) FROM stdin;
\.


--
-- Data for Name: email_confirmations; Type: TABLE DATA; Schema: site; Owner: postgres
--

COPY site.email_confirmations (id, user_id, token, created_at, expires_at, confirmed, confirmed_at) FROM stdin;
fb983ef3-574a-4240-9743-cd8cca86d3a4	93721bbe-d0de-433c-8546-eae0e1a8e064	\N	2025-10-27 22:12:53.494372+03	2025-10-28 22:12:53.494372+03	t	2025-10-27 22:13:31.115297+03
e555e96b-3e0d-4037-a4d0-7c3cb9f281b8	dacd6f80-8db7-4cea-9605-9b2d5fe4a784	\N	2025-11-03 15:55:10.441226+03	2025-11-04 15:55:10.441226+03	t	2025-11-03 15:55:31.754007+03
\.


--
-- Data for Name: listing_images; Type: TABLE DATA; Schema: site; Owner: postgres
--

COPY site.listing_images (id, listing_id, image_id, "position", uploaded_at) FROM stdin;
\.


--
-- Data for Name: listings; Type: TABLE DATA; Schema: site; Owner: postgres
--

COPY site.listings (id, user_id, category_id, title, description, price, location_id, condition, status, created_at, updated_at, expiration_date, views_count, contact_phone) FROM stdin;
\.


--
-- Data for Name: locations; Type: TABLE DATA; Schema: site; Owner: postgres
--

COPY site.locations (id, country, region, city, postal_code) FROM stdin;
ac728115-b4ce-4cf9-95f1-e38904e9a48e	Беларусь	Минская область	Минск	220000
96281f4c-a010-4658-bbba-2c93312e35dd	Беларусь	Минская область	Борисов	222001
dd7bef66-40ea-4acd-91b5-6fb7dc3a62cf	Беларусь	Минская область	Солигорск	223001
b065fc1d-fdfc-46e0-8bcc-0694733dcd2d	Беларусь	Минская область	Молодечно	222301
b8938925-c3a0-4b83-940c-62713ee3f5e0	Беларусь	Минская область	Слуцк	223801
6901ae4e-26cc-4806-ae1d-b4f47d320c50	Беларусь	Гомельская область	Гомель	246001
57a84f32-e7f7-42fa-b649-02de4bee762d	Беларусь	Гомельская область	Мозырь	247750
043a6354-24bc-4cc9-9065-c15831b10e3d	Беларусь	Гомельская область	Жлобин	247001
58e9cf1a-54f1-49ad-ab6b-687b16517264	Беларусь	Гомельская область	Речица	247501
01127888-d4cd-48c3-9209-08cde154b32d	Беларусь	Гомельская область	Петриков	247601
18570587-28ed-4942-8d00-32c221bf591a	Беларусь	Могилёвская область	Могилёв	212001
11fdf79e-51f8-4303-b6c0-d6aeb1e1584b	Беларусь	Могилёвская область	Орша	211001
9de552e5-e2c7-4f6b-9978-c7d2e0ffd1eb	Беларусь	Могилёвская область	Шклов	213001
124d547d-f021-4cba-88aa-2ec29ccf361b	Беларусь	Могилёвская область	Кричев	213201
a93c05d6-b928-4bbc-b59f-905c97dcb33b	Беларусь	Могилёвская область	Славгород	213301
8b430862-e3df-497b-9720-d788b7dcdacb	Беларусь	Витебская область	Витебск	210001
d96576f9-42a4-4579-91c7-d63fc028da1d	Беларусь	Витебская область	Новополоцк	211401
c5d57d67-5c9e-4d54-8b75-b62438cee372	Беларусь	Витебская область	Полоцк	211451
01b8f6c1-00f9-46ad-8304-a0a0fe1e9d2d	Беларусь	Витебская область	Браслав	211301
eef59260-9f49-42c1-a1df-9d13d8d38558	Беларусь	Витебская область	Орша	211001
cf8f67a1-4328-44bb-9971-a7d4c6924d84	Беларусь	Гродненская область	Гродно	230001
804e5af6-e6b5-47ae-abb3-8a3578061a2a	Беларусь	Гродненская область	Лида	231301
bb4b79ac-a976-423f-ab4a-cf568d48ada1	Беларусь	Гродненская область	Слоним	231201
6e83cad7-f545-4fcb-8b0e-296a6987a825	Беларусь	Гродненская область	Щучин	231401
4812c971-0fc3-4ea0-8676-b452bdc78bc8	Беларусь	Гродненская область	Волковыск	231101
7c6a9f76-55ba-421f-87a8-9c5c76f3c290	Беларусь	Брестская область	Брест	224001
0952c152-3d38-4ef1-ac59-e545f69a2907	Беларусь	Брестская область	Пинск	225001
ac502f38-6ebf-4530-8a15-7c66c371e60b	Беларусь	Брестская область	Барановичи	225301
32287015-29e1-476d-9451-eaea6c63769a	Беларусь	Брестская область	Иваново	225801
ae76e4d6-68f4-4be3-8b72-f2cac567dd3d	Беларусь	Брестская область	Кобрин	225351
\.


--
-- Data for Name: messages; Type: TABLE DATA; Schema: site; Owner: postgres
--

COPY site.messages (id, sender_id, message_text, created_at, chat_id) FROM stdin;
\.


--
-- Data for Name: password_reset_tokens; Type: TABLE DATA; Schema: site; Owner: postgres
--

COPY site.password_reset_tokens (id, user_id, token, created_at, expires_at) FROM stdin;
8fb12abe-a5ed-4238-8445-e3712968f228	93721bbe-d0de-433c-8546-eae0e1a8e064	bbccefae6914eaced55b70ca03981c5f4409980342685c8d9014c249913213cf	2025-11-03 19:55:35.360844+03	2025-11-03 20:55:35.360675+03
\.


--
-- Data for Name: user_favorites; Type: TABLE DATA; Schema: site; Owner: postgres
--

COPY site.user_favorites (id, user_id, listing_id, created_at) FROM stdin;
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: site; Owner: postgres
--

COPY site.users (id, username, email, password_hash, role, status, phone_number, registration_date, last_login, profile_picture_id) FROM stdin;
dacd6f80-8db7-4cea-9605-9b2d5fe4a784	asdfzxcv	1238216@mtp.by	$2a$14$kalCti72XY3QX8ZlMJZhDuvbF4/6Wj/Vv6csJT1eHC/GyXUx84amC	user	active	\N	2025-11-03 15:55:08.655264+03	\N	\N
93721bbe-d0de-433c-8546-eae0e1a8e064	Матвей	zyazyulyam@bk.ru	$2a$14$aQi2fICu8nqjn/yf1KqQR.HHsl5taUacynjc4NWytbMsL4NE.VrMC	admin	active	+375291110235	2025-10-27 22:12:52.094949+03	19:57:34.828474+03	54f2af9e-bd39-469f-b120-900dd34bee1b
\.


--
-- Name: categories categories_pkey; Type: CONSTRAINT; Schema: site; Owner: postgres
--

ALTER TABLE ONLY site.categories
    ADD CONSTRAINT categories_pkey PRIMARY KEY (id);


--
-- Name: categories categoriesid; Type: CONSTRAINT; Schema: site; Owner: postgres
--

ALTER TABLE ONLY site.categories
    ADD CONSTRAINT categoriesid UNIQUE (id);


--
-- Name: chats chats_pkey; Type: CONSTRAINT; Schema: site; Owner: postgres
--

ALTER TABLE ONLY site.chats
    ADD CONSTRAINT chats_pkey PRIMARY KEY (id);


--
-- Name: chats chats_user1_id_user2_id_listing_id_key; Type: CONSTRAINT; Schema: site; Owner: postgres
--

ALTER TABLE ONLY site.chats
    ADD CONSTRAINT chats_user1_id_user2_id_listing_id_key UNIQUE (seller_id, customer_id, listing_id);


--
-- Name: complaints complaints_pkey; Type: CONSTRAINT; Schema: site; Owner: postgres
--

ALTER TABLE ONLY site.complaints
    ADD CONSTRAINT complaints_pkey PRIMARY KEY (id);


--
-- Name: email_confirmations email_confirmations_pkey; Type: CONSTRAINT; Schema: site; Owner: postgres
--

ALTER TABLE ONLY site.email_confirmations
    ADD CONSTRAINT email_confirmations_pkey PRIMARY KEY (id);


--
-- Name: email_confirmations email_confirmations_token_key; Type: CONSTRAINT; Schema: site; Owner: postgres
--

ALTER TABLE ONLY site.email_confirmations
    ADD CONSTRAINT email_confirmations_token_key UNIQUE (token);


--
-- Name: listing_images imageid; Type: CONSTRAINT; Schema: site; Owner: postgres
--

ALTER TABLE ONLY site.listing_images
    ADD CONSTRAINT imageid UNIQUE (id);


--
-- Name: listings list; Type: CONSTRAINT; Schema: site; Owner: postgres
--

ALTER TABLE ONLY site.listings
    ADD CONSTRAINT list UNIQUE (id);


--
-- Name: listing_images listing_images_pkey; Type: CONSTRAINT; Schema: site; Owner: postgres
--

ALTER TABLE ONLY site.listing_images
    ADD CONSTRAINT listing_images_pkey PRIMARY KEY (id);


--
-- Name: listings listings_pkey; Type: CONSTRAINT; Schema: site; Owner: postgres
--

ALTER TABLE ONLY site.listings
    ADD CONSTRAINT listings_pkey PRIMARY KEY (id);


--
-- Name: locations locations_pkey; Type: CONSTRAINT; Schema: site; Owner: postgres
--

ALTER TABLE ONLY site.locations
    ADD CONSTRAINT locations_pkey PRIMARY KEY (id);


--
-- Name: locations locationsid; Type: CONSTRAINT; Schema: site; Owner: postgres
--

ALTER TABLE ONLY site.locations
    ADD CONSTRAINT locationsid UNIQUE (id);


--
-- Name: messages messages_pkey; Type: CONSTRAINT; Schema: site; Owner: postgres
--

ALTER TABLE ONLY site.messages
    ADD CONSTRAINT messages_pkey PRIMARY KEY (id);


--
-- Name: messages messagesid; Type: CONSTRAINT; Schema: site; Owner: postgres
--

ALTER TABLE ONLY site.messages
    ADD CONSTRAINT messagesid UNIQUE (id);


--
-- Name: password_reset_tokens password_reset_tokens_pkey; Type: CONSTRAINT; Schema: site; Owner: postgres
--

ALTER TABLE ONLY site.password_reset_tokens
    ADD CONSTRAINT password_reset_tokens_pkey PRIMARY KEY (id);


--
-- Name: password_reset_tokens password_reset_tokens_token_unique; Type: CONSTRAINT; Schema: site; Owner: postgres
--

ALTER TABLE ONLY site.password_reset_tokens
    ADD CONSTRAINT password_reset_tokens_token_unique UNIQUE (token);


--
-- Name: complaints uniqcomplaints; Type: CONSTRAINT; Schema: site; Owner: postgres
--

ALTER TABLE ONLY site.complaints
    ADD CONSTRAINT uniqcomplaints UNIQUE (id);


--
-- Name: user_favorites user_favorites_pkey; Type: CONSTRAINT; Schema: site; Owner: postgres
--

ALTER TABLE ONLY site.user_favorites
    ADD CONSTRAINT user_favorites_pkey PRIMARY KEY (id);


--
-- Name: user_favorites user_favorites_unique; Type: CONSTRAINT; Schema: site; Owner: postgres
--

ALTER TABLE ONLY site.user_favorites
    ADD CONSTRAINT user_favorites_unique UNIQUE (user_id, listing_id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: site; Owner: postgres
--

ALTER TABLE ONLY site.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_password_reset_tokens_expires_at; Type: INDEX; Schema: site; Owner: postgres
--

CREATE INDEX idx_password_reset_tokens_expires_at ON site.password_reset_tokens USING btree (expires_at);


--
-- Name: idx_password_reset_tokens_token; Type: INDEX; Schema: site; Owner: postgres
--

CREATE INDEX idx_password_reset_tokens_token ON site.password_reset_tokens USING btree (token);


--
-- Name: categories categories_categories_fk; Type: FK CONSTRAINT; Schema: site; Owner: postgres
--

ALTER TABLE ONLY site.categories
    ADD CONSTRAINT categories_categories_fk FOREIGN KEY (parent_id) REFERENCES site.categories(id);


--
-- Name: chats chats_listings_fk; Type: FK CONSTRAINT; Schema: site; Owner: postgres
--

ALTER TABLE ONLY site.chats
    ADD CONSTRAINT chats_listings_fk FOREIGN KEY (listing_id) REFERENCES site.listings(id);


--
-- Name: complaints complaints_listings_fk; Type: FK CONSTRAINT; Schema: site; Owner: postgres
--

ALTER TABLE ONLY site.complaints
    ADD CONSTRAINT complaints_listings_fk FOREIGN KEY (listing_id) REFERENCES site.listings(id) ON DELETE CASCADE;


--
-- Name: complaints complaints_users_fk; Type: FK CONSTRAINT; Schema: site; Owner: postgres
--

ALTER TABLE ONLY site.complaints
    ADD CONSTRAINT complaints_users_fk FOREIGN KEY (complainant_id) REFERENCES site.users(id) ON DELETE CASCADE;


--
-- Name: email_confirmations fk_email_confirmation_user; Type: FK CONSTRAINT; Schema: site; Owner: postgres
--

ALTER TABLE ONLY site.email_confirmations
    ADD CONSTRAINT fk_email_confirmation_user FOREIGN KEY (user_id) REFERENCES site.users(id) ON DELETE CASCADE;


--
-- Name: listing_images listing_images_listings_fk; Type: FK CONSTRAINT; Schema: site; Owner: postgres
--

ALTER TABLE ONLY site.listing_images
    ADD CONSTRAINT listing_images_listings_fk FOREIGN KEY (listing_id) REFERENCES site.listings(id) ON DELETE CASCADE;


--
-- Name: listings listings_categories_fk; Type: FK CONSTRAINT; Schema: site; Owner: postgres
--

ALTER TABLE ONLY site.listings
    ADD CONSTRAINT listings_categories_fk FOREIGN KEY (category_id) REFERENCES site.categories(id);


--
-- Name: listings listings_locations_fk; Type: FK CONSTRAINT; Schema: site; Owner: postgres
--

ALTER TABLE ONLY site.listings
    ADD CONSTRAINT listings_locations_fk FOREIGN KEY (location_id) REFERENCES site.locations(id);


--
-- Name: listings listings_users_fk; Type: FK CONSTRAINT; Schema: site; Owner: postgres
--

ALTER TABLE ONLY site.listings
    ADD CONSTRAINT listings_users_fk FOREIGN KEY (user_id) REFERENCES site.users(id) ON DELETE CASCADE;


--
-- Name: messages messages_chats_fk; Type: FK CONSTRAINT; Schema: site; Owner: postgres
--

ALTER TABLE ONLY site.messages
    ADD CONSTRAINT messages_chats_fk FOREIGN KEY (chat_id) REFERENCES site.chats(id);


--
-- Name: password_reset_tokens password_reset_tokens_user_fk; Type: FK CONSTRAINT; Schema: site; Owner: postgres
--

ALTER TABLE ONLY site.password_reset_tokens
    ADD CONSTRAINT password_reset_tokens_user_fk FOREIGN KEY (user_id) REFERENCES site.users(id) ON DELETE CASCADE;


--
-- Name: user_favorites user_favorites_listing_id_fkey; Type: FK CONSTRAINT; Schema: site; Owner: postgres
--

ALTER TABLE ONLY site.user_favorites
    ADD CONSTRAINT user_favorites_listing_id_fkey FOREIGN KEY (listing_id) REFERENCES site.listings(id) ON DELETE CASCADE;


--
-- Name: user_favorites user_favorites_user_id_fkey; Type: FK CONSTRAINT; Schema: site; Owner: postgres
--

ALTER TABLE ONLY site.user_favorites
    ADD CONSTRAINT user_favorites_user_id_fkey FOREIGN KEY (user_id) REFERENCES site.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

\unrestrict vfDaLqC7kFxuknK9dVSke5kDzqBOdITaIoDPsRExxBbXSd6fjQxiQVhonTw4qcq

