DROP PROCEDURE IF EXISTS `sp_delete_department`;
DELIMITER $$

CREATE PROCEDURE `sp_delete_department`(IN p_department_id BIGINT)
BEGIN
    -- === Declarations ===
    DECLARE v_parent_department_id BIGINT DEFAULT NULL;
    DECLARE v_not_found INT DEFAULT 0;

    -- Khi SELECT ... INTO không ra dòng nào sẽ kích hoạt NOT FOUND
    DECLARE CONTINUE HANDLER FOR NOT FOUND SET v_not_found = 1;

    -- Bất kỳ lỗi SQL nào khác -> rollback & bắn lại lỗi
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
BEGIN
ROLLBACK;
RESIGNAL;
END;

START TRANSACTION;

-- Khóa hàng để tránh race
SELECT group_department_id
INTO v_parent_department_id
FROM departments
WHERE id = p_department_id
  AND deleted_at IS NULL
    FOR UPDATE;

-- Không tìm thấy department hợp lệ
IF v_not_found = 1 THEN
        ROLLBACK;
        SIGNAL SQLSTATE '45000'
            SET MESSAGE_TEXT = 'Department not found or already deleted';
END IF;

    -- Không cho xóa root
    IF v_parent_department_id IS NULL THEN
        ROLLBACK;
        SIGNAL SQLSTATE '45000'
            SET MESSAGE_TEXT = 'You cannot delete the root department';
END IF;

    -- Re-assign con sang parent của node bị xóa
UPDATE departments
SET group_department_id = v_parent_department_id
WHERE group_department_id = p_department_id
  AND deleted_at IS NULL;

-- Soft delete (chỉ khi chưa xóa)
UPDATE departments
SET deleted_at = NOW()
WHERE id = p_department_id
  AND deleted_at IS NULL;

COMMIT;
END$$

DELIMITER ;
