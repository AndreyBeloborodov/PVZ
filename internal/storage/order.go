package storage

import (
	"gitlab.ozon.dev/go/classroom-13/students/Homework/Homework-1/internal/models"
	"sync"
	"time"
)

type OrderRecord struct {
	OrderId     int       `db:"id"`
	UserId      int       `db:"user_id"`
	TimeEnd     time.Time `db:"time_end"`
	IsGiven     bool      `db:"is_given"`
	TimeGiven   time.Time `db:"time_given"`
	IsReturned  bool      `db:"is_returned"`
	Price       int       `db:"price"`
	Weight      int       `db:"weight"`
	PackageName string    `db:"name_package"`
	ResultPrice int       `db:"result_price"`
}

func (t OrderRecord) toDomain() models.Order {
	return models.Order{
		OrderId:     models.OrderId(t.OrderId),
		UserId:      models.UserId(t.UserId),
		TimeEnd:     models.TimeEnd(t.TimeEnd),
		IsGiven:     models.IsGiven(t.IsGiven),
		TimeGiven:   models.TimeGiven(t.TimeGiven),
		IsReturned:  models.IsReturned(t.IsReturned),
		FirstPrice:  models.FirstPrice(t.Price),
		Weight:      models.OrderWeight(t.Weight),
		PackageName: models.PackageName(t.PackageName),
		ResultPrice: models.ResultPrice(t.ResultPrice),
	}
}

func transform(order models.Order) OrderRecord {
	return OrderRecord{
		OrderId:     int(order.OrderId),
		UserId:      int(order.UserId),
		TimeEnd:     time.Time(order.TimeEnd),
		IsGiven:     bool(order.IsGiven),
		TimeGiven:   time.Time(order.TimeGiven),
		IsReturned:  bool(order.IsReturned),
		Price:       int(order.FirstPrice),
		Weight:      int(order.Weight),
		PackageName: string(order.PackageName),
		ResultPrice: int(order.ResultPrice),
	}
}

func toDomains(records []OrderRecord) []models.Order {
	const (
		numWorkers = 5
	)

	var (
		numJobs = len(records)
		jobs    = make(chan OrderRecord, numJobs) // буферизированный канал для задач (размер numJobs)
		result  = make(chan models.Order, numJobs)
		wg      = sync.WaitGroup{}
	)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(&wg, jobs, result)
	}

	go func() {
		for _, record := range records {
			jobs <- record
		}

		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(result)
	}()

	resultOrders := make([]models.Order, numJobs)
	it := 0

	for r := range result {
		resultOrders[it] = r
		it++
	}

	return resultOrders
}

func worker(wg *sync.WaitGroup, jobs <-chan OrderRecord, result chan<- models.Order) {
	defer wg.Done()

	for j := range jobs {
		result <- j.toDomain()
	}
}
