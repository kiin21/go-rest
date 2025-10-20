-- =============================================
-- VIEW
-- =============================================
CREATE
OR REPLACE VIEW v_departments_with_bu AS
WITH RECURSIVE dept_hierarchy AS (SELECT id,
                                         group_department_id,
                                         full_name,
                                         shortname,
                                         leader_id,
                                         business_unit_id,
                                         created_at,
                                         updated_at,
                                         deleted_at,
                                         business_unit_id AS actual_business_unit_id,
                                         0                AS level
                                  FROM departments
                                  WHERE group_department_id IS NULL
                                    AND deleted_at IS NULL

                                  UNION ALL

                                  SELECT d.id,
                                         d.group_department_id,
                                         d.full_name,
                                         d.shortname,
                                         d.leader_id,
                                         d.business_unit_id,
                                         d.created_at,
                                         d.updated_at,
                                         d.deleted_at,
                                         dh.actual_business_unit_id, -- Inherit từ parent
                                         dh.level + 1
                                  FROM departments d
                                           INNER JOIN dept_hierarchy dh ON d.group_department_id = dh.id
                                  WHERE d.deleted_at IS NULL)
SELECT id,
       group_department_id,
       full_name,
       shortname,
       leader_id,
       created_at,
       updated_at,
       deleted_at,
       actual_business_unit_id AS business_unit_id
FROM dept_hierarchy;

-- =============================================
-- VIEW
-- =============================================
CREATE
OR REPLACE VIEW v_departments_with_counts AS
SELECT d.id,
       d.group_department_id,
       d.full_name,
       d.shortname,
       d.leader_id,
       d.business_unit_id,
       d.created_at,
       d.updated_at,
       d.deleted_at,
       COUNT(DISTINCT s.id)  AS total_starters,
       COUNT(DISTINCT sd.id) AS total_subdepartments
FROM departments d
         LEFT JOIN starters s ON s.department_id = d.id AND s.deleted_at IS NULL
         LEFT JOIN departments sd ON sd.group_department_id = d.id AND sd.deleted_at IS NULL
WHERE d.deleted_at IS NULL
GROUP BY d.id,
         d.group_department_id,
         d.full_name,
         d.shortname,
         d.leader_id,
         d.business_unit_id,
         d.created_at,
         d.updated_at,
         d.deleted_at;
