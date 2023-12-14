CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "hashedPassword" varchar NOT NULL,
  "fullName" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "passwordChangedAt" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "createdAt" timestamptz NOT NULL DEFAULT (now())
);

CREATE UNIQUE INDEX ON "accounts" ("owner", "currency");
-- ALTER TABLE "accounts" ADD CONSTRAINT "owner_currency_key" UNIQUE ("owner", "currency");

ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");
