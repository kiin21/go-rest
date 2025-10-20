DROP PROCEDURE IF EXISTS `sp_delete_department`;

DELIMITER
$$
CREATE PROCEDURE `sp_delete_department`(IN p_department_id BIGINT)
BEGIN
    DECLARE
v_parent_department_id BIGINT DEFAULT NULL;
    
    -- Handle case when SELECT returns no rows
    DECLARE
EXIT HANDLER FOR NOT FOUND
BEGIN
ROLLBACK;
SIGNAL
SQLSTATE '45000'
            SET MESSAGE_TEXT = 'Department not found or already deleted';
END;
  
    -- Rollback transaction if any SQL error occurs
    DECLARE
EXIT HANDLER FOR SQLEXCEPTION
BEGIN
ROLLBACK;
RESIGNAL;
END;

START TRANSACTION;

-- Get the parent department ID and lock the row
SELECT group_department_id
INTO v_parent_department_id
FROM departments
WHERE id = p_department_id
  AND deleted_at IS NULL
    FOR UPDATE;

-- Prevent deleting root
IF
v_parent_department_id IS NULL THEN
        ROLLBACK;
        SIGNAL
SQLSTATE '45000'
            SET MESSAGE_TEXT = 'You cannot delete the root department';
END IF;

    -- Reassign all child departments to the parent department
UPDATE departments
SET group_department_id = v_parent_department_id
WHERE group_department_id = p_department_id
  AND deleted_at IS NULL;

-- Soft delete
UPDATE departments
SET deleted_at = NOW()
WHERE id = p_department_id;

COMMIT;
END$$
