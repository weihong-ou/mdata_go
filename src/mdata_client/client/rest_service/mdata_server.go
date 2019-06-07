package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type ProductState int

const (
	INACTIVE     ProductState = 0
	ACTIVE       ProductState = 1
	DISCONTINUED ProductState = 2
)

func (ps ProductState) String() string {
	switch ps {
	case INACTIVE:
		return "INACTIVE"
	case ACTIVE:
		return "ACTIVE"
	case DISCONTINUED:
		return "DISCONTINUED"
	default:
		return "UNKNOWN"
	}
}

func convertToState(s string) (ProductState, error) {
	switch strings.ToUpper(s) {
	case "INACTIVE":
		return INACTIVE, nil
	case "ACTIVE":
		return ACTIVE, nil
	case "DISCONTINUED":
		return DISCONTINUED, nil
	default:
		return INACTIVE, fmt.Errorf("Unknown product state string: %s", s)

	}
}

type ProductAttributes struct {
	Name        string `json:"name" xml:"name" form:"name" query:"name"`
	Description string `json:"desc" xml:"desc" form:"desc" query:"desc"`
	Quantity    int    `json:"qty" xml:"qty" form:"qty" query:"qty"`
	UOM         string `json:"uom" xml:"uom" form:"uom" query:"uom"`
}

type Product struct {
	Gtin       string            `json:"gtin" xml:"gtin" form:"gtin" query:"gtin"`
	Attributes ProductAttributes `json:"attrs" xml:"attrs" form:"attrs" query:"attrs"`
	State      string            `json:"state" xml:"state" form:"state" query:"state"`
}

type ProductStateUpdate struct {
	State string `json:"new_state" xml:"new_state" form:"new_state" query:"new_state"`
}

func (p Product) isValid() bool {
	gtin := p.Gtin
	pattern := regexp.MustCompile(`^\d{14}$`)
	if !pattern.MatchString(gtin) {
		return false
	}

	//attributes are all optional, so no need to validate
	_, err := convertToState(p.State)
	return err == nil
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.CORS()) //for now open to all origins

	e.POST("/product/create/:gtin", func(c echo.Context) error {
		p := Product{}
		if err := c.Bind(&p); err != nil {
			//first logs the error
			e.Logger.Errorf("error binding request product: %v", err)
			return echo.NewHTTPError(http.StatusBadRequest, "error binding request message to product")
		}
		p.Gtin = c.Param("gtin")
		if !p.isValid() {
			return echo.NewHTTPError(http.StatusBadRequest, "product payload missing required fields or has invalid field values")
		}

		//DO THE NORMAL MDATA CLIENT PRODUCT CREATEION LOGIC HERE

		return c.JSON(http.StatusCreated, p)
	})

	e.PUT("/product/update/:gtin", func(c echo.Context) error {
		p := new(Product)
		if err := c.Bind(p); err != nil {
			//first logs the error
			e.Logger.Errorf("error binding request product: %v", err)
			return echo.NewHTTPError(http.StatusBadRequest, "error binding request message to product")
		}
		p.Gtin = c.Param("gtin")
		if !p.isValid() {
			return echo.NewHTTPError(http.StatusBadRequest, "product payload missing required fields or has invalid field values")
		}

		//DO THE NORMAL MDATA CLIENT PRODUCT UPDATE LOGIC HERE
		//WE MAY NEED TO DIFFERANTIATE THE ERROR OF NON-EXISTING GTIN AND
		//CHAIN STATE WRITE ERROR (ONE IS 404 AND THE OTHER IS 500 RESPONSE CODE)
		gtin := c.Param("gtin")
		e.Logger.Debugf("Update product (%s) state to %v", gtin, p)

		//assuming everything is fine
		return c.JSON(http.StatusOK, p)
	})

	e.PUT("/product/setstate/:gtin", func(c echo.Context) error {
		//assuming the json payload has the following format {new_state: "discontinued"}
		ps := new(ProductStateUpdate)
		if err := c.Bind(ps); err != nil {
			//first logs the error
			e.Logger.Errorf("error binding request product state: %v", err)
			return echo.NewHTTPError(http.StatusBadRequest, "error binding request message to product state")
		}

		//validate the state string
		pState, err := convertToState(ps.State)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid product state: "+ps.State)
		}

		//DO THE NORMAL MDATA CLIENT PRODUCT STATE UPDATE LOGIC HERE
		//WE MAY NEED TO DIFFERANTIATE THE ERROR OF NON-EXISTING GTIN AND
		//CHAIN STATE WRITE ERROR (ONE IS 404 AND THE OTHER IS 500 RESPONSE CODE)
		gtin := c.Param("gtin")
		e.Logger.Debugf("Update product (%s) state to %v", gtin, pState)

		//assuming everything is fine
		return c.JSON(http.StatusOK, &struct {
			Message string `json:"message"`
		}{Message: fmt.Sprintf("Product %s state has been set to %v", gtin, pState)})
	})

	e.DELETE("/product/delete/:gtin", func(c echo.Context) error {
		gtin := c.Param("gtin")
		e.Logger.Debugf("Delete product %s", gtin)

		//DO THE NORMAL MDATA CLIENT PRODUCT DELETION LOGIC HERE
		//WE MAY NEED TO DIFFERANTIATE THE ERROR OF NON-EXISTING GTIN AND
		//CHAIN STATE WRITE ERROR (ONE IS 404 AND THE OTHER IS 500 RESPONSE CODE)

		//assuming everything is ok
		return c.JSON(http.StatusOK, &struct {
			Message string `json:"message"`
		}{Message: fmt.Sprintf("Product %s deleted.", gtin)})
	})

	e.GET("/product/list", func(c echo.Context) error {
		//DO THE NORMAL MDATA CLIENT PRODUCT LISTING LOGIC HERE

		//for now, just a mockup list of products:
		plist := []Product{
			Product{
				Gtin: "12345678901234",
				Attributes: ProductAttributes{
					Name:        "Chicken Wing",
					Description: "Tyson chicken wing",
					Quantity:    10,
					UOM:         "case",
				},
				State: "active",
			},
			Product{
				Gtin: "12345678905678",
				Attributes: ProductAttributes{
					Name:        "Beef Steak",
					Description: "Tyson beef steak",
					Quantity:    100,
					UOM:         "pound",
				},
				State: "discontinued",
			},
		}

		return c.JSON(http.StatusOK, plist)
	})

	port := 8888
	var err error
	if len(os.Args) > 1 {
		port, err = strconv.Atoi(os.Args[1])
		if err != nil {
			port = 8888
			fmt.Fprintf(os.Stderr, "port should be an integer, using default port 8888")
		}
	}
	e.Logger.Fatal(e.Start(":" + strconv.Itoa(port)))

}
