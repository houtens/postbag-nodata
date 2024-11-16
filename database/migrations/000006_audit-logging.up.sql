CREATE TABLE "audit_log" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "event" int not null,
  "user_id" uuid not null,
  "remote_ip" varchar not null,
  "user_agent" varchar not null,
  "message" varchar not null,
  "created_at" timestamp NOT NULL DEFAULT (now())
);
