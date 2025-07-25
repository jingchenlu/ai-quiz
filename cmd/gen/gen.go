package main

import (
	"aiquiz/config"
	"aiquiz/dao/model"
	"aiquiz/migrations"
	"gorm.io/gen"
	"log"
	"strings"
)

func tableToStructName(table string) string {
	name := table
	// 下划线转驼峰
	parts := strings.Split(name, "_")
	for i, p := range parts {
		parts[i] = strings.Title(p)
	}
	return strings.Join(parts, "")
}

func main() {
	appConfig := config.GetConfig(false)

	// 初始化数据库连接
	db := migrations.InitDB(appConfig.DBPath)
	log.Println("数据库初始化完成")

	g := gen.NewGenerator(gen.Config{
		OutPath: "./dao/query",

		Mode: gen.WithDefaultQuery | gen.WithQueryInterface,
		// 表字段可为null时，对应字段使用指针类型
		FieldNullable: true,
	})

	g.UseDB(db)

	// 从数据库表生成所有模型
	// 获取所有表名，排除sqlite系统表

	var tables []string
	db.Raw("SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'").Scan(&tables)

	// 只生成业务表的 model
	var models []interface{}
	for _, table := range tables {
		structName := tableToStructName(table)
		model := g.GenerateModelAs(table, structName)
		models = append(models, model)
	}

	// 生成 dao 到 dao 目录
	//g.ApplyBasic(models...)
	g.ApplyBasic(
		&model.User{},
		&model.Paper{},
		&model.Question{},
		&model.PaperQuestion{},
	)

	// 执行代码生成
	g.Execute()

}
