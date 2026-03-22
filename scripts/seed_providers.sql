INSERT INTO providers (id, name, image, created_at, updated_at) VALUES
  (gen_random_uuid(), 'Crunchyroll',        'https://www.crunchyroll.com/build/assets/img/favicons/favicon-192x192.png',                    NOW(), NOW()),
  (gen_random_uuid(), 'Netflix',            'https://assets.nflxext.com/us/ffe/siteui/common/icons/nficon2016.png',                          NOW(), NOW()),
  (gen_random_uuid(), 'Funimation',         'https://www.funimation.com/favicon.ico',                                                         NOW(), NOW()),
  (gen_random_uuid(), 'HiDive',             'https://www.hidive.com/favicon.ico',                                                             NOW(), NOW()),
  (gen_random_uuid(), 'Amazon Prime Video', 'https://m.media-amazon.com/images/G/01/primevideo/seo/primevideo-seo-logo.png',                   NOW(), NOW()),
  (gen_random_uuid(), 'Disney+',            'https://cnbl-cdn.bamgrid.com/assets/7ecc8bcb60ad77193058d63e321bd21cbac2fc67/original',          NOW(), NOW()),
  (gen_random_uuid(), 'Hulu',               'https://assetshuluimcom-a.akamaihd.net/h5/default_v3/static/hulu-logo.png',                      NOW(), NOW()),
  (gen_random_uuid(), 'Apple TV+',          'https://tv.apple.com/assets/atv-web/atvweb.png',                                                 NOW(), NOW()),
  (gen_random_uuid(), 'Bilibili',           'https://www.bilibili.com/favicon.ico',                                                           NOW(), NOW()),
  (gen_random_uuid(), 'Aniplus',            'https://www.aniplus-asia.com/favicon.ico',                                                       NOW(), NOW())
ON CONFLICT DO NOTHING;
