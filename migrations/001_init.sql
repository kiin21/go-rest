-- =============================================
-- SCHEMA
-- =============================================

CREATE TABLE IF NOT EXISTS `companies` (
                                           `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
                                           `name` VARCHAR(255) NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    ) ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4
    COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `business_units` (
                                                `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
                                                `name` VARCHAR(255) NOT NULL,
    `shortname` VARCHAR(50) NOT NULL,
    `company_id` BIGINT NOT NULL,
    `leader_id` BIGINT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    ) ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4
    COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `departments` (
                                             `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
                                             `group_department_id` BIGINT NULL,
                                             `full_name` VARCHAR(255) NOT NULL,
    `shortname` VARCHAR(100) NOT NULL,
    `business_unit_id` BIGINT NULL,
    `leader_id` BIGINT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    ) ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4
    COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `starters` (
                                          `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
                                          `domain` VARCHAR(25) NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `email` VARCHAR(100) NOT NULL,
    `mobile` VARCHAR(20) NOT NULL,
    `work_phone` VARCHAR(20) NOT NULL,
    `job_title` VARCHAR(100) NOT NULL,
    `department_id` BIGINT NULL,
    `line_manager_id` BIGINT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    ) ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4
    COLLATE = utf8mb4_unicode_ci;

-- =============================================
-- FOREIGN KEYS
-- =============================================

ALTER TABLE `business_units`
    ADD CONSTRAINT `fk_bu_company`
        FOREIGN KEY (`company_id`) REFERENCES `companies` (`id`)
            ON DELETE RESTRICT
            ON UPDATE CASCADE,
    ADD CONSTRAINT `fk_bu_leader`
        FOREIGN KEY (`leader_id`) REFERENCES `starters` (`id`)
            ON DELETE SET NULL
            ON UPDATE CASCADE;

ALTER TABLE `departments`
    ADD CONSTRAINT `fk_dep_bu`
        FOREIGN KEY (`business_unit_id`) REFERENCES `business_units` (`id`)
            ON DELETE SET NULL
            ON UPDATE CASCADE,
    ADD CONSTRAINT `fk_dep_group`
        FOREIGN KEY (`group_department_id`) REFERENCES `departments` (`id`)
            ON DELETE SET NULL
            ON UPDATE CASCADE,
                                                                      ADD CONSTRAINT `fk_dep_leader`
                                                                      FOREIGN KEY (`leader_id`) REFERENCES `starters` (`id`)
               ON DELETE SET NULL
            ON UPDATE CASCADE;

ALTER TABLE `starters`
    ADD CONSTRAINT `fk_starter_department`
        FOREIGN KEY (`department_id`) REFERENCES `departments` (`id`)
            ON DELETE SET NULL
            ON UPDATE CASCADE,
    ADD CONSTRAINT `fk_starter_manager`
        FOREIGN KEY (`line_manager_id`) REFERENCES `starters` (`id`)
            ON DELETE SET NULL
            ON UPDATE CASCADE;

-- =============================================
-- INSERT DATA
-- =============================================

-- Insert company
INSERT INTO companies (id, name)
VALUES (1, 'VNG Corporation')
    ON DUPLICATE KEY UPDATE
                         name = VALUES(name);

-- Insert leader starters first (without department and line manager initially)
INSERT INTO starters (id, domain, name, email, mobile, work_phone, job_title, department_id, line_manager_id)
VALUES
    (1, 'minhlh', 'Le Hoang Minh', 'minhlh@vng.com.vn', '(+84) 0913875329', '0913875329', 'Founder, Chairman of VNG', NULL, NULL),
    (2, 'kelly', 'Kelly Le', 'kelly@vng.com.vn', '(+84) 0919392288', '0919392288', 'Chief Executive Officer', NULL, NULL),
    (3, 'chill', 'Chi Le', 'chill@vng.com.vn', '(+84) 08607548666', '08607548666', 'CEO of Zalopay', NULL, NULL),
    (4, 'thanhnl', 'Nguyen Le Thanh', 'thanhnl@vng.com.vn', '(+84) 0764450920', '0764450920 (6200)', 'Vice President of VNG, CEO of Digital Business', NULL, NULL),
    (5, 'khaivq', 'Vu Quang Khai', 'khaivq@vng.com.vn', '(+84) 0904101242', '0904101242 (1001)', 'Co-founder, Executive Vice President of VNG', NULL, NULL),
    (6, 'hawkinsp', 'Pham Hawkins', 'hawkinsp@vng.com.vn', '(+84) 0937972323', '0937972323', 'Director, Head of Publishing Platform', NULL, NULL),
    (7, 'vietph', 'Pham Hoang Viet', 'vietph@vng.com.vn', '(+84) 0909987790', '0909987790', 'Group Manager, Head of Engineering', NULL, NULL),
    (8, 'locth2', 'Tran Hoang Loc', 'locth2@vng.com.vn', '(+84) 0934482999', '0934482999 (35500)', 'Technical Manager', NULL, NULL),
    (9, 'duydh2', 'Dang Huu Duy', 'duydh2@vng.com.vn', '(+84) 0904157432', '0904157432', 'Lead Software Engineer', NULL, NULL),
    (10, 'khoapda', 'Pham Dinh Anh Khoa', 'khoapda@vng.com.vn', '(+84) 0359754697', '0359754697', 'Software Intern', NULL, NULL),
    (11, 'duyct', 'Cao Tien Duy', 'duyct@vng.com.vn', '(+84) 0908608731', '0908608731', 'Senior Data Engineer', NULL, NULL),
    (12, 'trungnk', 'Nguyen Kim Trung', 'trungnk@vng.com.vn', '(+84) 0908226847', '0908226847', 'Senior Data Engineer', NULL, NULL),
    (13, 'trungnk', 'Nguyen Kim Trung', 'trungnk@vng.com.vn', '(+84) 0908226847', '0908226847', 'Senior Data Engineer', NULL, NULL)
    ON DUPLICATE KEY UPDATE
                         domain = VALUES(domain);

-- Insert business units (without leader first to avoid FK constraint)
INSERT INTO business_units (id, name, shortname, company_id, leader_id)
VALUES
    (1, 'VNGGames', 'VNGG', 1, NULL),
    (2, 'Zalo', 'ZALO', 1, NULL),
    (3, 'Zalopay', 'ZLP', 1, NULL),
    (4, 'Digital Business', 'DB', 1, NULL)
    ON DUPLICATE KEY UPDATE
                         name = VALUES(name),
                         shortname = VALUES(shortname);

-- Insert departments (without leader first)
INSERT INTO departments (id, group_department_id, business_unit_id, full_name, shortname, leader_id)
VALUES
    (1, NULL, NULL, 'Senior Management Team', 'SMT', NULL),
    (2, NULL, NULL, 'Management Team', 'MT', NULL),
    (3, NULL, NULL, 'Executive Services', 'ES', NULL),
    (4, NULL, 3, 'Zalopay Management Team', 'ZMT', NULL),
    (5, NULL, 1, 'Games Publishing Platform', 'GPP', NULL),
    (6, NULL, 1, 'Games Studios', 'GSs', NULL),
    (7, NULL, 1, 'Games Developments', 'GDs', NULL),
    (8, 5, NULL, 'Product Core', 'PRO', NULL),
    (9, 5, NULL, 'Publishing Platform Engineering', 'PEN', NULL),
    (10, 5, NULL, 'Platform Integration', 'PIN', NULL),
    (11, 5, NULL, 'Game Infrastructure Operation', 'GIO', NULL),
    (12, 5, NULL, 'Game Data Studio', 'GDS', NULL)
    ON DUPLICATE KEY UPDATE
                         full_name = VALUES(full_name),
                         shortname = VALUES(shortname);

-- Update department assignments for starters
UPDATE starters
SET department_id = CASE
                        WHEN domain = 'minhlh' THEN 1
    WHEN domain = 'kelly' THEN 1
    WHEN domain = 'chill' THEN 4
    WHEN domain = 'thanhnl' THEN 2
    WHEN domain = 'khaivq' THEN 2
    WHEN domain = 'hawkinsp' THEN 5
    WHEN domain = 'vietph' THEN 9
    WHEN domain = 'locth2' THEN 9
    WHEN domain = 'duydh2' THEN 9
    WHEN domain = 'khoapda' THEN 9
    WHEN domain = 'duyct' THEN 12
    WHEN domain = 'trungnk' THEN 4
END
WHERE domain IN ('minhlh', 'kelly', 'chill', 'thanhnl', 'khaivq', 'hawkinsp', 'vietph', 'locth2', 'duydh2', 'khoapda', 'duyct', 'trungnk');

-- Update line managers for starters
UPDATE starters
SET line_manager_id = CASE
                          WHEN domain = 'kelly' THEN 1
    WHEN domain = 'chill' THEN 2
    WHEN domain = 'thanhnl' THEN 1
    WHEN domain = 'khaivq' THEN 1
    WHEN domain = 'hawkinsp' THEN 2
    WHEN domain = 'vietph' THEN 6
    WHEN domain = 'locth2' THEN 7
    WHEN domain = 'duydh2' THEN 8
    WHEN domain = 'khoapda' THEN 9
    WHEN domain = 'duyct' THEN 6
    WHEN domain = 'trungnk' THEN 1
END
WHERE domain IN ('kelly', 'chill', 'thanhnl', 'khaivq', 'hawkinsp', 'vietph', 'locth2', 'duydh2', 'khoapda', 'duyct', 'trungnk');

-- Update leaders for departments
UPDATE departments
SET leader_id = CASE
                    WHEN id = 1 THEN 1   -- SMT -> minhlh
                    WHEN id = 2 THEN 1   -- MT -> minhlh
                    WHEN id = 3 THEN 1   -- ES -> minhlh
                    WHEN id = 4 THEN 3   -- ZMT -> chill
                    WHEN id = 5 THEN 6   -- GPP -> hawkinsp
                    WHEN id = 6 THEN 2   -- GSs -> kelly
                    WHEN id = 7 THEN 2   -- GDs -> kelly
                    WHEN id = 8 THEN 6   -- PRO -> hawkinsp
                    WHEN id = 9 THEN 7   -- PEN -> vietph
                    WHEN id = 10 THEN 6  -- PIN -> hawkinsp
                    WHEN id = 11 THEN 6  -- GIO -> hawkinsp
                    WHEN id = 12 THEN 11 -- GDS -> duyct
    END
WHERE id IN (1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12);

-- Update leaders for business units
UPDATE business_units
SET leader_id = CASE
                    WHEN id = 1 THEN 2 -- VNGGames -> kelly
                    WHEN id = 2 THEN 5 -- Zalo -> khaivq
                    WHEN id = 3 THEN 3 -- Zalopay -> chill
                    WHEN id = 4 THEN 4 -- Digital Business -> thanhnl
    END
WHERE id IN (1, 2, 3, 4);

-- =============================================
-- INDEXES
-- =============================================

-- Basic indexes will be created here
-- Index for deleted_at will be created in 003_add_soft_delete.sql after the column is added