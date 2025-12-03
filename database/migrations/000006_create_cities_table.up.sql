CREATE TABLE IF NOT EXISTS cities (
    id          UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    name        TEXT NOT NULL,
    is_disabled BOOLEAN NOT NULL DEFAULT FALSE,
    location    GEOMETRY (POINT, 4326),
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    UNIQUE (name)
);
CREATE OR REPLACE FUNCTION app_func_update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER app_trigger_update_cities_update_at
BEFORE UPDATE ON cities FOR EACH ROW
EXECUTE PROCEDURE app_func_update_updated_at();

INSERT INTO cities
(id, name, location)
VALUES
('3bc7f707-1a4d-48c1-a5bd-0d243b655aa7', 'بنغازي', (SELECT st_geomfromgeojson('{"type":"Point","coordinates":[20.0881772, 32.1189829]}'))),
('adf0d179-2967-4bbe-a7a3-df2ae94d0362', 'طرابلس', (SELECT st_geomfromgeojson('{"type":"Point","coordinates":[13.1887656, 32.8876938]}'))),
('1e77a009-8d9b-4253-9b0d-ab3752b3209c', 'الجغبوب', (SELECT st_geomfromgeojson('{"type":"Point","coordinates":[24.5156454, 29.7431618]}'))),
('2cbea011-bef4-4c00-b90e-be285327750a', 'البريقة', (SELECT st_geomfromgeojson('{"type":"Point","coordinates":[19.5692142, 30.405783]}'))),
('30b3dbb8-16d7-4349-be33-74ce37cdddc7', 'مصراته', (SELECT st_geomfromgeojson('{"type":"Point","coordinates":[15.0940628, 32.3265127]}'))),
('3894a2bb-cd07-49f5-888e-62e0096b51fc', 'إغدامس', (SELECT st_geomfromgeojson('{"type":"Point","coordinates":[9.4938173, 30.1301998]}'))),
('07d0521a-9d88-4005-849a-c343028fac66', 'لانوف', (SELECT st_geomfromgeojson('{"type":"Point","coordinates":[18.5494458, 30.5076417]}'))),
('3c2ae040-7cd2-4a59-96bb-e9c3211973fb', 'جالو', (SELECT st_geomfromgeojson('{"type":"Point","coordinates":[21.4941389, 29.0411565]}'))),
('4136e239-3c7e-4f49-a5dc-bf1e1a674749', 'اوباري', (SELECT st_geomfromgeojson('{"type":"Point","coordinates":[12.7907482, 26.581723]}'))),
('43ab7849-1b4c-4a57-b010-e23c94c7a126', 'سرت', (SELECT st_geomfromgeojson('{"type":"Point","coordinates":[16.5664641, 31.1887704]}'))),
('4ad91a99-a4a5-4e98-aca6-019dfff37c96', 'إجدابيا', (SELECT st_geomfromgeojson('{"type":"Point","coordinates":[20.1059132, 30.2133625]}'))),
('5d4c7286-7e2a-45b9-8b50-0be113e2976f', 'هون', (SELECT st_geomfromgeojson('{"type":"Point","coordinates":[15.9409452, 28.6662964]}'))),
('6d500e9a-4a47-43b9-acb0-1ca4c0d0f015', 'البيضاء', (SELECT st_geomfromgeojson('{"type":"Point","coordinates":[21.7342299, 32.7572094]}'))),
('6dc8f605-4092-40eb-aa6d-a373c35942cc', 'زليتن', (SELECT st_geomfromgeojson('{"type":"Point","coordinates":[14.4449538, 32.2850161]}'))),
('6e60b2b6-5c94-46a6-a151-7209f2f09335', 'المرج', (SELECT st_geomfromgeojson('{"type":"Point","coordinates":[20.8173194, 32.4970411]}'))),
('78557158-7fa6-4368-b2cf-3757766dd9a6', 'الزاوية', (SELECT st_geomfromgeojson('{"type":"Point","coordinates":[12.7303464, 32.7662368]}'))),
('7fdc0002-7d35-49cd-ae8a-4eacfa4608b1', 'الخمس', (SELECT st_geomfromgeojson('{"type":"Point","coordinates":[14.2610275, 32.6486679]}'))),
('a6dbabe6-3379-4dff-8327-f8cc3caf9305', 'امساعد', (SELECT st_geomfromgeojson('{"type":"Point","coordinates":[25.0457125, 31.605082]}'))),
('12e9a871-819f-4821-98d1-7f4583d4898b', 'سبها', (SELECT st_geomfromgeojson('{"type":"Point","coordinates":[14.4243868, 27.0360739]}'))),
('b34c7c81-c72b-496d-82d1-57115f65d993', 'طبرق', (SELECT st_geomfromgeojson('{"type":"Point","coordinates":[23.9406318, 32.0690011]}'))),
('b3888d04-47be-4cc1-81b5-90e7374efb0d', 'بني وليد', (SELECT st_geomfromgeojson('{"type":"Point","coordinates":[13.9775257, 31.7602015]}'))),
('b81cb64e-2db3-4602-a97e-78e1741717b1', 'الكفرة', (SELECT st_geomfromgeojson('{"type":"Point","coordinates":[23.2598148, 24.1325483]}'))),
('ca2eb404-4a61-4284-9221-2eb5bd427d09', 'شحات', (SELECT st_geomfromgeojson('{"type":"Point","coordinates":[21.8592076, 32.8049799]}'))),
('d1170490-368b-47f9-9520-ca8c37704bdb', 'درنة', (SELECT st_geomfromgeojson('{"type":"Point","coordinates":[22.6364145, 32.757385]}'))),
('46e594c7-1d5b-4e8a-8f0a-6f4e3c9d2a1b', 'مرزق', (SELECT st_geomfromgeojson('{"type":"Point","coordinates":[25.92237457490026, 13.92566401344491]}'))),
('ee4f03bd-d271-4ddf-ae5f-857e1e00c647', 'اوجلة', (SELECT st_geomfromgeojson('{"type":"Point","coordinates":[21.2911149, 29.1326919]}')));
