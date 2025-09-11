package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

type GenGenerator struct {
	gen *gen.Generator
}

func MakeGenGenerator(OutPath string, dsn string) *GenGenerator {
	g := gen.NewGenerator(gen.Config{
		OutPath:       OutPath,
		Mode:          gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable: true,
	})

	db, _ := gorm.Open(mysql.Open(dsn))

	g.UseDB(db)

	return &GenGenerator{gen: g}
}

func (tg *GenGenerator) GenerateModelWithExecute(models []string) {
	var genModels []any
	for _, model := range models {
		genModels = append(genModels, tg.gen.GenerateModel(model))
	}

	tg.gen.ApplyBasic(genModels)
	tg.gen.Execute()
}

func main() {
	fmt.Println("Generate Model from SQL running...")
}
