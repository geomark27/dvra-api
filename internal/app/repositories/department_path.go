package repositories

import (
	"dvra-api/internal/app/models"

	"gorm.io/gorm"
)

type departmentPathRow struct {
	LeafID uint
	ID     uint
	Name   string
	Code   string
	Depth  int
}

// findDepartmentPaths resolves all requested paths in one PostgreSQL recursive
// query. Results are ordered root-to-leaf and never cross tenant boundaries.
func findDepartmentPaths(db *gorm.DB, departmentIDs []uint) (map[uint][]models.DepartmentPathItem, error) {
	paths := make(map[uint][]models.DepartmentPathItem, len(departmentIDs))
	if len(departmentIDs) == 0 {
		return paths, nil
	}
	var rows []departmentPathRow
	err := db.Raw(`
		WITH RECURSIVE department_ancestors AS (
			SELECT d.id AS leaf_id, d.company_id, d.id, d.parent_department_id,
				d.name, d.code, 0 AS depth, ARRAY[d.id]::bigint[] AS visited
			FROM departments AS d
			WHERE d.id IN ? AND d.deleted_at IS NULL

			UNION ALL

			SELECT child.leaf_id, parent.company_id, parent.id, parent.parent_department_id,
				parent.name, parent.code, child.depth + 1,
				child.visited || parent.id
			FROM departments AS parent
			JOIN department_ancestors AS child
				ON child.parent_department_id = parent.id
				AND child.company_id = parent.company_id
			WHERE parent.deleted_at IS NULL
				AND NOT parent.id = ANY(child.visited)
		)
		SELECT leaf_id, id, name, code, depth
		FROM department_ancestors
		ORDER BY leaf_id, depth DESC`, departmentIDs).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		paths[row.LeafID] = append(paths[row.LeafID], models.DepartmentPathItem{
			ID: row.ID, Name: row.Name, Code: row.Code,
		})
	}
	return paths, nil
}

func attachDepartmentPaths(departments []models.Department, paths map[uint][]models.DepartmentPathItem) {
	for i := range departments {
		departments[i].Path = paths[departments[i].ID]
	}
}

func attachJobDepartmentPaths(db *gorm.DB, jobs []models.Job) error {
	departmentIDs := make([]uint, 0, len(jobs))
	seen := make(map[uint]bool)
	for i := range jobs {
		if jobs[i].Department == nil || seen[jobs[i].Department.ID] {
			continue
		}
		seen[jobs[i].Department.ID] = true
		departmentIDs = append(departmentIDs, jobs[i].Department.ID)
	}
	paths, err := findDepartmentPaths(db, departmentIDs)
	if err != nil {
		return err
	}
	for i := range jobs {
		if jobs[i].Department != nil {
			jobs[i].Department.Path = paths[jobs[i].Department.ID]
		}
	}
	return nil
}
