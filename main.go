package main

import (
  "fmt"
  "net/http"
  "html/template"
)

type User struct {
  Name string //char?
  Age uint16 // >0
  Trust_lvl int16 //can be less then 0
  Avg_g, Hpnss float64
  Skills [] string
}

func (u User)  getAllInfo() string {
  return fmt.Sprintf("Username is %s. She is %d, her trust " +
    "level is %d", u.Name, u.Age, u.Trust_lvl)

}

func (u *User) setNewName(newName string)  {
  u.Name = newName
}

func index(w http.ResponseWriter, r *http.Request)  {
  // u1 := User{"Nancee", 25, 100, 4.1, 0, []string{"C", "Go", "Oil Painting"}}

  t, err := template.ParseFiles("templates/header.html", "templates/footer.html", "templates/index.html")

  if err != nil {
    fmt.Fprintf(w, err.Error())
    return
  }

  t.ExecuteTemplate(w, "index", nil)
}

func handleFunc()  {
  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
  http.HandleFunc("/", index)
  // http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
  http.ListenAndServe(":8080", nil)
}


func main(){

  handleFunc()
}
