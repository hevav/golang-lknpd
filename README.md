## (Пока что) не протестировано

# API для lknpd.nalog.ru (Мой налог)
Эта библиотка позволяет создавать/отменять чеки самозанятого <br>
Создано для [elling-npd](https://github.com/Elytrium/elling-npd)

#### ⚠ Автор не несет ответственности за ваше использование этой библиотеки

Для установки в ваш проект, введите команду:
```shell
go get -u github.com/hevav/golang-lknpd
```
## Использование
### Создание клиента
```go
NalogClient := lknpd.CreateClient(randomDeviceId) // randomDeviceId - 20-22х значная строка
err := NalogClient.Auth(login, password)
```
### Добавление дохода
```go
income := lknpd.DefaultIncome()

income.AddService(Service{
	Name: "Предоставление услуг #123456",
	Amount: amount
    Quantity: 1,
})

income.SetClientType(lknpd.LegalEntity) // по дефолту - Individual
income.SetClientName("ИП Иванов Иван Иванович") // только для юр. лица (LEGAL_ENTITY)
income.SetClientINN("123456789012") // только для юр. лица зарегистрированного в РФ
income.SetOperationTime(time.Date(2021, time.February, 13, 23, 59, 59, 0, time.Local)) // по дефолту - time.Now()

receipt, err := NalogClient.AddIncome(income)

if err != nil {
	panic(err)
}

fmt.Println(receipt.UUID) // ID чека
fmt.Println(receipt.InfoURL) // Ссылка на информацию о чеке
fmt.Println(receipt.PrintURL) // Ссылка на изображение чека
```