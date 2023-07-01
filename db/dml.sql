-- =============================================
-- Author:      William Wibowo Ciptono
-- Create date: 22 Jan 2023
-- Description: data seeding for blanche DB
-- =============================================


INSERT INTO roles (
id, 
role_name
)VALUES 
    (1, 'user'),
    (2, 'admin')
;

INSERT INTO users (
role_id,
username,
email,
password,
wallet_pin
)VALUES
(1, 'user', 'user@mail.com', 'user', '123456'),
(2, 'admin', 'admin@mail.com', 'admin', '123456');
-- password & wallet_pin not encrpyted yet

INSERT INTO user_details (
user_id,
fullname,
phone,
gender,
birthdate,
profile_picture
)VALUES
(1, 'user user', '0572342934', 'male', NULL, NULL),
(2, 'admin admin', '0572342935', 'female', NULL, NULL);

INSERT INTO user_addresses(
user_id,
province,
city,
sub_district,
zip_code,
label,
details,
name,
phone_number,
is_default
)VALUES
(1, 'Jawa Tengah', 'Semarang', 'Pedurungan', 50192, 'rumah', '', 'will', '085375627432',false),
(1, 'Jawa Tengah', 'Semarang', 'Pedurungan', 50192, 'rumah', '', 'kris', '0234823842',true),
(2, 'Jawa Tengah', 'Semarang', 'Pedurungan', 50192, 'rumah', '', 'will', '085375627432',true);
