package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/shigaichi/top-sites-ranking-api/internal/domain/model"

	"github.com/shigaichi/top-sites-ranking-api/internal/adapter/http/handler/dto"

	"github.com/shigaichi/top-sites-ranking-api/internal/usecase"
	log "github.com/sirupsen/logrus"
)

type GetRanking interface {
	GetDailyRanking(w http.ResponseWriter, r *http.Request)
	GetMonthlyRanking(w http.ResponseWriter, r *http.Request)
}

type GetRankingImpl struct {
	u usecase.RankHistoryUseCase
}

func NewGetRankingImpl(u usecase.RankHistoryUseCase) *GetRankingImpl {
	return &GetRankingImpl{u: u}
}

func (g GetRankingImpl) GetDailyRanking(w http.ResponseWriter, r *http.Request) {
	domain := r.URL.Query().Get("domain")
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	if domain == "" || startDateStr == "" || endDateStr == "" {
		http.Error(w, "Bad Request: Missing or invalid query parameters", http.StatusBadRequest)
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		http.Error(w, "Bad Request: Invalid start_date format", http.StatusBadRequest)
		return
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		http.Error(w, "Bad Request: Invalid end_date format", http.StatusBadRequest)
		return
	}

	if startDate.After(endDate) {
		http.Error(w, "start_date should be before or equal to end_date", http.StatusBadRequest)
		return
	}

	ranks, err := g.u.GetDailyRanking(r.Context(), domain, startDate, endDate)
	if err != nil {
		log.WithContext(r.Context()).WithFields(log.Fields{"error": err, "domain": domain, "stat_date": startDateStr, "end_date": endDateStr}).Error("GetDailyRanking usecase returned error while processing daily ranking")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if len(ranks) == 0 {
		http.Error(w, "Not Found: No rankings available for the given period", http.StatusNotFound)
		return
	}

	responseRanks := make([]dto.ResponseRank, len(ranks))
	for i, rank := range ranks {
		responseRanks[i] = dto.ResponseRank{
			Rank: rank.Rank,
			Date: rank.Date.UTC().Format("2006-01-02"),
		}
	}

	resp := struct {
		Ranks  []dto.ResponseRank `json:"ranks"`
		Domain string             `json:"domain"`
	}{
		Ranks:  responseRanks,
		Domain: domain,
	}

	if isIncludingEveryDayRecord(startDate, endDate, ranks) {
		w.Header().Set("Cache-Control", "max-age=86400")
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.WithContext(r.Context()).WithFields(log.Fields{"error": err, "domain": domain, "stat_date": startDateStr, "end_date": endDateStr}).Error("cannot marshall to resp json while processing daily ranking")
		http.Error(w, "Failed to encode the resp", http.StatusInternalServerError)
	}
}

func isIncludingEveryDayRecord(start, end time.Time, ranks []model.DailyRank) bool {
	d := end.Add(time.Hour * 24).Sub(start)
	return int(d.Hours()/24) == len(ranks)
}

func (g GetRankingImpl) GetMonthlyRanking(w http.ResponseWriter, r *http.Request) {
	domain := r.URL.Query().Get("domain")
	startMonthStr := r.URL.Query().Get("start_month")
	endMonthStr := r.URL.Query().Get("end_month")

	if domain == "" || startMonthStr == "" || endMonthStr == "" {
		http.Error(w, "Bad Request: Missing or invalid query parameters", http.StatusBadRequest)
		return
	}

	// "2023-12" -> 2023-12-01 00:00:00 +0000
	startMonth, err := time.Parse("2006-01", startMonthStr)
	if err != nil {
		http.Error(w, "Invalid start_month format", http.StatusBadRequest)
		return
	}

	endMonth, err := time.Parse("2006-01", endMonthStr)
	if err != nil {
		http.Error(w, "Invalid end_month format", http.StatusBadRequest)
		return
	}

	if startMonth.After(endMonth) {
		http.Error(w, "start_month should be before or equal to end_month", http.StatusBadRequest)
		return
	}

	ranks, err := g.u.GetMonthlyRanking(r.Context(), domain, getLastDayOfMonth(startMonth), getLastDayOfMonth(endMonth))
	if err != nil {
		log.WithContext(r.Context()).WithFields(log.Fields{"error": err, "domain": domain, "stat_date": startMonthStr, "end_date": endMonthStr}).Error("GetMonthlyRanking usecase returned error while processing monthly ranking")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if len(ranks) == 0 {
		http.Error(w, "No ranks found for the given domain and date range", http.StatusNotFound)
		return
	}

	responseRanks := make([]dto.ResponseRank, len(ranks))
	for i, rank := range ranks {
		responseRanks[i] = dto.ResponseRank{
			Rank: rank.Rank,
			Date: rank.Date.UTC().Format("2006-01-02"),
		}
	}

	resp := struct {
		Ranks  []dto.ResponseRank `json:"ranks"`
		Domain string             `json:"domain"`
	}{
		Ranks:  responseRanks,
		Domain: domain,
	}

	if isIncludingEveryMonthRecord(getLastDayOfMonth(startMonth), getLastDayOfMonth(endMonth), ranks) {
		w.Header().Set("Cache-Control", "max-age=86400")
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.WithContext(r.Context()).WithFields(log.Fields{"error": err, "domain": domain, "stat_date": startMonthStr, "end_date": endMonthStr}).Error("cannot marshall to response json while processing monthly ranking")
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func getLastDayOfMonth(t time.Time) time.Time {
	nextMonth := t.AddDate(0, 1, 0)
	lastDay := nextMonth.AddDate(0, 0, -1)
	return lastDay
}

func isIncludingEveryMonthRecord(start, end time.Time, ranks []model.DailyRank) bool {
	d := countMonthEnds(start, end)
	return d == len(ranks)
}

// countMonthEnds counts how many month ends are there between two dates, inclusive.
func countMonthEnds(start, end time.Time) int {
	count := 0
	current := start

	for current.Before(end) || current.Equal(end) {
		// Check if tomorrow's month is different, which means today is the end of the month.
		if current.AddDate(0, 0, 1).Month() != current.Month() {
			count++
		}
		current = current.AddDate(0, 0, 1)
	}

	return count
}
