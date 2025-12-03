
🕌 Prayer Notification System (Athan Project)
The backend includes a sophisticated, schedule-driven system to send timely and hyper-localized push notifications for daily prayer times (Athan). This feature ensures high accuracy and minimal latency by integrating directly with the application scheduler and using PostgreSQL's time functionality.

Key Features
Minute-Level Precision: The core job runs every minute to check against the scheduled prayer times, guaranteeing notifications are sent precisely when the Athan is due.

Topic-Based Targeting (Anonymous Users): Notifications are sent using FCM Topics derived from the city and prayer name (e.g., Tripoli_fajr, Benghazi_asr). This allows for highly scalable messaging that reaches all anonymous subscribers without managing individual user tokens.

Dynamic, Localized Messaging: The system dynamically constructs the Arabic notification message based on the prayer time that is due:

Title: Formatted as "حان الآن موعد صلاة [Prayer Name]" (e.g., “The time for Asr prayer is now.”).

Body: Includes the specific prayer time and city (e.g., “العصر (17:05) في طرابلس”).

Data Models: This system introduced specific PostgreSQL tables to manage the required data:

daily_prayer_times: Stores the exact time for each of the five daily prayers for a given city, day, and month.

adhkars, hadiths, special_topics, categories: Provides the content foundation for future features related to Islamic content and reminders.

📜 Database Schemas Overview
The template includes robust, production-ready schemas for core content and features:

daily_prayer_times: The central table for the Athan scheduler, tracking daily prayer schedules by city.

categories: A complex, hierarchical model supporting multi-level categorization with automated path tracking for fast retrieval.

hadiths & adhkars: Structured content tables for religious texts and supplications, linked to the categories model.

special_topics: Content for general or non-categorized announcements and information.

---

- [Changelog](CHANGELOG.md)

---

## resources

- [go.dev](https://go.dev/) Go Programming Language
- [pgx](https://github.com/jackc/pgx) Replaces the standard `database/sql` driver
- [sqlx](https://github.com/jmoiron/sqlx) Extends the standard `database/sql` functions
- [echo](https://github.com/labstack/echo) minimalist web framework
- [goi18n](https://github.com/nicksnyder/go-i18n) internationalization
- [squirrel](https://github.com/Masterminds/squirrel) Query Builder
- [jsonSchema](https://github.com/santhosh-tekuri/jsonschema) Json Schema validation

---

### Lab

- `make lab.down && make lab` removes the old container and builds the new one

---

## Makefile

this project utilises docker to run `builds` and `migration`

`make [command]`

### Commands

- `init` install go development dependencies
- `build` build binary
- `run` run built binary
- `test` run tests
- `dev` build a docker image on local machine
- `dev.down` stops and remove dev docker image
- `migrate.up n=1` migrate database `n` steps
- `migrate.up.all` migrate database to latest
- `migrate.down n=1` rolls back `n` migrations
- `migrate.down.all` rolls back all migration
- `migration n=create_somethings_table` creates up and down sql in migrations
- `migrate.force n=23` force back failed migration version
- `refresh` runs down.all + up + seed
- `prune` prunes unused volumes, images and build caches
- `docker.ps` better format for docker ps command
- `audit` runs audit with go utilities on the project
- `list` lists update-able dependencies
- `update` downloads and updates project dependencies
- `swag` format and generate swag docs
- translations:
    - `translate.extract` update the `active.en.toml` file
    - `translate.merge` creates the `translate.ar.toml` file with new variables
    - translate the content of `translate.ar.toml` values
    - `translate.merge.done` merges translations to the `active.ar.toml` file

### Notes

installing psql on mac without starting the service:

1. `brew search postgres`
2. `brew install postgresql`
3. `echo 'export PATH="/opt/homebrew/opt/postgresql@16/bin:$PATH"' >> ~/.zshrc`
