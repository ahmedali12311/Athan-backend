CREATE TABLE IF NOT EXISTS  daily_prayer_times (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    city_id uuid NOT NULL REFERENCES Cities(id) ON DELETE CASCADE,
    day INTEGER NOT NULL CHECK (day >= 1 AND day <= 31),
    month INTEGER NOT NULL CHECK (month >= 1 AND month <= 12),
    fajr_first_time TIME NOT NULL,
    fajr_second_time TIME NOT NULL,
    sunrise_time TIME NOT NULL,
    dhuhr_time TIME NOT NULL,
    asr_time TIME NOT NULL,
    maghrib_time TIME NOT NULL,
    isha_time TIME NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT prayer_times_day_month_city_id_key UNIQUE (day, month, city_id)
);


CREATE TRIGGER app_trigger_update_daily_prayer_Times_update_at
BEFORE UPDATE ON daily_prayer_times FOR EACH ROW
EXECUTE PROCEDURE app_func_update_updated_at();

