package schedule

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, schedule Schedule) (Schedule, error) {
	err := r.db.QueryRow(ctx, `
SELECT branch_id
FROM groups
WHERE id = $1
  AND organization_id = $2
  AND status = 'active'
`, schedule.GroupID, schedule.OrganizationID).Scan(&schedule.BranchID)
	if err != nil {
		return Schedule{}, ErrGroupNotFound
	}

	err = r.db.QueryRow(ctx, `
INSERT INTO schedules (
id,
organization_id,
branch_id,
group_id,
weekday,
start_time,
end_time,
room
)
VALUES ($1, $2, $3, $4, $5, $6::time, $7::time, $8)
RETURNING
id,
organization_id,
branch_id,
group_id,
weekday,
start_time::text,
end_time::text,
COALESCE(room, '')
`,
		schedule.ID,
		schedule.OrganizationID,
		schedule.BranchID,
		schedule.GroupID,
		schedule.Weekday,
		schedule.StartTime,
		schedule.EndTime,
		schedule.Room,
	).Scan(
		&schedule.ID,
		&schedule.OrganizationID,
		&schedule.BranchID,
		&schedule.GroupID,
		&schedule.Weekday,
		&schedule.StartTime,
		&schedule.EndTime,
		&schedule.Room,
	)
	if err != nil {
		return Schedule{}, err
	}

	return schedule, nil
}

func (r *Repository) ListByOrganizationID(ctx context.Context, organizationID uuid.UUID) ([]Schedule, error) {
	rows, err := r.db.Query(ctx, `
SELECT
id,
organization_id,
branch_id,
group_id,
weekday,
start_time::text,
end_time::text,
COALESCE(room, '')
FROM schedules
WHERE organization_id = $1
ORDER BY weekday ASC, start_time ASC
`, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	schedules := make([]Schedule, 0)

	for rows.Next() {
		var item Schedule

		if err := rows.Scan(
			&item.ID,
			&item.OrganizationID,
			&item.BranchID,
			&item.GroupID,
			&item.Weekday,
			&item.StartTime,
			&item.EndTime,
			&item.Room,
		); err != nil {
			return nil, err
		}

		schedules = append(schedules, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return schedules, nil
}

func (r *Repository) GenerateLessons(
	ctx context.Context,
	organizationID uuid.UUID,
	scheduleID uuid.UUID,
	fromDate time.Time,
	toDate time.Time,
	topic string,
) ([]GeneratedLesson, []string, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var schedule Schedule

	err = tx.QueryRow(ctx, `
SELECT
s.id,
s.organization_id,
s.branch_id,
s.group_id,
s.weekday,
s.start_time::text,
s.end_time::text,
COALESCE(s.room, ''),
g.subject_id,
COALESCE(g.teacher_id::text, '')
FROM schedules s
JOIN groups g ON g.id = s.group_id
WHERE s.id = $1
  AND s.organization_id = $2
  AND g.status = 'active'
`, scheduleID, organizationID).Scan(
		&schedule.ID,
		&schedule.OrganizationID,
		&schedule.BranchID,
		&schedule.GroupID,
		&schedule.Weekday,
		&schedule.StartTime,
		&schedule.EndTime,
		&schedule.Room,
		&schedule.SubjectID,
		&schedule.TeacherID,
	)
	if err != nil {
		return nil, nil, ErrScheduleNotFound
	}

	createdLessons := make([]GeneratedLesson, 0)
	skippedDates := make([]string, 0)

	for currentDate := fromDate; !currentDate.After(toDate); currentDate = currentDate.AddDate(0, 0, 1) {
		if weekdayNumber(currentDate) != schedule.Weekday {
			continue
		}

		lessonDate := currentDate.Format("2006-01-02")
		lessonID := uuid.New()

		var teacherID interface{}
		if schedule.TeacherID != "" {
			teacherID = schedule.TeacherID
		}

		var created GeneratedLesson

		err = tx.QueryRow(ctx, `
INSERT INTO lessons (
	id,
	organization_id,
	branch_id,
	group_id,
	teacher_id,
	planned_teacher_id,
	actual_teacher_id,
	subject_id,
	schedule_id,
	lesson_date,
	start_time,
	end_time,
	topic,
	status
)
VALUES ($1, $2, $3, $4, $5::uuid, $5::uuid, $5::uuid, $6, $7, $8::date, $9::time, $10::time, $11, 'planned')
ON CONFLICT (schedule_id, lesson_date) WHERE schedule_id IS NOT NULL
DO NOTHING
RETURNING
id,
schedule_id,
organization_id,
branch_id,
group_id,
COALESCE(teacher_id::text, ''),
subject_id,
lesson_date::text,
start_time::text,
end_time::text,
COALESCE(topic, ''),
status
`,
			lessonID,
			organizationID,
			schedule.BranchID,
			schedule.GroupID,
			teacherID,
			schedule.SubjectID,
			schedule.ID,
			lessonDate,
			schedule.StartTime,
			schedule.EndTime,
			topic,
		).Scan(
			&created.ID,
			&created.ScheduleID,
			&created.OrganizationID,
			&created.BranchID,
			&created.GroupID,
			&created.TeacherID,
			&created.SubjectID,
			&created.LessonDate,
			&created.StartTime,
			&created.EndTime,
			&created.Topic,
			&created.Status,
		)

		if err != nil {
			if err == pgx.ErrNoRows {
				skippedDates = append(skippedDates, lessonDate)
				continue
			}

			return nil, nil, err
		}

		createdLessons = append(createdLessons, created)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, nil, err
	}

	return createdLessons, skippedDates, nil
}

func weekdayNumber(date time.Time) int {
	weekday := int(date.Weekday())
	if weekday == 0 {
		return 7
	}

	return weekday
}
