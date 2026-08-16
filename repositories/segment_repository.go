package repositories

import (
	"fmt"
	"gamebiller/helpers"
	"gamebiller/models"
	"strconv"
)

func CreateSegment(exec QueryExecutor, s *models.Segment) (int64, error) {
	query := `INSERT INTO segments (segment_name, agent_id, created_at, created_by, updated_at, updated_by)
	          VALUES (?, ?, ?, ?, ?, ?) RETURNING id`
	query = helpers.QuerySupport(query)
	var id int64
	err := exec.QueryRow(query, s.SegmentName, s.AgentID, s.CreatedAt, s.CreatedBy, s.UpdatedAt, s.UpdatedBy).Scan(&id)
	if err != nil {
		return 0, err
	}
	s.ID = id
	return id, nil
}

func GetSegmentByID(exec QueryExecutor, id int64) (*models.Segment, error) {
	query := `SELECT s.id, s.segment_name, s.agent_id, s.created_at, s.created_by, s.updated_at, s.updated_by,
	                 COALESCE(a.agent_name, '')
	          FROM segments s
	          LEFT JOIN agents a ON a.id = s.agent_id
	          WHERE s.id = $1`
	var s models.Segment
	err := exec.QueryRow(query, id).Scan(
		&s.ID, &s.SegmentName, &s.AgentID, &s.CreatedAt, &s.CreatedBy, &s.UpdatedAt, &s.UpdatedBy,
		&s.AgentName,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func GetSegmentByName(exec QueryExecutor, name string) (*models.Segment, error) {
	query := `SELECT s.id, s.segment_name, s.agent_id, s.created_at, s.created_by, s.updated_at, s.updated_by,
	                 COALESCE(a.agent_name, '')
	          FROM segments s
	          LEFT JOIN agents a ON a.id = s.agent_id
	          WHERE s.segment_name = $1 LIMIT 1`
	var s models.Segment
	err := exec.QueryRow(query, name).Scan(
		&s.ID, &s.SegmentName, &s.AgentID, &s.CreatedAt, &s.CreatedBy, &s.UpdatedAt, &s.UpdatedBy,
		&s.AgentName,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func UpdateSegment(exec QueryExecutor, s *models.Segment) error {
	query := `UPDATE segments SET segment_name = ?, agent_id = ?, updated_at = ?, updated_by = ? WHERE id = ?`
	query = helpers.QuerySupport(query)
	_, err := exec.Exec(query, s.SegmentName, s.AgentID, s.UpdatedAt, s.UpdatedBy, s.ID)
	return err
}

func DeleteSegment(exec QueryExecutor, id int64) error {
	query := `DELETE FROM segments WHERE id = $1`
	_, err := exec.Exec(query, id)
	return err
}

func GetSegmentsList(exec QueryExecutor, search string, start, length int, order, sort string, filters models.SegmentFilters) ([]models.Segment, int64, error) {
	var (
		count int64
		whr   string
	)

	if filters.ID != 0 {
		whr += " AND s.id = " + strconv.FormatInt(filters.ID, 10)
	}
	if filters.AgentID != 0 {
		whr += " AND s.agent_id = " + strconv.FormatInt(filters.AgentID, 10)
	}
	if filters.SegmentName != "" {
		whr += " AND s.segment_name ILIKE '%" + filters.SegmentName + "%'"
	}
	if search != "" {
		whr += " AND (s.segment_name ILIKE '%" + search + "%')"
	}

	countQuery := `SELECT COUNT(*) FROM segments s LEFT JOIN agents a ON a.id = s.agent_id WHERE true` + whr
	err := exec.QueryRow(countQuery).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	orderBy := "s.id"
	if order != "" {
		orderBy = order
	}
	sortOrder := "DESC"
	if sort != "" {
		sortOrder = sort
	}

	query := `SELECT s.id, s.segment_name, s.agent_id, s.created_at, s.created_by, s.updated_at, s.updated_by,
	                 COALESCE(a.agent_name, '')
	          FROM segments s
	          LEFT JOIN agents a ON a.id = s.agent_id
	          WHERE true` + whr + " ORDER BY " + orderBy + " " + sortOrder
	if length > 0 {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", length, start)
	}
	query = helpers.QuerySupport(query)

	rows, err := exec.Query(query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.Segment
	for rows.Next() {
		var s models.Segment
		err = rows.Scan(
			&s.ID, &s.SegmentName, &s.AgentID, &s.CreatedAt, &s.CreatedBy, &s.UpdatedAt, &s.UpdatedBy,
			&s.AgentName,
		)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, s)
	}
	return list, count, nil
}
