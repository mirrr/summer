# Summer panel
Simple control panel for [Golang](https://golang.org/) based on [Gin framework](https://gin-gonic.github.io/gin/) and [MongoDB](https://www.mongodb.com/)
    

## How To Install   
```bash
go get -u gopkg.in/night-codes/summer.v1/...

```


## Getting Started

1) Create new project with demo-modules:
```bash
summerGen project --name myProject --title="My project" --db "project" --port=8080 --views="templates/main" --views-dot="templates/dot" --demo
```
   
    
2) Create new module:
```bash
cd myProject/
summerGen module --name tasks --title="My tasks" --menu=MainMenu --add-sort --add-search --add-pages --add-tabs
```

3) Start project
```bash
go build
./myProject
```

### Result: 
![summer](https://cloud.githubusercontent.com/assets/2770221/21749479/7d293c42-d5d1-11e6-846a-654afbf1288a.png)


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
