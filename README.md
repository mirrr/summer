# Summer panel
Simple control panel for [Golang](https://golang.org/) based on [Gin framework](https://gin-gonic.github.io/gin/) and [MongoDB](https://www.mongodb.com/)
    

## How To Install   
```bash
go get gopkg.in/night-codes/summer.v1
```


## Getting Started

```go
package main

import (
    "gopkg.in/night-codes/summer.v1"
)

var (
    panel = summer.Create(summer.Settings{
        Port:        8080,
        AuthSalt:    "myappSalt123",
        AuthPrefix:  "myapp-",
        DefaultPage: "news",
        Path:        "/panel",
        DBName:      "mypanel",
        Views:       "views",
        Files:       "files",
    })
)

func main() {
    summer.Wait()
}

```
   

## Create section in Summer
 
**/news.go**   
   
```go
package main

import (
    "gopkg.in/gin-gonic/gin.v1"
    "gopkg.in/night-codes/summer.v1"
)

type (
    obj map[string]interface{}

    NewsModule struct {
        summer.Module
    }
)

var (
    subMenu = panel.MainMenu.Add("Submenu", 1)
    news = panel.AddModule(
        &summer.ModuleSettings{
            Name:         "news",
            Title:        "News",
            Menu:         subMenu, // or panel.MainMenu
            MenuOrder:    1,
            TemplateName: "news/index",
        },
        &NewsModule{},
    )
)

func (module *NewsModule) Page(c *gin.Context) {
    settings := module.Settings
    c.HTML(200, settings.TemplateName+".html", gin.H{
        "title": settings.Title,
        "user":  c.MustGet("user"), // must be for correct username in the header
        "data":  obj{"text": "This is backend", "check": "Check Me"},
    })
}
```
   
**/views/news/index.html**

```html
<h1>{{.title}}:</h1>
<div class="right-panel">
    <button id="addItem"><span class="fa fa-plus"></span> Add news</button>
</div>
<div class="wow">
    <img src="https://goo.gl/lLmbVR" />
    <h3>{{.data.text}}</h3>
    <input type="checkbox" id="c" checked /> {{.data.check}}
</div>
```
   
**Result:**
(On [http://localhost:8080/panel/news/](http://localhost:8080/panel/news/) )
![Summer screenshot](https://cloud.githubusercontent.com/assets/2770221/20869933/3a89113e-ba86-11e6-9a22-2967cb1eac05.png)

## Examples
Coming soon...
   
   
## People

Author and developer is [Oleksiy Chechel](https://github.com/night-codes)    
   


## License
   
MIT License   
   
Copyright (C) 2016-2017 Oleksiy Chechel (alex.mirrr@gmail.com)   
   
Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:   
   
The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.   
   
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
