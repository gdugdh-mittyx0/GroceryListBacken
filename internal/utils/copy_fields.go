package utils

import (
	"errors"
	"fmt"
	"reflect"
)

func CopyFields(src interface{}, dest interface{}) error {
	if src == nil || dest == nil {
		return errors.New("src and dest cannot be nil")
	}

	srcVal := reflect.ValueOf(src)
	destVal := reflect.ValueOf(dest)

	if srcVal.Kind() != reflect.Ptr || destVal.Kind() != reflect.Ptr {
		return errors.New("src and dest must be pointers to structs")
	}

	srcVal = srcVal.Elem()
	destVal = destVal.Elem()

	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Type().Field(i)
		destField := destVal.FieldByName(srcField.Name)

		// Проверка на тип поля
		if destField.IsValid() && destField.CanSet() {
			if destField.Type() == srcField.Type {
				destField.Set(srcVal.Field(i))
			} else {
				return errors.New("err_different_types_fields_struct")
			}
		}
	}
	return nil
}

func СopyIfExist(src, dest interface{}) error {
	if src == nil || dest == nil {
		return errors.New("src and dest cannot be nil")
	}

	srcVal := reflect.ValueOf(src)
	destVal := reflect.ValueOf(dest)

	if srcVal.Kind() != reflect.Ptr || destVal.Kind() != reflect.Ptr {
		return errors.New("src and dest must be pointers to structs")
	}

	srcVal = srcVal.Elem()
	destVal = destVal.Elem()

	for i := 0; i < destVal.NumField(); i++ {
		destField := destVal.Type().Field(i)
		srcField := srcVal.FieldByName(destField.Name)

		// Проверяем, существует ли поле в src
		if srcField.IsValid() {
			// Является ли оно указателем
			if srcField.Kind() == reflect.Ptr {
				// Если указатель не nil, копируем значение
				if !srcField.IsNil() {
					fmt.Printf("Копирование поля: %s\n", destField.Name)
					destFieldValue := destVal.Field(i)
					if destFieldValue.CanSet() {
						destFieldValue.Set(srcField.Elem())
					}
				}
			} else {
				// Если поле не указатель, просто копируем значение
				destFieldValue := destVal.Field(i)
				if destFieldValue.CanSet() {
					destFieldValue.Set(srcField)
				}
			}
		}
	}

	return nil
}
