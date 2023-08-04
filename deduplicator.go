package main

import "idempotency/model"

func unique(sample []model.Product) []model.Product {
	var unique []model.Product
	type key struct{ value1, value2, value3 string } // choose the key for deduplication
	m := make(map[key]int)
	for _, v := range sample {
		k := key{v.BrandID, v.Name, v.ParentSKU}
		if i, ok := m[k]; ok {
			// Overwrite previous value per requirement in
			// question to keep last matching value.
			unique[i] = v
		} else {
			// Unique key found. Record position and collect
			// in result.
			m[k] = len(unique)
			unique = append(unique, v)
		}
	}
	return unique
}
