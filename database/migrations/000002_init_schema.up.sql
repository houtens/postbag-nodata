CREATE TABLE "countries" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "name" varchar NOT NULL,
  "flag" varchar,
  "code" varchar,
  "priority" boolean NOT NULL DEFAULT false,
  "x_id" varchar
);

CREATE TABLE "clubs" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "name" varchar NOT NULL,
  "county" varchar,
  "website" varchar,
  "is_active" boolean NOT NULL DEFAULT false,
  "phone" varchar,
  "email" varchar,
  "contact_name" varchar,
  "country_id" uuid,
  "x_id" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "users" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "alt_name" varchar,
  "email" varchar,
  "password_hash" varchar,
  "absp_num" int,
  "club_id" uuid,
  "title_id" uuid,
  "role_id" uuid NOT NULL,
  "x_life" boolean NOT NULL,
  "x_post" boolean NOT NULL,
  "x_id" varchar NOT NULL
);

CREATE TABLE "titles" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "name" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "contacts" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "address1" varchar,
  "address2" varchar,
  "address3" varchar,
  "address4" varchar,
  "postcode" varchar,
  "country_id" uuid,
  "phone" varchar,
  "mobile" varchar,
  "user_id" uuid NOT NULL,
  "notes" varchar,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "tournament_state" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "name" varchar NOT NULL,
  "code" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "tournaments" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "name" varchar NOT NULL,
  "short_name" varchar,
  "start_date" timestamp,
  "end_date" timestamp,
  "state" uuid NOT NULL,
  "num_divisions" int,
  "num_rounds" int,
  "num_entries" int,
  "is_pc" boolean NOT NULL DEFAULT false,
  "is_fc" boolean NOT NULL DEFAULT false,
  "is_rr" boolean NOT NULL DEFAULT false,
  "is_wespa" boolean NOT NULL DEFAULT false,
  "is_invitational" boolean NOT NULL DEFAULT false,
  "is_locked" boolean NOT NULL DEFAULT false,
  "creator_id" uuid,
  "organiser_id" uuid,
  "director_id" uuid,
  "coperator_id" uuid,
  "x_id" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

-- Create once, can be refreshed from tournament/results data before issuing.
CREATE TABLE "invoices" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "tournament_id" uuid NOT NULL,
  "num_players" int NOT NULL DEFAULT 0,
  "num_non_members" int NOT NULL DEFAULT 0,
  "num_games" int NOT NULL DEFAULT 0,
  "is_multiday" boolean NOT NULL DEFAULT false,
  "is_overseas" boolean NOT NULL DEFAULT false,
  "levy_cost" real NOT NULL,
  "extras_cost" real NOT NULL,
  "total_cost" real NOT NULL,
  "is_paid" boolean NOT NULL DEFAULT false,
  -- description of the calculation
  "description" varchar,
  -- custom comment if necessary
  "extras_comment" varchar,
  "comment" varchar,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "ratings" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "user_id" uuid NOT NULL,
  "tournament_id" uuid NOT NULL,
  "division" int NOT NULL DEFAULT 1,
  "num_games" int,
  "start_rating" int,
  "end_rating" int,
  "rating_points" int,
  "opp_ratings_sum" int,
  "num_wins" real NOT NULL DEFAULT 0,
  "is_locked" boolean NOT NULL DEFAULT false,
  "x_id" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "results" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "player1_id" uuid,
  "player2_id" uuid,
  "score1" int NOT NULL,
  "score2" int NOT NULL,
  "spread" int NOT NULL,
  "tournament_id" uuid NOT NULL,
  "round_num" int NOT NULL,
  "type" int NOT NULL,
  "is_locked" boolean NOT NULL DEFAULT false,
  "x_id" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "payment_types" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "name" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

-- is_life may be redundant - life members have a valid (< expires_at) membership and years = 100
CREATE TABLE "membership_types" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "name" varchar NOT NULL,
  "code" varchar NOT NULL,
  "num_years" int not null,
  "is_junior" boolean not null default false,
  "is_post" boolean not null default false,
  "is_life" boolean not null default false,
  "created_at" timestamp not null default (now()),
  "updated_at" timestamp not null default (now())
);

create table "memberships" (
  "id" uuid primary key default (uuid_generate_v4()),
  "user_id" uuid not null,
  "cost" real not null,
  "membership_type_id" uuid not null,
  "payment_type_id" uuid not null,
  "expires_at" timestamp not null,
  "created_at" timestamp not null default (now()),
  "updated_at" timestamp not null default (now())
);

create table "auth_roles" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "name" varchar NOT NULL,
  "can_login" boolean NOT NULL DEFAULT false,
  "is_guest" boolean NOT NULL DEFAULT false,
  "is_members_admin" boolean NOT NULL DEFAULT false,
  "is_clubs_admin" boolean NOT NULL DEFAULT false,
  "is_ratings_admin" boolean NOT NULL DEFAULT false,
  "is_tournaments_admin" boolean NOT NULL DEFAULT false,
  "is_super_admin" boolean NOT NULL DEFAULT false,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);


ALTER TABLE "clubs" ADD FOREIGN KEY ("country_id") REFERENCES "countries" ("id");

ALTER TABLE "users" ADD FOREIGN KEY ("club_id") REFERENCES "clubs" ("id");
ALTER TABLE "users" ADD FOREIGN KEY ("title_id") REFERENCES "titles" ("id");
ALTER TABLE "users" ADD FOREIGN KEY ("role_id") REFERENCES "auth_roles" ("id");

ALTER TABLE "contacts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "tournaments" ADD FOREIGN KEY ("creator_id") REFERENCES "users" ("id");
ALTER TABLE "tournaments" ADD FOREIGN KEY ("organiser_id") REFERENCES "users" ("id");
ALTER TABLE "tournaments" ADD FOREIGN KEY ("director_id") REFERENCES "users" ("id");
ALTER TABLE "tournaments" ADD FOREIGN KEY ("coperator_id") REFERENCES "users" ("id");
ALTER TABLE "tournaments" ADD FOREIGN KEY ("state") REFERENCES "tournament_state" ("id");

ALTER TABLE "invoices" ADD FOREIGN KEY ("tournament_id") REFERENCES "tournaments" ("id");

ALTER TABLE "ratings" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "ratings" ADD FOREIGN KEY ("tournament_id") REFERENCES "tournaments" ("id");

ALTER TABLE "results" ADD FOREIGN KEY ("player1_id") REFERENCES "users" ("id");
ALTER TABLE "results" ADD FOREIGN KEY ("player2_id") REFERENCES "users" ("id");
ALTER TABLE "results" ADD FOREIGN KEY ("tournament_id") REFERENCES "tournaments" ("id");

ALTER TABLE "memberships" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "memberships" ADD FOREIGN KEY ("membership_type_id") REFERENCES "membership_types" ("id");
ALTER TABLE "memberships" ADD FOREIGN KEY ("payment_type_id") REFERENCES "payment_types" ("id");

