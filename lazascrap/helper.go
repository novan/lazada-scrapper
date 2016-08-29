package lazascrap

import (
	"errors"
	"fmt"
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func HttpRequest(BaseUrl string, param *map[string]string) (content []byte, err error) {

	finalUrl, err := buildUrl(BaseUrl, param)

	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", finalUrl.String(), nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", randomUserAgents())
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	content, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return
}

func ScrapPage(BaseUrl string, page int) (product *Product, err error) {
	param := make(map[string]string)
	param["itemperpage"] = "120"
	param["page"] = fmt.Sprintf("%d", page)

	content, err := HttpRequest(BaseUrl, &param)

	doc, err := gokogiri.ParseHtml(content)

	if err != nil {
		return nil, err
	}

	// getting total items
	const schXPath = "/html/body/div[2]/div[2]/div[2]/div[1]/div[2]/div"

	root := doc.Root()

	if root == nil {
		return nil, errors.New("Empty document?")
	}

	html := doc.Root().FirstChild()

	results, err := html.Search(schXPath)
	if err != nil {
		return nil, err
	} else if len(results) <= 0 {
		return nil, errors.New("Parsing failed")
	}

	defer doc.Free()

	product = &Product{}
	product.Items = make([]ProductItem, len(results)-1)

	log.Println("Length Result: ", len(results))
	for i, result := range results[0 : len(results)-1] {
		log.Println("Parsing result-", i, ":", results[i])
		product.Items[i], _ = trNodeToProduct(result)
	}
	log.Println("Products: ", product.Items)

	// getting total data
	const totalSchXPath = "/html/body/div[2]/div[2]/div[2]/div[1]/div[1]/div[1]/div/span/text()"

	totalResults, err := html.Search(totalSchXPath)

	if err != nil {
		return nil, err
	}
	product.TotalItems, err =
		strconv.Atoi(strings.Split(totalResults[0].String(), " ")[0])

	if err != nil {
		return nil, err
	}

	return
}

func mapToQuery(m *map[string]string) string {
	if m == nil || len(*m) == 0 {
		return ""
	} else {
		params := url.Values{}
		for k, v := range *m {
			params.Add(k, v)
		}
		return params.Encode()
	}
}

func buildUrl(base string, qs *map[string]string) (url *url.URL, err error) {
	baseUrl, err := url.Parse(base)
	if err != nil {
		return nil, err
	} else {
		if qs == nil {
			return baseUrl, nil
		} else {
			baseUrl.RawQuery = mapToQuery(qs)
			return baseUrl, nil
		}
	}
}

var rndInitialised bool = false

func randomUserAgents() string {

	if !rndInitialised {
		rand.Seed(time.Now().UTC().UnixNano())
	}
	// Randomise User Agents just for fun, we'll use Console's UA, and OLD OS
	var userAgents = [...]string{
		"Mozilla/5.0 (PlayStation 4 2.57) AppleWebKit/536.26 (KHTML, like Gecko)",
		"Opera/9.50 (Nintendo DSi; Opera/507; U; en-US)",
		"Mozilla/5.0 (Nintendo 3DS; U; ; en) Version/1.7498.US",
		"AmigaVoyager/3.2 (AmigaOS/MC680x0)",
		"NCSA Mosaic/3.0 (Windows 95)",
		"Mozilla/3.0 (Planetweb/2.100 JS SSL US; Dreamcast US)",
		"Mozilla/5.0 (PlayStation Vita 1.80) AppleWebKit/531.22.8 (KHTML, like Gecko) Silk/3.2",
		"Mozilla/4.0 (compatible; MSIE 6.1; Windows XP; .NET CLR 1.1.4322; .NET CLR 2.0.50727)",
	}
	idx := rand.Intn(len(userAgents))
	return userAgents[idx]
}

func trNodeToProduct(productNode xml.Node) (item ProductItem, err error) {

	results, err := productNode.Search("./a/div[3]/div[1]/span/text()")
	productTitle := results[0].String()

	results, err = productNode.Search("./a/div[3]/div[2]/div[3]/div/text()")
	productOriginalPrice := ""
	productDiscountedPrice := ""
	if len(results) == 0 {
		results, err = productNode.Search("./a/div[3]/div[2]/div[1]/text()")
		if len(results) > 0 {
			productOriginalPrice = results[0].String()
		}
	} else {
		productOriginalPrice = results[0].String()

		results, err = productNode.Search("./a/div[3]/div[2]/div[1]/text()")
		if len(results) > 0 {
			productDiscountedPrice = results[0].String()
		}
	}

	results, err = productNode.Search("./a/div[2]/img")
	productImage := results[0].Attributes()["data-original"].String()

	if err != nil {
		return ProductItem{}, err
	}

	item = ProductItem{
		Title:           productTitle,
		Price:           productOriginalPrice,
		DiscountedPrice: productDiscountedPrice,
		Image:           productImage,
	}

	return
}
