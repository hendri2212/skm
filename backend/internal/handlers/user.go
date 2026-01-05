package handlers

import (
	"context"
	"fmt"
	"net/http"
	"skm/internal/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

// ---------- RESPONSE DTOs ----------
type ChoiceDisplay struct {
	ID          uint   `json:"id"`
	DisplayText string `json:"display_text"` // akan di-bold bila dipilih
	RawText     string `json:"raw_text"`
	Selected    bool   `json:"selected"`
}

type QuestionDetail struct {
	QuestionID   uint            `json:"question_id"`
	QuestionText string          `json:"question_text"`
	Choices      []ChoiceDisplay `json:"choices"`
	SelectedID   *uint           `json:"selected_choice_id,omitempty"`
}

type UserAnswerDetailResponse struct {
	User struct {
		ID                uint   `json:"id"`
		FullName          string `json:"full_name"`
		PlaceOfBirth      string `json:"place_of_birth,omitempty"`
		DateOfBirth       string `json:"date_of_birth,omitempty"`
		IsMale            *bool  `json:"is_male,omitempty"`
		LastEducationID   *uint  `json:"last_education_id,omitempty"`
		LastEducationName string `json:"last_education_name,omitempty"`
		OccupationID      *uint  `json:"main_occupation_id,omitempty"`
		OccupationName    string `json:"main_occupation_name,omitempty"`
	} `json:"user"`
	Answers []QuestionDetail `json:"answers"`
}

type UserOption struct {
	Value uint   `json:"value"`
	Label string `json:"label"`
}

func UsersHandler(db *gorm.DB) *UserHandler {
	// NOTE: Di production, pertimbangkan memindahkan AutoMigrate ke proses migrasi terpisah.
	db.AutoMigrate(&models.User{}, &models.Education{}, &models.Occupation{})
	return &UserHandler{db: db}
}

// ----------------- Helpers (aman & tanpa case nil) -----------------

// fmtPct: format float ke 2 desimal sebagai string
func fmtPct(v float64) string { return fmt.Sprintf("%.2f", v) }

// toDateString: terima either time.Time atau *time.Time dari field model yang bisa bervariasi
func toDateString(v any) (string, bool) {
	if t, ok := v.(*time.Time); ok && t != nil && !t.IsZero() {
		return t.Format("2006-01-02"), true
	}
	if t, ok := v.(time.Time); ok && !t.IsZero() {
		return t.Format("2006-01-02"), true
	}
	return "", false
}

// toBoolPtr: support bool atau *bool → *bool
func toBoolPtr(v any) *bool {
	if p, ok := v.(*bool); ok {
		return p
	}
	if b, ok := v.(bool); ok {
		bb := b
		return &bb
	}
	return nil
}

// toUintPtr: support uint atau *uint → *uint
func toUintPtr(v any) *uint {
	if p, ok := v.(*uint); ok {
		return p
	}
	if u, ok := v.(uint); ok {
		uu := u
		return &uu
	}
	return nil
}

// dbx: shortcut WithContext untuk semua query (cancellable)
func (h *UserHandler) dbx(ctx context.Context) *gorm.DB {
	return h.db.WithContext(ctx)
}

// ----------------- Handlers -----------------

// CountAge returns total respondents and breakdown from youngest to oldest
func (h *UserHandler) CountAge(c *gin.Context) {
	var res models.AgeCountResult
	err := h.dbx(c.Request.Context()).Raw(`
        SELECT
          SUM(CASE WHEN TIMESTAMPDIFF(YEAR, date_of_birth, CURDATE()) <= 19 THEN 1 ELSE 0 END) AS age_19_under,
          SUM(CASE WHEN TIMESTAMPDIFF(YEAR, date_of_birth, CURDATE()) BETWEEN 20 AND 30 THEN 1 ELSE 0 END) AS age_20_30,
          SUM(CASE WHEN TIMESTAMPDIFF(YEAR, date_of_birth, CURDATE()) BETWEEN 31 AND 49 THEN 1 ELSE 0 END) AS age_31_49,
          SUM(CASE WHEN TIMESTAMPDIFF(YEAR, date_of_birth, CURDATE()) >= 50 THEN 1 ELSE 0 END) AS age_50_above,
          COUNT(*) AS total
        FROM users
    `).Scan(&res).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var pct19, pct20_30, pct31_49, pct50 float64
	if res.Total > 0 {
		t := float64(res.Total)
		pct19 = float64(res.Age19Under) / t * 100
		pct20_30 = float64(res.Age20_30) / t * 100
		pct31_49 = float64(res.Age31_49) / t * 100
		pct50 = float64(res.Age50Above) / t * 100
	}

	type ageGroup struct {
		AgeRange string `json:"age_range"`
		Count    int64  `json:"count"`
		Percent  string `json:"percent"`
	}
	groups := make([]ageGroup, 0, 4)
	groups = append(groups,
		ageGroup{AgeRange: "<=19", Count: res.Age19Under, Percent: fmtPct(pct19)},
		ageGroup{AgeRange: "20-30", Count: res.Age20_30, Percent: fmtPct(pct20_30)},
		ageGroup{AgeRange: "31-49", Count: res.Age31_49, Percent: fmtPct(pct31_49)},
		ageGroup{AgeRange: ">=50", Count: res.Age50Above, Percent: fmtPct(pct50)},
	)

	c.JSON(http.StatusOK, gin.H{
		"total_respondents": res.Total,
		"age_groups":        groups,
	})
}

func (h *UserHandler) CountGender(c *gin.Context) {
	var res models.GenderCountResult

	err := h.dbx(c.Request.Context()).Raw(`
		SELECT
			SUM(CASE WHEN is_male = 1 THEN 1 ELSE 0 END) AS male,
			SUM(CASE WHEN is_male = 0 THEN 1 ELSE 0 END) AS female,
			COUNT(*) AS total
		FROM users
	`).Scan(&res).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var pctMale, pctFemale float64
	if res.Total > 0 {
		t := float64(res.Total)
		pctMale = float64(res.Male) / t * 100
		pctFemale = float64(res.Female) / t * 100
	}

	type genderGroup struct {
		Gender  string `json:"gender"`
		Count   int64  `json:"count"`
		Percent string `json:"percent"`
	}

	result := []genderGroup{
		{Gender: "Laki-laki", Count: res.Male, Percent: fmtPct(pctMale)},
		{Gender: "Perempuan", Count: res.Female, Percent: fmtPct(pctFemale)},
	}

	c.JSON(http.StatusOK, gin.H{
		"total_respondents": res.Total,
		"gender_groups":     result,
	})
}

func (h *UserHandler) CountEducation(c *gin.Context) {
	type row struct {
		EducationID *uint  `json:"id,omitempty"`
		Label       string `json:"label"`
		Count       int64  `json:"count"`
	}

	var rows []row
	err := h.dbx(c.Request.Context()).Raw(`
		SELECT
			e.id AS education_id,
			COALESCE(e.name, 'Tidak diisi') AS label,
			COUNT(u.id) AS count
		FROM users u
		LEFT JOIN educations e ON e.id = u.last_education_id
		GROUP BY e.id, e.name
		ORDER BY CASE WHEN e.id IS NULL THEN 1 ELSE 0 END, e.id ASC
	`).Scan(&rows).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var total int64
	for i := range rows {
		total += rows[i].Count
	}

	type eduGroup struct {
		ID      *uint  `json:"id,omitempty"`
		Label   string `json:"label"`
		Count   int64  `json:"count"`
		Percent string `json:"percent"`
	}

	out := make([]eduGroup, 0, len(rows))
	var denom float64 = 1
	if total > 0 {
		denom = float64(total)
	}
	for _, r := range rows {
		pct := float64(r.Count) / denom * 100
		out = append(out, eduGroup{
			ID:      r.EducationID,
			Label:   r.Label,
			Count:   r.Count,
			Percent: fmtPct(pct),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"total_respondents": total,
		"education_groups":  out,
	})
}

func (h *UserHandler) CountOccupation(c *gin.Context) {
	type row struct {
		OccupationID *uint  `json:"id,omitempty"`
		Label        string `json:"label"`
		Count        int64  `json:"count"`
	}

	var rows []row
	err := h.dbx(c.Request.Context()).Raw(`
		SELECT
			o.id AS occupation_id,
			COALESCE(o.name, 'Tidak diisi') AS label,
			COUNT(u.id) AS count
		FROM users u
		LEFT JOIN occupations o ON o.id = u.main_occupation_id
		GROUP BY o.id, o.name
		ORDER BY CASE WHEN o.id IS NULL THEN 1 ELSE 0 END, o.id ASC
	`).Scan(&rows).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var total int64
	for i := range rows {
		total += rows[i].Count
	}

	type occGroup struct {
		ID      *uint  `json:"id,omitempty"`
		Label   string `json:"label"`
		Count   int64  `json:"count"`
		Percent string `json:"percent"`
	}

	out := make([]occGroup, 0, len(rows))
	var denom float64 = 1
	if total > 0 {
		denom = float64(total)
	}
	for _, r := range rows {
		pct := float64(r.Count) / denom * 100
		out = append(out, occGroup{
			ID:      r.OccupationID,
			Label:   r.Label,
			Count:   r.Count,
			Percent: fmtPct(pct),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"total_respondents": total,
		"occupation_groups": out,
	})
}

// Report is a placeholder endpoint for aggregated reporting.
// Adjust the response structure to your needs.
func (h *UserHandler) Report(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "status":  "ok",
        "message": "report endpoint",
    })
}

func (h *UserHandler) GetUserAnswers(c *gin.Context) {
	type row struct {
		UserName string `json:"user_name" gorm:"column:user_name"`
		Soal1    *int   `json:"soal_1"    gorm:"column:soal_1"`
		Soal2    *int   `json:"soal_2"    gorm:"column:soal_2"`
		Soal3    *int   `json:"soal_3"    gorm:"column:soal_3"`
		Soal4    *int   `json:"soal_4"    gorm:"column:soal_4"`
		Soal5    *int   `json:"soal_5"    gorm:"column:soal_5"`
		Soal6    *int   `json:"soal_6"    gorm:"column:soal_6"`
		Soal7    *int   `json:"soal_7"    gorm:"column:soal_7"`
		Soal8    *int   `json:"soal_8"    gorm:"column:soal_8"`
		Soal9    *int   `json:"soal_9"    gorm:"column:soal_9"`
	}

	var rows []row
	const sql = `
		SELECT 
			u.full_name AS user_name,
			MAX(CASE WHEN a.question_id = 1 THEN c.points END) AS soal_1,
			MAX(CASE WHEN a.question_id = 2 THEN c.points END) AS soal_2,
			MAX(CASE WHEN a.question_id = 3 THEN c.points END) AS soal_3,
			MAX(CASE WHEN a.question_id = 4 THEN c.points END) AS soal_4,
			MAX(CASE WHEN a.question_id = 5 THEN c.points END) AS soal_5,
			MAX(CASE WHEN a.question_id = 6 THEN c.points END) AS soal_6,
			MAX(CASE WHEN a.question_id = 7 THEN c.points END) AS soal_7,
			MAX(CASE WHEN a.question_id = 8 THEN c.points END) AS soal_8,
			MAX(CASE WHEN a.question_id = 9 THEN c.points END) AS soal_9
		FROM users u
		LEFT JOIN answers a ON u.id = a.user_id
		LEFT JOIN choices c ON c.id = a.choice_id
		GROUP BY u.id, u.full_name
	`
	if err := h.dbx(c.Request.Context()).Raw(sql).Scan(&rows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	totalUsers := len(rows)

	type agg struct {
		Soal1 int64 `gorm:"column:soal_1"`
		Soal2 int64 `gorm:"column:soal_2"`
		Soal3 int64 `gorm:"column:soal_3"`
		Soal4 int64 `gorm:"column:soal_4"`
		Soal5 int64 `gorm:"column:soal_5"`
		Soal6 int64 `gorm:"column:soal_6"`
		Soal7 int64 `gorm:"column:soal_7"`
		Soal8 int64 `gorm:"column:soal_8"`
		Soal9 int64 `gorm:"column:soal_9"`
	}
	var a agg

	const aggSQL = `
		SELECT
			SUM(CASE WHEN a.question_id = 1 THEN COALESCE(c.points,0) ELSE 0 END) AS soal_1,
			SUM(CASE WHEN a.question_id = 2 THEN COALESCE(c.points,0) ELSE 0 END) AS soal_2,
			SUM(CASE WHEN a.question_id = 3 THEN COALESCE(c.points,0) ELSE 0 END) AS soal_3,
			SUM(CASE WHEN a.question_id = 4 THEN COALESCE(c.points,0) ELSE 0 END) AS soal_4,
			SUM(CASE WHEN a.question_id = 5 THEN COALESCE(c.points,0) ELSE 0 END) AS soal_5,
			SUM(CASE WHEN a.question_id = 6 THEN COALESCE(c.points,0) ELSE 0 END) AS soal_6,
			SUM(CASE WHEN a.question_id = 7 THEN COALESCE(c.points,0) ELSE 0 END) AS soal_7,
			SUM(CASE WHEN a.question_id = 8 THEN COALESCE(c.points,0) ELSE 0 END) AS soal_8,
			SUM(CASE WHEN a.question_id = 9 THEN COALESCE(c.points,0) ELSE 0 END) AS soal_9
		FROM users u
		LEFT JOIN answers a ON u.id = a.user_id
		LEFT JOIN choices c ON c.id = a.choice_id
	`
	if err := h.dbx(c.Request.Context()).Raw(aggSQL).Scan(&a).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	type stat struct {
		SoalID        int    `json:"soal_id"`
		Total         int64  `json:"total"`
		Average       string `json:"average"`
		NRRTertimbang string `json:"nrr_tertimbang"`
		IKM           string `json:"ikm"`
	}
	stats := make([]stat, 0, 9)

	denom := float64(totalUsers)
	if totalUsers == 0 {
		denom = 1 // hindari div-by-zero; hasil avg akan 0.00
	}

	appendStat := func(id int, total int64) {
		avg := float64(total) / denom
		nrrT := avg / 9.0
		ikm := nrrT * 25.0
		stats = append(stats, stat{
			SoalID:        id,
			Total:         total,
			Average:       fmtPct(avg),
			NRRTertimbang: fmtPct(nrrT),
			IKM:           fmtPct(ikm),
		})
	}

	appendStat(1, a.Soal1)
	appendStat(2, a.Soal2)
	appendStat(3, a.Soal3)
	appendStat(4, a.Soal4)
	appendStat(5, a.Soal5)
	appendStat(6, a.Soal6)
	appendStat(7, a.Soal7)
	appendStat(8, a.Soal8)
	appendStat(9, a.Soal9)

	c.JSON(http.StatusOK, gin.H{
		"total_users": totalUsers,
		"data":        rows,
		"stats":       stats,
	})
}

func (h *UserHandler) GetUserAnswerByID(c *gin.Context) {
	idStr := c.Param("id")
	uid, err := strconv.Atoi(idStr)
	if err != nil || uid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id tidak valid"})
		return
	}

	// 1) Ambil data user
	var user models.User
	if err := h.dbx(c.Request.Context()).
		Preload("LastEducation").
		Preload("MainOccupation").
		First(&user, uid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 2) Ambil semua jawaban user + relasi Question.Choices dan Choice terpilih
	var answers []models.Answer
	if err := h.dbx(c.Request.Context()).
		Preload("Question").
		Preload("Question.Choices").
		Preload("Choice").
		Where("user_id = ?", uid).
		Find(&answers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3) Susun response user (tanpa type switch case nil)
	resp := UserAnswerDetailResponse{}
	resp.User.ID = user.ID
	resp.User.FullName = user.FullName
	resp.User.PlaceOfBirth = user.PlaceOfBirth

	if s, ok := toDateString(any(user.DateOfBirth)); ok {
		resp.User.DateOfBirth = s
	}
	resp.User.IsMale = toBoolPtr(any(user.IsMale))
	resp.User.LastEducationID = toUintPtr(any(user.LastEducationID))
	resp.User.OccupationID = toUintPtr(any(user.MainOccupationID))

	if user.LastEducation != nil {
		resp.User.LastEducationName = user.LastEducation.Name
	}
	if user.MainOccupation != nil {
		resp.User.OccupationName = user.MainOccupation.Name
	}

	// 4) Detail per soal
	resp.Answers = make([]QuestionDetail, 0, len(answers))
	for _, a := range answers {
		q := a.Question

		var selectedID *uint
		if a.ChoiceID != 0 {
			selectedID = &a.ChoiceID
		}

		qd := QuestionDetail{
			QuestionID:   q.ID,
			QuestionText: q.QuestionText, // sesuaikan jika nama field berbeda
			SelectedID:   selectedID,
			Choices:      make([]ChoiceDisplay, 0, len(q.Choices)),
		}

		for _, ch := range q.Choices {
			text := ch.ChoiceText // ganti jika field berbeda
			cd := ChoiceDisplay{
				ID:       ch.ID,
				RawText:  text,
				Selected: selectedID != nil && ch.ID == *selectedID,
			}
			if cd.Selected {
				cd.DisplayText = "**" + text + "**"
			} else {
				cd.DisplayText = text
			}
			qd.Choices = append(qd.Choices, cd)
		}

		resp.Answers = append(resp.Answers, qd)
	}

	c.JSON(http.StatusOK, resp)
}

// GetUserAnswerAll mengembalikan daftar UserAnswerDetailResponse untuk semua user.
// Struktur elemen di dalam "data" identik dengan GetUserAnswerByID.
func (h *UserHandler) GetUserAnswerAll(c *gin.Context) {
	ctx := c.Request.Context()

	// 1) Ambil semua user + relasi ringkas yg dibutuhkan untuk header user
	var users []models.User
	if err := h.dbx(ctx).
		Preload("LastEducation").
		Preload("MainOccupation").
		Order("full_name ASC").
		Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(users) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"count": 0,
			"data":  []UserAnswerDetailResponse{},
		})
		return
	}

	// 2) Ambil seluruh jawaban untuk semua user sekaligus + preload relasi yang sama seperti GetUserAnswerByID
	ids := make([]uint, 0, len(users))
	for _, u := range users {
		ids = append(ids, u.ID)
	}

	var answers []models.Answer
	if err := h.dbx(ctx).
		Preload("Question").
		Preload("Question.Choices").
		Preload("Choice").
		Where("user_id IN ?", ids).
		Find(&answers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3) Kelompokkan jawaban per user_id agar assembling O(n)
	ansByUser := make(map[uint][]models.Answer, len(users))
	for _, a := range answers {
		ansByUser[a.UserID] = append(ansByUser[a.UserID], a)
	}

	// 4) Susun response per user (identik dengan GetUserAnswerByID)
	out := make([]UserAnswerDetailResponse, 0, len(users))
	for _, user := range users {
		resp := UserAnswerDetailResponse{}
		resp.User.ID = user.ID
		resp.User.FullName = user.FullName
		resp.User.PlaceOfBirth = user.PlaceOfBirth
		if s, ok := toDateString(any(user.DateOfBirth)); ok {
			resp.User.DateOfBirth = s
		}
		resp.User.IsMale = toBoolPtr(any(user.IsMale))
		resp.User.LastEducationID = toUintPtr(any(user.LastEducationID))
		resp.User.OccupationID = toUintPtr(any(user.MainOccupationID))
		if user.LastEducation != nil {
			resp.User.LastEducationName = user.LastEducation.Name
		}
		if user.MainOccupation != nil {
			resp.User.OccupationName = user.MainOccupation.Name
		}

		// Detail per soal
		ua := ansByUser[user.ID]
		resp.Answers = make([]QuestionDetail, 0, len(ua))
		for _, a := range ua {
			q := a.Question

			var selectedID *uint
			if a.ChoiceID != 0 {
				selectedID = &a.ChoiceID
			}

			qd := QuestionDetail{
				QuestionID:   q.ID,
				QuestionText: q.QuestionText,
				SelectedID:   selectedID,
				Choices:      make([]ChoiceDisplay, 0, len(q.Choices)),
			}

			for _, ch := range q.Choices {
				text := ch.ChoiceText
				cd := ChoiceDisplay{
					ID:       ch.ID,
					RawText:  text,
					Selected: selectedID != nil && ch.ID == *selectedID,
				}
				if cd.Selected {
					cd.DisplayText = "**" + text + "**"
				} else {
					cd.DisplayText = text
				}
				qd.Choices = append(qd.Choices, cd)
			}
			resp.Answers = append(resp.Answers, qd)
		}

		out = append(out, resp)
	}

	// 5) Response
	c.JSON(http.StatusOK, gin.H{
		"count": len(out),
		"data":  out,
	})
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	ctx := c.Request.Context()

	// Ambil hanya kolom yang diperlukan
	var rows []struct {
		ID       uint   `gorm:"column:id"`
		FullName string `gorm:"column:full_name"`
	}

	if err := h.dbx(ctx).
		Model(&models.User{}).
		Select("id, full_name").
		Order("full_name ASC").
		Scan(&rows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Map ke format options
	opts := make([]UserOption, len(rows))
	for i, r := range rows {
		opts[i] = UserOption{
			Value: r.ID,
			Label: r.FullName,
		}
	}

	// Response ringkas & siap pakai di select
	c.JSON(http.StatusOK, gin.H{
		"options": opts,
		"count":   len(opts),
	})
}
