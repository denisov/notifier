package kengu

// DataSource это источник данных
type DataSource interface {
	GetData() (string, error)
}
