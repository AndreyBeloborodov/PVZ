package console

import (
	"context"
	"fmt"
	"gitlab.ozon.dev/go/classroom-13/students/Homework/Homework-1/internal/models"
	"gitlab.ozon.dev/go/classroom-13/students/Homework/Homework-1/internal/services"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Console struct {
	Service        *services.OrderService
	HistoryCommand []string
	Data           map[string]any
	mu             sync.Mutex
}

const offsetReturns = 2

func NewConsole(service *services.OrderService) *Console {
	return &Console{
		Service:        service,
		HistoryCommand: make([]string, 0),
		Data:           map[string]any{},
	}
}

func getCommandList(commandLine string) []string {
	commandLine = strings.Replace(commandLine, "\n", "", -1)
	commandLine = strings.Replace(commandLine, "\r", "", -1)
	for strings.Contains(commandLine, "  ") {
		commandLine = strings.Replace(commandLine, "  ", " ", -1)
	}
	return strings.Split(commandLine, " ")
}

func (c *Console) Solve(ctx context.Context, commandLine string) {
	command := getCommandList(commandLine)
	c.HistoryCommand = append(c.HistoryCommand, command[0])

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Проверьте параметры запроса")
			log.Fatal(r)
		}
	}()

	switch command[0] {
	case "help":
		fmt.Println("список доступных комаед:")
		fmt.Println("add id1 id2 time weight price pack - Принять заказ от курьера, time - кол-во дней (add 123 321 5 10 120 (film/box/package))")
		fmt.Println("delete id - Вернуть заказ курьеру (delete 123)")
		fmt.Println("give id1 1d2 1d3... - Выдать заказы клиенту (give 123 321 12 21)")
		fmt.Println("orders id n ok - Получить список заказов (orders 123 5 true, orders 123 -1 true), (-1) - значит вернуть все")
		fmt.Println("return id1 1d2 - Принять возврат от клиента (return 123 321)")
		fmt.Println("returns - Получить список возвратов, чтобы посмотреть больше нужно снова ввести команду (returns)")
		return
	case "add":
		orderId, _ := strconv.Atoi(command[1])
		userId, _ := strconv.Atoi(command[2])
		timeSafe, _ := strconv.Atoi(command[3])
		weight, _ := strconv.Atoi(command[4])
		price, _ := strconv.Atoi(command[5])
		packageName := command[6]

		params := services.ParamsAddOrder{
			OrderId:     orderId,
			UserId:      userId,
			TimeEnd:     time.Now().Add(time.Hour * 24 * time.Duration(timeSafe)),
			Weight:      weight,
			Price:       price,
			PackageName: packageName,
		}

		if err := c.Service.AddOrder(ctx, params); err != nil {
			fmt.Println(command[0] + " error : " + err.Error())
		} else {
			fmt.Println("заказ принят")
		}
		return
	case "delete":
		orderId, _ := strconv.Atoi(command[1])
		if err := c.Service.ReturnOrderToCourier(ctx, orderId); err != nil {
			fmt.Println(command[0] + " error : " + err.Error())
		} else {
			fmt.Println("заказ возвращён курьеру")
		}
		return
	case "give":
		ordersId := make([]int, 0)
		for i := 1; i < len(command); i++ {
			orderId, _ := strconv.Atoi(command[i])
			ordersId = append(ordersId, orderId)
		}
		if err := c.Service.GiveOrderToUser(ctx, ordersId); err != nil {
			fmt.Println(command[0] + " error : " + err.Error())
		} else {
			fmt.Println("заказ выдан клиенту")
		}
		return
	case "orders":
		userId, _ := strconv.Atoi(command[1])

		limit, _ := strconv.Atoi(command[2])
		onlyExist := command[3] == "true"
		options := services.ParamsGetOrders{
			Limit:     limit,
			OnlyExist: onlyExist,
		}

		orders, err := c.Service.GetOrders(ctx, userId, options)
		if err != nil {
			fmt.Println(command[0] + " error : " + err.Error())
		} else {
			if len(orders) == 0 {
				fmt.Println("заказов по вашему запросу не нашлось")
				return
			}
			fmt.Println("заказы:")
			for _, order := range orders {
				fmt.Println(order)
			}
		}
		return
	case "return":
		orderId, _ := strconv.Atoi(command[1])
		userId, _ := strconv.Atoi(command[2])
		if err := c.Service.AcceptReturnOrder(ctx, orderId, userId); err != nil {
			fmt.Println(command[0] + " error : " + err.Error())
		} else {
			fmt.Println("возврат успешно принят")
		}
		return
	case "returns":
		fmt.Println("возвраты:")

		c.mu.Lock()
		c.Data["func-next-returns"] = c.Service.GetReturnsPagination(ctx)
		c.mu.Unlock()

		c.Solve(ctx, "more")
		return
	case "more":

		c.mu.Lock()
		f, ok := c.Data["func-next-returns"].(func(count int) []models.Order)
		c.mu.Unlock()

		if !ok {
			fmt.Println("Сначала введите команду returns")
			return
		}
		for _, ret := range f(offsetReturns) {
			fmt.Println(ret)
		}
		return
	case "exit":
		os.Exit(1)
		return
	default:
		fmt.Println("нет команды: " + command[0])
	}
}
