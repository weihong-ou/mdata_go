# Summary
This document provides a brief specification of **MDATA** Rest service endpoints built specifically for Demo puprose. Currently, the Rest service provides endpoints to perform:
* Create a product
* Update a product
* Set state for the specified product (valid states are `active`, `inactive`, and `discontinued`)
* Delete a product
* List all products currently saved on chain

## Go Rest service development framework used
All Rest service endpoints are build via **[Echo](https://echo.labstack.com/)** Go web framework.

## Endpoints Specification
Each product is uniquely identified via its 14-digit **GTIN** number. All service endpoints except *List All Product* require GTIN number to be the last segment of the endpoit path

### Create a Product
HTTP Method: **POST**
Endpoint Path: /product/create/:gtin (e.g. /product/create/*12345678901234*)
Request Payload:
```JSON
{
    "gtin": "12345678901234",
    "attrs": {
        "name": "Chicken Wing",
        "desc": "Tyson chicken wing",
        "qty": 10,
        "uom": "case"
    },
    "state": "active"
}
```
Response: `201 Created`
```JSON
{
    "gtin": "12345678901234",
    "attrs": {
        "name": "Chicken Wing",
        "desc": "Tyson chicken wing",
        "qty": 10,
        "uom": "case"
    },
    "state": "active"
}
```

### Update a Product
HTTP Method: **PUT**
Endpoint Path: /product/update/:gtin (e.g. /product/update/*12345678901234*)
Request Payload:
```JSON
{
    "gtin": "12345678901234",
    "attrs": {
        "name": "Chicken Wing",
        "desc": "Tyson chicken wing",
        "qty": 150,
        "uom": "case"
    },
    "state": "active"
}
```
Response: `200 OK`
```JSON
{
    "gtin": "12345678901234",
    "attrs": {
        "name": "Chicken Wing",
        "desc": "Tyson chicken wing",
        "qty": 150,
        "uom": "case"
    },
    "state": "active"
}
```

### Set Product State
HTTP Method: **PUT**
Endpoint Path: /product/setstate/:gtin (e.g. /product/setstate/*12345678901234*)
Request Payload:
```JSON
{
    "new_state": "discontinued"
}
```
Response: `200 OK`
```JSON
{
    "message": "Product 12345678901234 state has been set to DISCONTINUED"
}
```

### Delete a Product
HTTP Method: **DELETE**
Endpoint Path: /product/delete/:gtin (e.g. /product/delete/*12345678901234*)
Request Payload: *empty*
Response: `200 OK`
```JSON
{
    "message": "Product 12345678901238 deleted."
}
```

### List Products
HTTP Method: **GET**
Endpoint Path: /product/list
Request Payload: *empty*
Response: `200 OK`
```JSON
[
    {
        "gtin": "12345678901234",
        "attrs": {
            "name": "Chicken Wing",
            "desc": "Tyson chicken wing",
            "qty": 10,
            "uom": "case"
        },
        "state": "active"
    },
    {
        "gtin": "12345678905678",
        "attrs": {
            "name": "Beef Steak",
            "desc": "Tyson beef steak",
            "qty": 100,
            "uom": "pound"
        },
        "state": "discontinued"
    }
]
```