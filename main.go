package main

import "fmt"

func main() {
	baseURL := "https://azami-jp.com/live"

	// Goquery
	resp, err := fetch(baseURL)

	if err != nil {
		panic(err)
	}

	infosByGoquery, err := parseByGoquery(resp)
	if err != nil {
		panic(err)
	}

	fmt.Println("Goquery")
	for _, item := range infosByGoquery {
		fmt.Println(item)
		// fmt.Println(info, info.URL) // TODO ランタイムのバグか？出力が異なる
	}

	// Colly
	infosByColly, err := parseByColly(baseURL)
	if err != nil {
		panic(err)
	}

	fmt.Println("====================================")
	fmt.Println("Colly")
	for _, info := range infosByColly {
		fmt.Println(info)
		// fmt.Println(info, info.URL) // TODO ランタイムのバグか？出力が異なる
	}
}
