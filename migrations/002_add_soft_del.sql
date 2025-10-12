-- =============================================
-- ADD SOFT DELETE SUPPORT
-- Add deleted_at column to all main tables
-- =============================================

-- Add deleted_at to companies
ALTER TABLE `companies`
    ADD COLUMN `deleted_at` TIMESTAMP NULL DEFAULT NULL AFTER `updated_at`;

-- Add deleted_at to business_units
ALTER TABLE `business_units`
    ADD COLUMN `deleted_at` TIMESTAMP NULL DEFAULT NULL AFTER `updated_at`;

-- Add deleted_at to departments
ALTER TABLE `departments`
    ADD COLUMN `deleted_at` TIMESTAMP NULL DEFAULT NULL AFTER `updated_at`;

-- Add deleted_at to starters
ALTER TABLE `starters`
    ADD COLUMN `deleted_at` TIMESTAMP NULL DEFAULT NULL AFTER `updated_at`;

-- =============================================
-- CREATE INDEXES FOR SOFT DELETE QUERIES
-- =============================================

-- Index for querying active (non-deleted) records
CREATE INDEX `idx_companies_deleted_at` ON `companies` (`deleted_at`);
CREATE INDEX `idx_business_units_deleted_at` ON `business_units` (`deleted_at`);
CREATE INDEX `idx_departments_deleted_at` ON `departments` (`deleted_at`);
CREATE INDEX `idx_starters_deleted_at` ON `starters` (`deleted_at`);

-- Composite indexes for common queries
CREATE INDEX `idx_starters_domain_deleted` ON `starters` (`domain`, `deleted_at`);
CREATE INDEX `idx_starters_dept_deleted` ON `starters` (`department_id`, `deleted_at`);
CREATE INDEX `idx_departments_bu_deleted` ON `departments` (`business_unit_id`, `deleted_at`);

-- =============================================
-- NOTES:
-- - deleted_at = NULL means record is active
-- - deleted_at = timestamp means record is soft deleted
-- - To query active records: WHERE deleted_at IS NULL
-- - To query deleted records: WHERE deleted_at IS NOT NULL
-- - To query all records: no WHERE clause on deleted_at
-- =============================================
